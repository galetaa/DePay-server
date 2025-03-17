package services

import (
	"errors"
	"user-service/internal/models"
	"user-service/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// UserService определяет бизнес-логику для работы с пользователями
type UserService interface {
	Register(req models.RegisterRequest) (*models.User, error)
	Login(req models.LoginRequest) (*models.User, error)
	RefreshToken(token string) (string, error)
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
		ID:           generateID(), // В продакшене используйте надёжный UUID генератор
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(req models.LoginRequest) (*models.User, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) RefreshToken(token string) (string, error) {
	// Здесь должна быть логика валидации и обновления токена
	newToken := token + "_refreshed" // Заглушка для демонстрации
	return newToken, nil
}

func generateID() string {
	// Для продакшена используйте генерацию UUID, например, через github.com/google/uuid
	return "user-id-123"
}
