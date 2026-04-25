package services

import (
	"context"
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
	Broadcast(transactionID string) (*models.Transaction, error)
}

type transactionService struct {
	repo        repositories.TransactionRepository
	amqpChannel *amqp091.Channel
	queue       amqp091.Queue
	broadcaster BlockchainBroadcaster
	dispatcher  WebhookDispatcher
}

// NewTransactionService создаёт новый экземпляр TransactionService.
// Если переменная окружения SKIP_RABBITMQ установлена в "true", то соединение с RabbitMQ не производится.
func NewTransactionService(repo repositories.TransactionRepository, opts ...Option) TransactionService {
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

	svc := &transactionService{
		repo:        repo,
		amqpChannel: ch,
		queue:       q,
		broadcaster: NewBlockchainBroadcasterFromEnv(),
		dispatcher:  noopWebhookDispatcher{},
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

type Option func(*transactionService)

func WithWebhookDispatcher(dispatcher WebhookDispatcher) Option {
	return func(s *transactionService) {
		if dispatcher != nil {
			s.dispatcher = dispatcher
		}
	}
}

func WithBlockchainBroadcaster(broadcaster BlockchainBroadcaster) Option {
	return func(s *transactionService) {
		if broadcaster != nil {
			s.broadcaster = broadcaster
		}
	}
}

func (s *transactionService) Initiate(tx models.Transaction) error {
	// Сохраняем транзакцию в репозитории.
	if err := s.repo.Save(tx); err != nil {
		return err
	}
	s.dispatch("transaction.created", tx)

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
	if err := s.repo.UpdateStatus(transactionID, status, failureReason); err != nil {
		return err
	}
	tx, err := s.repo.Get(transactionID)
	if err != nil {
		return nil
	}
	s.dispatch(eventTypeForStatus(status), *tx)
	return nil
}

func (s *transactionService) Broadcast(transactionID string) (*models.Transaction, error) {
	tx, err := s.repo.Get(transactionID)
	if err != nil {
		return nil, err
	}
	if tx.Status != "validated" {
		return nil, fmt.Errorf("transaction must be validated before broadcast, current status is %s", tx.Status)
	}
	result, err := s.broadcaster.Broadcast(context.Background(), *tx)
	if err != nil {
		_ = s.repo.UpdateStatus(transactionID, "failed", err.Error())
		tx.Status = "failed"
		tx.FailureReason = err.Error()
		s.dispatch("transaction.failed", *tx)
		return tx, err
	}
	if err := s.repo.MarkBroadcasted(transactionID, result.TxHash); err != nil {
		return nil, err
	}
	tx, err = s.repo.Get(transactionID)
	if err != nil {
		return nil, err
	}
	s.dispatch("transaction.broadcasted", *tx)
	return tx, nil
}

func (s *transactionService) dispatch(eventType string, tx models.Transaction) {
	if eventType == "" {
		return
	}
	s.dispatcher.Dispatch(context.Background(), eventType, tx)
}

func eventTypeForStatus(status string) string {
	switch status {
	case "submitted":
		return "transaction.submitted"
	case "validated":
		return "transaction.validated"
	case "broadcasted":
		return "transaction.broadcasted"
	case "confirmed":
		return "transaction.confirmed"
	case "failed":
		return "transaction.failed"
	default:
		return ""
	}
}
