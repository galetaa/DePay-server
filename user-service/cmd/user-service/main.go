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
	userRepo := repositories.NewUserRepository()
	if os.Getenv("DATABASE_URL") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := db.Open(ctx)
		if err != nil {
			logging.Logger.Warn("PostgreSQL is unavailable, falling back to in-memory user repository", zap.Error(err))
		} else {
			defer conn.Close()
			userRepo = repositories.NewPostgresUserRepository(conn)
		}
	}
	userSvc := services.NewUserService(userRepo)
	userCtrl := controllers.NewUserController(userSvc)

	// Инициализируем Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("user-service"))

	// Эндпоинт проверки работоспособности сервиса
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.Handler())

	// Определяем маршруты для работы с пользователями
	router.POST("/user/register", userCtrl.Register)
	router.POST("/user/login", userCtrl.Login)
	router.POST("/user/refresh-token", userCtrl.RefreshToken)

	api := router.Group("/api")
	userAPI := api.Group("/user")
	userAPI.POST("/register", userCtrl.Register)
	userAPI.POST("/login", userCtrl.Login)
	userAPI.POST("/refresh-token", userCtrl.RefreshToken)
	protectedUserAPI := userAPI.Group("")
	protectedUserAPI.Use(middleware.JWTAuthMiddleware())
	protectedUserAPI.GET("/me", userCtrl.Me)
	protectedUserAPI.PUT("/me", userCtrl.UpdateMe)
	protectedUserAPI.POST("/kyc", userCtrl.SubmitKYC)
	protectedUserAPI.GET("/kyc/status", userCtrl.KYCStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
