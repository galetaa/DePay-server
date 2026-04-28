package main

import (
	"log"
	"os"

	"gas-info-service/internal/controllers"
	"gas-info-service/internal/services"
	"shared/config"
	"shared/logging"
	"shared/middleware"
	"shared/observability"

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

	gasSvc := services.NewGasService()
	gasCtrl := controllers.NewGasController(gasSvc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("gas-info-service"))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.Handler())
	// Эндпоинт для получения информации о газе. Ожидается query параметр "network"
	router.GET("/gas-info", gasCtrl.GetGasInfo)
	router.GET("/api/transactions/gas-info", gasCtrl.GetGasInfo)
	router.GET("/api/gas-info/history", gasCtrl.GetGasHistory)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
