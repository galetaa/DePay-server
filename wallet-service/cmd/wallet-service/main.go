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
	if os.Getenv("DATABASE_URL") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := db.Open(ctx)
		if err != nil {
			logging.Logger.Warn("PostgreSQL is unavailable, falling back to in-memory wallet repository", zap.Error(err))
		} else {
			defer conn.Close()
			walletRepo = repositories.NewPostgresWalletRepository(conn)
		}
	}
	walletSvc := services.NewWalletService(walletRepo)
	walletCtrl := controllers.NewWalletController(walletSvc)

	// Настройка маршрутов
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("wallet-service"))

	// Эндпоинт проверки работоспособности сервиса
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.Handler())

	// Эндпоинты Wallet Service
	router.GET("/wallet/export", walletCtrl.ExportWallets)
	router.POST("/wallet/balance", walletCtrl.GetBalance)

	api := router.Group("/api")
	wallets := api.Group("/wallets")
	wallets.Use(middleware.JWTAuthMiddleware())
	wallets.POST("", walletCtrl.CreateWallet)
	wallets.GET("", walletCtrl.ExportWallets)
	wallets.GET("/:wallet_id/balances", walletCtrl.GetWalletBalances)
	wallets.POST("/:wallet_id/sync", walletCtrl.SyncWallet)
	wallets.GET("/:wallet_id/balance", walletCtrl.GetBalanceByAddress)
	wallets.GET("/:wallet_id", walletCtrl.GetWallet)
	wallets.DELETE("/:wallet_id", walletCtrl.DeleteWallet)

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
