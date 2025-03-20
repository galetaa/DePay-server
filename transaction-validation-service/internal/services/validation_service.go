package services

import (
	"errors"

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
	// Имитируем базовую проверку адресов: они должны иметь определённую длину (например, 42 символа для Ethereum)
	if len(req.SenderAddress) != 42 || len(req.RecipientAddress) != 42 {
		return errors.New("invalid address format")
	}
	// Если всё корректно, возвращаем nil (успешно)
	return nil
}
