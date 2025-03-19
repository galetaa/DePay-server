package main

import (
	"log"
	"os"

	"kyc-service/internal/controllers"
	"kyc-service/internal/services"
	"kyc-service/shared/config"
	"kyc-service/shared/logging"
	"kyc-service/shared/middleware"

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
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Эндпоинт для проверки работоспособности
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Эндпоинт для обработки KYC запросов
	router.POST("/kyc", kycCtrl.ProcessKYC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера с использованием zap.Error для корректного логирования ошибки
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
