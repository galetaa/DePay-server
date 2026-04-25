package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"wallet-service/internal/models"
	"wallet-service/internal/repositories"

	"github.com/go-redis/redis/v8"
)

// WalletService описывает бизнес-логику Wallet Service
type WalletService interface {
	ExportWallets() ([]models.Wallet, error)
	GetBalance(req models.BalanceRequest) (models.BalanceResponse, error)
	CreateWallet(req models.CreateWalletRequest) (models.Wallet, error)
	GetWallet(id string) (models.Wallet, error)
	DeleteWallet(id string) error
	GetWalletBalances(walletID string) ([]models.WalletBalance, error)
	SyncWallet(walletID string) (models.BalanceResponse, error)
}

type walletService struct {
	repo            repositories.WalletRepository
	redisClient     *redis.Client
	ctx             context.Context
	balanceProvider BalanceProvider
}

func NewWalletService(repo repositories.WalletRepository) WalletService {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &walletService{
		repo:            repo,
		redisClient:     rdb,
		ctx:             context.Background(),
		balanceProvider: NewBalanceProviderFromEnv(),
	}
}

// ExportWallets возвращает список кошельков
func (s *walletService) ExportWallets() ([]models.Wallet, error) {
	return s.repo.GetAllWallets(s.ctx)
}

// GetBalance возвращает баланс указанного кошелька, используя кэш Redis
func (s *walletService) GetBalance(req models.BalanceRequest) (models.BalanceResponse, error) {
	cacheKey := fmt.Sprintf("balance:%s", req.Address)
	balance, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	if err != nil { // если ошибка (в том числе redis.Nil) – получаем значение из Blockchain Module
		balance, err = s.balanceProvider.GetBalance(s.ctx, req.Address)
		if err != nil {
			return models.BalanceResponse{}, err
		}
		// Пытаемся кэшировать результат, ошибки игнорируем
		_ = s.redisClient.Set(s.ctx, cacheKey, balance, 60*time.Second)
	}

	return models.BalanceResponse{
		Address:    req.Address,
		Blockchain: "ethereum",
		Balance:    balance,
	}, nil
}

func (s *walletService) CreateWallet(req models.CreateWalletRequest) (models.Wallet, error) {
	return s.repo.Create(s.ctx, req)
}

func (s *walletService) GetWallet(id string) (models.Wallet, error) {
	return s.repo.GetByID(s.ctx, id)
}

func (s *walletService) DeleteWallet(id string) error {
	return s.repo.Delete(s.ctx, id)
}

func (s *walletService) GetWalletBalances(walletID string) ([]models.WalletBalance, error) {
	return s.repo.GetBalances(s.ctx, walletID)
}

func (s *walletService) SyncWallet(walletID string) (models.BalanceResponse, error) {
	wallet, err := s.repo.GetByID(s.ctx, walletID)
	if err != nil {
		return models.BalanceResponse{}, err
	}
	resp, err := s.GetBalance(models.BalanceRequest{Address: wallet.Address})
	if err != nil {
		return models.BalanceResponse{}, err
	}
	if err := s.repo.UpsertNativeBalance(s.ctx, walletID, resp.Balance); err != nil {
		return models.BalanceResponse{}, err
	}
	return resp, nil
}
