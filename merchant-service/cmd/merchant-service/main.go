package main

import (
	"context"
	"log"
	"os"
	"time"

	"merchant-service/internal/controllers"
	"merchant-service/internal/repositories"
	"merchant-service/internal/services"
	"shared/config"
	"shared/db"
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

	repo := repositories.NewMerchantRepository()
	if os.Getenv("DATABASE_URL") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := db.Open(ctx)
		if err != nil {
			logging.Logger.Warn("PostgreSQL is unavailable, falling back to in-memory merchant repository", zap.Error(err))
		} else {
			defer conn.Close()
			repo = repositories.NewPostgresMerchantRepository(conn)
		}
	}
	svc := services.NewMerchantService(repo)
	ctrl := controllers.NewMerchantController(svc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("merchant-service"))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.Handler())

	api := router.Group("/api/merchant")
	api.POST("/register", ctrl.Register)
	api.POST("/login", ctrl.Login)
	protected := api.Group("")
	protected.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("merchant"))
	protected.GET("/me", ctrl.Me)
	protected.POST("/verification", ctrl.SubmitVerification)
	protected.GET("/verification/status", ctrl.VerificationStatus)
	protected.POST("/invoices", ctrl.CreateInvoice)
	protected.GET("/invoices", ctrl.ListInvoices)
	protected.POST("/terminals", ctrl.CreateTerminal)
	protected.GET("/terminals", ctrl.ListTerminals)
	protected.POST("/webhooks", ctrl.CreateWebhook)
	protected.GET("/webhooks", ctrl.ListWebhooks)
	protected.GET("/webhooks/:webhook_id", ctrl.GetWebhook)
	protected.POST("/webhooks/:webhook_id/test", ctrl.TestWebhook)
	protected.DELETE("/webhooks/:webhook_id", ctrl.DeleteWebhook)
	protected.POST("/api-keys", ctrl.CreateAPIKey)
	protected.GET("/api-keys", ctrl.ListAPIKeys)
	protected.DELETE("/api-keys/:key_id", ctrl.RevokeAPIKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
