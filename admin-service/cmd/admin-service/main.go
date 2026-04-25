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
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/admin/tables", adminCtrl.ListTables)
	api.GET("/admin/tables/:table_name", adminCtrl.GetTableRows)
	api.POST("/admin/functions/:function_name/execute", adminCtrl.ExecuteFunction)
	api.GET("/admin/audit-logs", adminCtrl.AuditLogs)
	api.GET("/admin/risk-alerts", adminCtrl.RiskAlerts)
	api.GET("/analytics/store-turnover", adminCtrl.StoreTurnover)
	api.GET("/analytics/transaction-statuses", adminCtrl.TransactionStatuses)
	api.GET("/analytics/failed-transactions", adminCtrl.FailedTransactions)
	api.GET("/analytics/rpc-health", adminCtrl.RPCHealth)
	api.POST("/admin/demo/invoices", adminCtrl.CreateDemoInvoice)
	api.POST("/admin/demo/payments", adminCtrl.SubmitDemoPayment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	if err := router.Run(":" + port); err != nil {
		logging.Logger.Fatal("Failed to run server", zap.Error(err))
	}
}
