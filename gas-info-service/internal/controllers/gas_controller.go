package controllers

import (
	"net/http"

	"gas-info-service/internal/services"

	"github.com/gin-gonic/gin"
)

type GasController struct {
	service services.GasService
}

func NewGasController(s services.GasService) *GasController {
	return &GasController{service: s}
}

// GetGasInfo возвращает информацию о газе для указанного блока сети
func (gc *GasController) GetGasInfo(c *gin.Context) {
	network := c.Query("network")
	if network == "" {
		network = "ethereum" // значение по умолчанию
	}
	gasInfo, err := gc.service.GetGasInfo(network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gasInfo)
}
