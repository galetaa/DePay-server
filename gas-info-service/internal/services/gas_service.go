package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"

	"gas-info-service/internal/models"
)

// GasService описывает бизнес-логику для получения информации о газе.
type GasService interface {
	GetGasInfo(network string) (models.GasInfo, error)
}

type gasService struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewGasService() GasService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	return &gasService{
		redisClient: rdb,
		ctx:         context.Background(),
	}
}

func (s *gasService) GetGasInfo(network string) (models.GasInfo, error) {
	cacheKey := fmt.Sprintf("gasinfo:%s", network)
	val, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	var gasInfo models.GasInfo
	if err != nil { // если нет кэша или ошибка, используем заглушку
		gasInfo = getGasInfoFromExternalService(network)
		// кэшируем значение в Redis на 30 секунд
		serialized := fmt.Sprintf("%f|%d|%s", gasInfo.GasPrice, gasInfo.EstimatedTime, gasInfo.NetworkStatus)
		s.redisClient.Set(s.ctx, cacheKey, serialized, 30*time.Second)
	} else {
		// Простейший парсинг закэшированного значения (в продакшене использовать сериализацию в JSON)
		var gasPrice float64
		var estimatedTime int
		var networkStatus string
		_, err := fmt.Sscanf(val, "%f|%d|%s", &gasPrice, &estimatedTime, &networkStatus)
		if err != nil {
			return models.GasInfo{}, err
		}
		gasInfo = models.GasInfo{
			Network:       network,
			GasPrice:      gasPrice,
			EstimatedTime: estimatedTime,
			NetworkStatus: networkStatus,
		}
	}
	return gasInfo, nil
}

// getGasInfoFromExternalService имитирует вызов внешнего API для получения данных о газе.
func getGasInfoFromExternalService(network string) models.GasInfo {
	// Возвращаем фиксированные значения для демонстрации.
	return models.GasInfo{
		Network:       network,
		GasPrice:      50.0,
		EstimatedTime: 30,
		NetworkStatus: "normal",
	}
}
