package main

import (
	"log"
	"os"

	"shared/config"
	"shared/logging"
	"shared/middleware"
	"user-service/internal/controllers"
	"user-service/internal/repositories"
	"user-service/internal/services"

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

	// Устанавливаем режим Gin (в продакшене — Release)
	gin.SetMode(gin.ReleaseMode)

	// Инициализируем слои приложения
	userRepo := repositories.NewUserRepository() // В продакшене подключение к PostgreSQL через GORM/sqlx
	userSvc := services.NewUserService(userRepo)
	userCtrl := controllers.NewUserController(userSvc)

	// Инициализируем Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Эндпоинт проверки работоспособности сервиса
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Определяем маршруты для работы с пользователями
	router.POST("/user/register", userCtrl.Register)
	router.POST("/user/login", userCtrl.Login)
	router.POST("/user/refresh-token", userCtrl.RefreshToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
