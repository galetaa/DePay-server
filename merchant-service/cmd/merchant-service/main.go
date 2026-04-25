package main

import (
	"log"
	"os"

	"merchant-service/internal/controllers"
	"merchant-service/internal/repositories"
	"merchant-service/internal/services"
	"shared/config"
	"shared/logging"
	"shared/middleware"

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
	svc := services.NewMerchantService(repo)
	ctrl := controllers.NewMerchantController(svc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
