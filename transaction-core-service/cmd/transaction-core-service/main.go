package main

import (
	"context"
	"log"
	"os"
	"time"

	"shared/config"
	"shared/db"
	"shared/logging"
	"shared/middleware"
	"shared/observability"
	"transaction-core-service/internal/controllers"
	"transaction-core-service/internal/repositories"
	"transaction-core-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Загрузка конфигурации и инициализация логирования
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := logging.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logging.Sync()

	gin.SetMode(gin.ReleaseMode)

	// Инициализация зависимостей
	txRepo := repositories.NewTransactionRepository()
	var txOptions []services.Option
	if os.Getenv("DATABASE_URL") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := db.Open(ctx)
		if err != nil {
			logging.Logger.Warn("PostgreSQL is unavailable, falling back to in-memory transaction repository", zap.Error(err))
		} else {
			defer conn.Close()
			txRepo = repositories.NewPostgresTransactionRepository(conn)
			txOptions = append(txOptions, services.WithWebhookDispatcher(services.NewPostgresWebhookDispatcher(conn)))
		}
	}
	txSvc := services.NewTransactionService(txRepo, txOptions...)
	txCtrl := controllers.NewTransactionController(txSvc)

	// Настройка маршрутов
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("transaction-core-service"))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.Handler())
	legacyProtected := router.Group("/transaction")
	legacyProtected.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("merchant", "admin"))
	legacyProtected.POST("/initiate", txCtrl.InitiateTransaction)

	api := router.Group("/api/transaction")
	api.Use(middleware.JWTAuthMiddleware())
	api.GET("/:transaction_id/status", txCtrl.GetStatus)
	api.GET("/:transaction_id", txCtrl.GetTransaction)
	protected := api.Group("")
	protected.Use(middleware.RoleMiddleware("merchant", "admin"))
	protected.POST("/initiate", txCtrl.InitiateTransaction)
	protected.POST("/submit", txCtrl.SubmitTransaction)
	protected.POST("/:transaction_id/cancel", txCtrl.CancelTransaction)
	protected.POST("/:transaction_id/submit", txCtrl.SubmitTransaction)
	protected.POST("/:transaction_id/validate", txCtrl.ValidateTransaction)
	protected.POST("/:transaction_id/broadcast", txCtrl.BroadcastTransaction)
	protected.POST("/:transaction_id/confirm", txCtrl.ConfirmTransaction)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
