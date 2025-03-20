package main

import (
	"log"
	"os"

	"shared/config"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
