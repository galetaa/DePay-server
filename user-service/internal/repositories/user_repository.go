package repositories

import (
	"errors"
	"user-service/internal/models"
)

// UserRepository описывает методы для работы с данными пользователей
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

// userRepo — пример реализации, в продакшене подключите базу данных (PostgreSQL через GORM или sqlx)
type userRepo struct {
	users map[string]*models.User
}

func NewUserRepository() UserRepository {
	return &userRepo{
		users: make(map[string]*models.User),
	}
}

func (r *userRepo) Create(user *models.User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
