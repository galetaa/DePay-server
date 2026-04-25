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
	"transaction-validation-service/internal/controllers"
	"transaction-validation-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := logging.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logging.Sync()

	gin.SetMode(gin.ReleaseMode)

	validationSvc := services.NewValidationService()
	if os.Getenv("DATABASE_URL") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := db.Open(ctx)
		if err != nil {
			logging.Logger.Warn("PostgreSQL is unavailable, using stateless transaction validation", zap.Error(err))
		} else {
			defer conn.Close()
			validationSvc = services.NewPostgresValidationService(conn)
		}
	}
	validationCtrl := controllers.NewValidationController(validationSvc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// Эндпоинт для проверки валидности транзакции
	router.POST("/transaction/validate", validationCtrl.ValidateTransaction)
	router.POST("/api/transaction/validate", validationCtrl.ValidateTransaction)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
