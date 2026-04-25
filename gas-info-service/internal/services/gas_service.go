package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gas-info-service/internal/models"

	"github.com/go-redis/redis/v8"
)

// GasService описывает бизнес-логику для получения информации о газе.
type GasService interface {
	GetGasInfo(network string) (models.GasInfo, error)
	GetGasHistory(network string) ([]models.GasHistoryPoint, error)
}

type gasService struct {
	redisClient *redis.Client
	ctx         context.Context
	provider    GasProvider
}

func NewGasService() GasService {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return &gasService{
		redisClient: rdb,
		ctx:         context.Background(),
		provider:    newGasProvider(),
	}
}

func (s *gasService) GetGasHistory(network string) ([]models.GasHistoryPoint, error) {
	cacheKey := fmt.Sprintf("gasinfo:history:%s", network)
	values, err := s.redisClient.LRange(s.ctx, cacheKey, 0, 23).Result()
	if err == nil && len(values) > 0 {
		points := make([]models.GasHistoryPoint, 0, len(values))
		for i := len(values) - 1; i >= 0; i-- {
			var point models.GasHistoryPoint
			if err := json.Unmarshal([]byte(values[i]), &point); err != nil {
				return nil, err
			}
			points = append(points, point)
		}
		return points, nil
	}

	points := make([]models.GasHistoryPoint, 0, 12)
	now := time.Now().UTC()
	for i := 11; i >= 0; i-- {
		points = append(points, models.GasHistoryPoint{
			Network:       network,
			GasPrice:      40 + float64((12-i)%5)*3.5,
			EstimatedTime: 20 + ((12 - i) % 4 * 5),
			NetworkStatus: "normal",
			CapturedAt:    now.Add(-time.Duration(i) * time.Hour).Format(time.RFC3339),
		})
	}
	return points, nil
}

func (s *gasService) GetGasInfo(network string) (models.GasInfo, error) {
	cacheKey := fmt.Sprintf("gasinfo:%s", network)
	val, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	var gasInfo models.GasInfo
	if err != nil {
		gasInfo, err = s.provider.GetGasInfo(s.ctx, network)
		if err != nil {
			return models.GasInfo{}, err
		}
		serialized, err := json.Marshal(gasInfo)
		if err != nil {
			return models.GasInfo{}, err
		}
		_ = s.redisClient.Set(s.ctx, cacheKey, serialized, 30*time.Second).Err()
		s.recordHistory(network, gasInfo)
	} else {
		if err := json.Unmarshal([]byte(val), &gasInfo); err != nil {
			gasInfo, err = s.provider.GetGasInfo(s.ctx, network)
			if err != nil {
				return models.GasInfo{}, err
			}
		}
	}
	return gasInfo, nil
}

func (s *gasService) recordHistory(network string, gasInfo models.GasInfo) {
	point := models.GasHistoryPoint{
		Network:       network,
		GasPrice:      gasInfo.GasPrice,
		EstimatedTime: gasInfo.EstimatedTime,
		NetworkStatus: gasInfo.NetworkStatus,
		CapturedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	payload, err := json.Marshal(point)
	if err != nil {
		return
	}
	cacheKey := fmt.Sprintf("gasinfo:history:%s", network)
	_ = s.redisClient.LPush(s.ctx, cacheKey, payload).Err()
	_ = s.redisClient.LTrim(s.ctx, cacheKey, 0, 23).Err()
	_ = s.redisClient.Expire(s.ctx, cacheKey, 24*time.Hour).Err()
}

type GasProvider interface {
	GetGasInfo(ctx context.Context, network string) (models.GasInfo, error)
}

func newGasProvider() GasProvider {
	if baseURL := os.Getenv("GAS_PROVIDER_URL"); baseURL != "" {
		return &httpGasProvider{
			baseURL: baseURL,
			client:  &http.Client{Timeout: 5 * time.Second},
		}
	}
	return mockGasProvider{}
}

type mockGasProvider struct{}

func (mockGasProvider) GetGasInfo(ctx context.Context, network string) (models.GasInfo, error) {
	return models.GasInfo{
		Network:       network,
		GasPrice:      50.0,
		EstimatedTime: 30,
		NetworkStatus: "normal",
	}, nil
}

type httpGasProvider struct {
	baseURL string
	client  *http.Client
}

func (p *httpGasProvider) GetGasInfo(ctx context.Context, network string) (models.GasInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL, nil)
	if err != nil {
		return models.GasInfo{}, err
	}
	query := req.URL.Query()
	query.Set("network", network)
	req.URL.RawQuery = query.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return models.GasInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return models.GasInfo{}, fmt.Errorf("gas provider returned status %d", resp.StatusCode)
	}

	var gasInfo models.GasInfo
	if err := json.NewDecoder(resp.Body).Decode(&gasInfo); err != nil {
		return models.GasInfo{}, err
	}
	if gasInfo.Network == "" {
		gasInfo.Network = network
	}
	return gasInfo, nil
}
