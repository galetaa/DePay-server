package services

import (
	"time"

	"kyc-service/internal/models"
)

// KYCService описывает бизнес-логику для обработки KYC запроса
type KYCService interface {
	ProcessKYC(req models.KYCRequest) (models.KYCResponse, error)
}

type kycService struct{}

// NewKYCService создаёт новый экземпляр KYCService
func NewKYCService() KYCService {
	return &kycService{}
}

// ProcessKYC имитирует интеграцию с внешним KYC-провайдером
func (s *kycService) ProcessKYC(req models.KYCRequest) (models.KYCResponse, error) {
	// Имитируем задержку для внешнего вызова
	time.Sleep(2 * time.Second)

	// Возвращаем успешный результат проверки
	resp := models.KYCResponse{
		UserID:    req.UserID,
		KYCStatus: "verified",
		Message:   "KYC verification successful",
	}
	return resp, nil
}
