package controllers

import (
	"net/http"
	"strconv"

	"admin-service/internal/models"
	"admin-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	service services.AdminService
}

func NewAdminController(service services.AdminService) *AdminController {
	return &AdminController{service: service}
}

func (ac *AdminController) ListTables(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": models.TableListResponse{Tables: ac.service.ListTables()}})
}

func (ac *AdminController) GetTableRows(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	rows, err := ac.service.GetTableRows(c.Request.Context(), c.Param("table_name"), limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) ExecuteFunction(c *gin.Context) {
	var req models.ExecuteFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": "invalid request", "details": gin.H{}}})
		return
	}

	rows, err := ac.service.ExecuteFunction(c.Request.Context(), c.Param("function_name"), req.Params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) AuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	rows, err := ac.service.AuditLogs(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "QUERY_FAILED", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) RiskAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	rows, err := ac.service.RiskAlerts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "QUERY_FAILED", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) StoreTurnover(c *gin.Context) {
	rows, err := ac.service.ExecuteFunction(c.Request.Context(), "get_store_turnover", []string{
		c.DefaultQuery("store_id", "1"),
		c.DefaultQuery("date_from", "2020-01-01"),
		c.DefaultQuery("date_to", "2100-01-01"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) TransactionStatuses(c *gin.Context) {
	rows, err := ac.service.GetTableRows(c.Request.Context(), "vw_store_transactions", 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "QUERY_FAILED", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) FailedTransactions(c *gin.Context) {
	rows, err := ac.service.ExecuteFunction(c.Request.Context(), "get_failed_transactions_analytics", []string{
		c.DefaultQuery("date_from", "2020-01-01"),
		c.DefaultQuery("date_to", "2100-01-01"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) RPCHealth(c *gin.Context) {
	rows, err := ac.service.GetTableRows(c.Request.Context(), "vw_rpc_node_status", 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "QUERY_FAILED", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (ac *AdminController) CreateDemoInvoice(c *gin.Context) {
	result, err := ac.service.CreateDemoInvoice(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "DEMO_FAILED", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (ac *AdminController) SubmitDemoPayment(c *gin.Context) {
	var req models.DemoPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_INPUT", "message": "invalid request", "details": gin.H{}}})
		return
	}
	result, err := ac.service.SubmitDemoPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "DEMO_FAILED", "message": err.Error(), "details": gin.H{}}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
