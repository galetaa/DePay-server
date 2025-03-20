package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"os"
	"transaction-core-service/internal/models"
	"transaction-core-service/internal/repositories"
)

// TransactionService описывает бизнес-логику для работы с транзакциями.
type TransactionService interface {
	Initiate(tx models.Transaction) error
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
		// Подключаемся к RabbitMQ (конфигурация должна задаваться через переменные окружения для продакшена)
		conn, err := amqp091.Dial("amqp://myuser:mypassword@rabbitmq:5672/")
		if err != nil {
			// В продакшене можно логировать и пытаться переподключиться, здесь panic для краткости.
			panic(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
		}
		ch, err = conn.Channel()
		if err != nil {
			panic(fmt.Sprintf("failed to open RabbitMQ channel: %v", err))
		}
		q, err = ch.QueueDeclare(
			"transaction_queue", // имя очереди
			true,                // durable
			false,               // auto-delete
			false,               // exclusive
			false,               // no-wait
			nil,                 // arguments
		)
		if err != nil {
			panic(fmt.Sprintf("failed to declare RabbitMQ queue: %v", err))
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
		return errors.New(fmt.Sprintf("failed to publish message: %v", err))
	}
	return nil
}
