package controllers

import (
	"net/http"

	"merchant-service/internal/models"
	"merchant-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type MerchantController struct {
	service services.MerchantService
}

func NewMerchantController(service services.MerchantService) *MerchantController {
	return &MerchantController{service: service}
}

func (mc *MerchantController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	merchant, pair, err := mc.service.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"merchant": merchant, "token": pair.AccessToken, "refresh_token": pair.RefreshToken}})
}

func (mc *MerchantController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	merchant, pair, err := mc.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"merchant": merchant, "token": pair.AccessToken, "refresh_token": pair.RefreshToken}})
}

func (mc *MerchantController) Me(c *gin.Context) {
	merchant, err := mc.service.GetMe(merchantID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": merchant})
}

func (mc *MerchantController) SubmitVerification(c *gin.Context) {
	var req models.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	merchant, err := mc.service.SubmitVerification(merchantID(c), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"data": merchant})
}

func (mc *MerchantController) VerificationStatus(c *gin.Context) {
	merchant, err := mc.service.GetMe(merchantID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"verification_status": merchant.VerificationStatus}})
}

func (mc *MerchantController) CreateInvoice(c *gin.Context) {
	var req models.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	invoice, err := mc.service.CreateInvoice(merchantID(c), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": invoice})
}

func (mc *MerchantController) ListInvoices(c *gin.Context) {
	invoices, err := mc.service.ListInvoices(merchantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": invoices})
}

func (mc *MerchantController) CreateTerminal(c *gin.Context) {
	var req models.CreateTerminalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	terminal, err := mc.service.CreateTerminal(merchantID(c), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": terminal})
}

func (mc *MerchantController) ListTerminals(c *gin.Context) {
	terminals, err := mc.service.ListTerminals(merchantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": terminals})
}

func (mc *MerchantController) CreateWebhook(c *gin.Context) {
	var req models.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	webhook, err := mc.service.CreateWebhook(merchantID(c), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": webhook})
}

func (mc *MerchantController) ListWebhooks(c *gin.Context) {
	webhooks, err := mc.service.ListWebhooks(merchantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": webhooks})
}

func (mc *MerchantController) GetWebhook(c *gin.Context) {
	webhook, err := mc.service.GetWebhook(merchantID(c), c.Param("webhook_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": webhook})
}

func (mc *MerchantController) DeleteWebhook(c *gin.Context) {
	if err := mc.service.DeleteWebhook(merchantID(c), c.Param("webhook_id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (mc *MerchantController) TestWebhook(c *gin.Context) {
	delivery, err := mc.service.TestWebhook(merchantID(c), c.Param("webhook_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"data": delivery})
}

func (mc *MerchantController) CreateAPIKey(c *gin.Context) {
	var req models.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	key, err := mc.service.CreateAPIKey(merchantID(c), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": key})
}

func (mc *MerchantController) ListAPIKeys(c *gin.Context) {
	keys, err := mc.service.ListAPIKeys(merchantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": keys})
}

func (mc *MerchantController) RevokeAPIKey(c *gin.Context) {
	if err := mc.service.RevokeAPIKey(merchantID(c), c.Param("key_id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func merchantID(c *gin.Context) string {
	rawClaims, exists := c.Get("claims")
	if !exists {
		return ""
	}
	claims, ok := rawClaims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return ""
	}
	return sub
}
