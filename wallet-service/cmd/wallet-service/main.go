package main

import (
	"log"
	"os"

	"shared/config"
	"shared/logging"
	"shared/middleware"
	"wallet-service/internal/controllers"
	"wallet-service/internal/repositories"
	"wallet-service/internal/services"

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

	// Запускаем Gin в Release-режиме
	gin.SetMode(gin.ReleaseMode)

	// Инициализируем зависимости
	walletRepo := repositories.NewWalletRepository()
	walletSvc := services.NewWalletService(walletRepo)
	walletCtrl := controllers.NewWalletController(walletSvc)

	// Настройка маршрутов
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Эндпоинт проверки работоспособности сервиса
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Эндпоинты Wallet Service
	router.GET("/wallet/export", walletCtrl.ExportWallets)
	router.POST("/wallet/balance", walletCtrl.GetBalance)

	// Чтение порта из переменной окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
