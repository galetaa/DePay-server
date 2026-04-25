package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"

	"transaction-core-service/internal/models"
	"transaction-core-service/internal/repositories"
)

// TransactionService описывает бизнес-логику для работы с транзакциями.
type TransactionService interface {
	Initiate(tx models.Transaction) error
	Get(transactionID string) (*models.Transaction, error)
	UpdateStatus(transactionID string, status string, failureReason string) error
}

type transactionService struct {
	repo        repositories.TransactionRepository
	amqpChannel *amqp091.Channel
	queue       amqp091.Queue
}

// NewTransactionService создаёт новый экземпляр TransactionService.
// Если переменная окружения SKIP_RABBITMQ установлена в "true", то соединение с RabbitMQ не производится.
func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	var ch *amqp091.Channel
	var q amqp091.Queue

	if os.Getenv("SKIP_RABBITMQ") != "true" {
		rabbitURL := os.Getenv("RABBITMQ_URL")
		if rabbitURL == "" {
			rabbitURL = "amqp://myuser:mypassword@rabbitmq:5672/"
		}
		conn, err := amqp091.Dial(rabbitURL)
		if err != nil {
			ch = nil
			q = amqp091.Queue{Name: "transaction_queue"}
		} else {
			ch, err = conn.Channel()
			if err != nil {
				ch = nil
				q = amqp091.Queue{Name: "transaction_queue"}
			} else {
				q, err = ch.QueueDeclare(
					"transaction_queue",
					true,
					false,
					false,
					false,
					nil,
				)
				if err != nil {
					ch = nil
					q = amqp091.Queue{Name: "transaction_queue"}
				}
			}
		}
	} else {
		// Тестовый режим: пропускаем подключение к RabbitMQ.
		ch = nil
		q = amqp091.Queue{Name: "dummy"}
	}

	return &transactionService{
		repo:        repo,
		amqpChannel: ch,
		queue:       q,
	}
}

func (s *transactionService) Initiate(tx models.Transaction) error {
	// Сохраняем транзакцию в репозитории.
	if err := s.repo.Save(tx); err != nil {
		return err
	}

	// Если RabbitMQ не используется (например, в тестовом режиме), пропускаем публикацию.
	if s.amqpChannel == nil {
		return nil
	}

	// Публикуем сообщение в RabbitMQ.
	body, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	err = s.amqpChannel.Publish(
		"",           // exchange
		s.queue.Name, // routing key (имя очереди)
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (s *transactionService) Get(transactionID string) (*models.Transaction, error) {
	return s.repo.Get(transactionID)
}

func (s *transactionService) UpdateStatus(transactionID string, status string, failureReason string) error {
	return s.repo.UpdateStatus(transactionID, status, failureReason)
}
