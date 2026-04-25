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
	if os.Getenv("DATABASE_URL") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := db.Open(ctx)
		if err != nil {
			logging.Logger.Warn("PostgreSQL is unavailable, falling back to in-memory transaction repository", zap.Error(err))
		} else {
			defer conn.Close()
			txRepo = repositories.NewPostgresTransactionRepository(conn)
		}
	}
	txSvc := services.NewTransactionService(txRepo)
	txCtrl := controllers.NewTransactionController(txSvc)

	// Настройка маршрутов
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// Эндпоинт для инициации транзакции
	router.POST("/transaction/initiate", txCtrl.InitiateTransaction)

	api := router.Group("/api/transaction")
	api.POST("/initiate", txCtrl.InitiateTransaction)
	api.POST("/submit", txCtrl.SubmitTransaction)
	api.GET("/:transaction_id/status", txCtrl.GetStatus)
	api.POST("/:transaction_id/cancel", txCtrl.CancelTransaction)
	api.POST("/:transaction_id/submit", txCtrl.SubmitTransaction)
	api.GET("/:transaction_id", txCtrl.GetTransaction)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
