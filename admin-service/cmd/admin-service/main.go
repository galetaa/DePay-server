package main

import (
	"context"
	"log"
	"os"
	"time"

	"admin-service/internal/controllers"
	"admin-service/internal/services"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := db.Open(ctx)
	if err != nil {
		logging.Logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer conn.Close()

	gin.SetMode(gin.ReleaseMode)

	adminSvc := services.NewAdminService(conn)
	adminCtrl := controllers.NewAdminController(adminSvc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(observability.Middleware("admin-service"))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/metrics", observability.DatabaseHandler(conn))

	api := router.Group("/api")
	adminAPI := api.Group("/admin", middleware.ProductionRoleGate("admin"))
	adminAPI.GET("/tables", adminCtrl.ListTables)
	adminAPI.GET("/tables/:table_name", adminCtrl.GetTableRows)
	adminAPI.POST("/functions/:function_name/execute", adminCtrl.ExecuteFunction)
	adminAPI.GET("/audit-logs", adminCtrl.AuditLogs)
	adminAPI.GET("/risk-alerts", adminCtrl.RiskAlerts)
	adminAPI.GET("/system-health", adminCtrl.SystemHealth)

	analyticsAPI := api.Group("/analytics", middleware.ProductionRoleGate("admin", "compliance"))
	analyticsAPI.GET("/store-turnover", adminCtrl.StoreTurnover)
	analyticsAPI.GET("/transaction-statuses", adminCtrl.TransactionStatuses)
	analyticsAPI.GET("/failed-transactions", adminCtrl.FailedTransactions)
	analyticsAPI.GET("/rpc-health", adminCtrl.RPCHealth)

	demoAPI := adminAPI.Group("/demo", middleware.DemoEndpointGuard())
	demoAPI.POST("/invoices", adminCtrl.CreateDemoInvoice)
	demoAPI.POST("/payments", adminCtrl.SubmitDemoPayment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
