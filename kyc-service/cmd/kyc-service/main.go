package main

import (
	"log"
	"os"

	"kyc-service/internal/controllers"
	"kyc-service/internal/services"
	"shared/config"
	"shared/logging"
	"shared/middleware"
	"shared/observability"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Загружаем конфигурацию
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логгер
	if err := logging.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logging.Sync()

	// Устанавливаем режим Gin для продакшена
	gin.SetMode(gin.ReleaseMode)

	// Инициализируем сервисы
	kycSvc := services.NewKYCService()
	kycCtrl := controllers.NewKYCController(kycSvc)

	// Инициализируем маршрутизатор
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("kyc-service"))

	// Эндпоинт для проверки работоспособности
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.Handler())

	// Эндпоинт для обработки KYC запросов
	router.POST("/kyc", kycCtrl.ProcessKYC)
	router.POST("/api/kyc", kycCtrl.ProcessKYC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера с использованием zap.Error для корректного логирования ошибки
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
