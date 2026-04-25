package services

import (
	"errors"
	"strings"

	"shared/validation"
	"transaction-validation-service/internal/models"
)

// ValidationService описывает логику валидации транзакций.
type ValidationService interface {
	Validate(req models.ValidationRequest) error
}

type validationService struct{}

func NewValidationService() ValidationService {
	return &validationService{}
}

// Validate имитирует проверку транзакции:
// Проверяет, достаточно ли средств (заглушка возвращает ошибку, если amount равен "0")
// и валидность адресов (например, проверка на длину)
func (s *validationService) Validate(req models.ValidationRequest) error {
	// Имитируем проверку баланса: если сумма "0", считаем, что средств недостаточно
	if req.Amount == "0" {
		return errors.New("insufficient funds")
	}
	if err := validation.PositiveAmount(req.Amount); err != nil {
		return err
	}
	if err := validation.EVMAddress(req.SenderAddress); err != nil {
		return err
	}
	if err := validation.EVMAddress(req.RecipientAddress); err != nil {
		return err
	}
	if strings.TrimSpace(req.SignedData) == "" {
		return errors.New("signature is required")
	}
	// Если всё корректно, возвращаем nil (успешно)
	return nil
}
