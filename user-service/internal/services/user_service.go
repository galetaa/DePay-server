package services

import (
	"context"
	"errors"
	"time"

	"shared/auth"
	"user-service/internal/models"
	"user-service/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// UserService определяет бизнес-логику для работы с пользователями
type UserService interface {
	Register(req models.RegisterRequest) (*models.User, error)
	Login(req models.LoginRequest) (*models.User, error)
	RefreshToken(token string) (string, error)
	GetMe(userID string) (*models.User, error)
	UpdateMe(userID string, req models.UpdateProfileRequest) (*models.User, error)
	SubmitKYC(userID string, req models.SubmitKYCRequest) error
	GetKYCStatus(userID string) (string, error)
	IssueTokenPair(user *models.User) (auth.TokenPair, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(req models.RegisterRequest) (*models.User, error) {
	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		KYCStatus:    "not_submitted",
		Roles:        []string{"user"},
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(context.Background(), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(req models.LoginRequest) (*models.User, error) {
	user, err := s.repo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) RefreshToken(token string) (string, error) {
	user, err := s.repo.GetByRefreshTokenHash(context.Background(), auth.HashToken(token))
	if err != nil {
		return "", err
	}
	pair, err := s.IssueTokenPair(user)
	if err != nil {
		return "", err
	}
	return pair.AccessToken, nil
}

func (s *userService) GetMe(userID string) (*models.User, error) {
	return s.repo.GetByID(context.Background(), userID)
}

func (s *userService) UpdateMe(userID string, req models.UpdateProfileRequest) (*models.User, error) {
	return s.repo.UpdateProfile(context.Background(), userID, req)
}

func (s *userService) SubmitKYC(userID string, req models.SubmitKYCRequest) error {
	return s.repo.SubmitKYC(context.Background(), userID, req)
}

func (s *userService) GetKYCStatus(userID string) (string, error) {
	return s.repo.GetKYCStatus(context.Background(), userID)
}

func (s *userService) IssueTokenPair(user *models.User) (auth.TokenPair, error) {
	roles := user.Roles
	if len(roles) == 0 {
		roles = []string{"user"}
	}
	pair, err := auth.NewTokenPair(user.ID, roles, 24*time.Hour, 30*24*time.Hour)
	if err != nil {
		return auth.TokenPair{}, err
	}
	if err := s.repo.SaveRefreshToken(context.Background(), user.ID, pair.RefreshTokenHash, pair.RefreshExpiresAt); err != nil {
		return auth.TokenPair{}, err
	}
	return pair, nil
}
