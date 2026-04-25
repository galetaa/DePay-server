package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"shared/auth"
	"user-service/internal/models"
)

// UserRepository описывает методы для работы с данными пользователей
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	UpdateProfile(ctx context.Context, id string, req models.UpdateProfileRequest) (*models.User, error)
	SubmitKYC(ctx context.Context, userID string, req models.SubmitKYCRequest) error
	GetKYCStatus(ctx context.Context, userID string) (string, error)
	SaveRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*models.User, error)
}

type userRepo struct {
	users         map[string]*models.User
	usersByID     map[string]*models.User
	refreshTokens map[string]string
	mu            sync.RWMutex
	nextID        int64
}

func NewUserRepository() UserRepository {
	return &userRepo{
		users:         make(map[string]*models.User),
		usersByID:     make(map[string]*models.User),
		refreshTokens: make(map[string]string),
		nextID:        1,
	}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	if user.ID == "" {
		user.ID = strconv.FormatInt(r.nextID, 10)
		r.nextID++
	}
	if user.Username == "" {
		user.Username = user.Email
	}
	if user.KYCStatus == "" {
		user.KYCStatus = "not_submitted"
	}
	if len(user.Roles) == 0 {
		user.Roles = []string{"user"}
	}
	user.CreatedAt = time.Now().UTC()
	r.users[user.Email] = user
	r.usersByID[user.ID] = user
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByID[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepo) UpdateProfile(ctx context.Context, id string, req models.UpdateProfileRequest) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.usersByID[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	return user, nil
}

func (r *userRepo) SubmitKYC(ctx context.Context, userID string, req models.SubmitKYCRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.usersByID[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.KYCStatus = "pending"
	return nil
}

func (r *userRepo) GetKYCStatus(ctx context.Context, userID string) (string, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	return user.KYCStatus, nil
}

func (r *userRepo) SaveRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByID[userID]; !exists {
		return errors.New("user not found")
	}
	r.refreshTokens[tokenHash] = userID
	return nil
}

func (r *userRepo) GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.refreshTokens[tokenHash]
	if !exists {
		return nil, errors.New("refresh token not found")
	}
	user, exists := r.usersByID[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

type postgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, user *models.User) error {
	username := user.Username
	if username == "" {
		username = user.Email
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users(username, email, phone_number, password_hash, first_name, last_name, kyc_status)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, 'not_submitted')
		RETURNING user_id, created_at
	`, username, user.Email, user.PhoneNumber, user.PasswordHash, user.FirstName, user.LastName).Scan(&userID, &user.CreatedAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_roles(user_id, role_id)
		SELECT $1, role_id FROM roles WHERE role_name = 'user'
	`, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	user.ID = strconv.FormatInt(userID, 10)
	user.Username = username
	user.KYCStatus = "not_submitted"
	user.Roles = []string{"user"}
	return nil
}

func (r *postgresUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.getOne(ctx, "u.email = $1", email)
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	return r.getOne(ctx, "u.user_id = $1", id)
}

func (r *postgresUserRepo) getOne(ctx context.Context, where string, arg any) (*models.User, error) {
	query := fmt.Sprintf(`
		SELECT
			u.user_id::text,
			u.username,
			u.email,
			COALESCE(u.phone_number, ''),
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			u.password_hash,
			u.kyc_status::text,
			u.created_at,
			COALESCE(string_agg(r.role_name, ','), '')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		WHERE %s
		GROUP BY u.user_id
	`, where)

	user := &models.User{}
	var rolesCSV string
	if err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.KYCStatus,
		&user.CreatedAt,
		&rolesCSV,
	); err != nil {
		return nil, err
	}
	if rolesCSV != "" {
		user.Roles = strings.Split(rolesCSV, ",")
	}
	return user, nil
}

func (r *postgresUserRepo) UpdateProfile(ctx context.Context, id string, req models.UpdateProfileRequest) (*models.User, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET username = COALESCE(NULLIF($2, ''), username),
		    phone_number = COALESCE(NULLIF($3, ''), phone_number),
		    first_name = COALESCE(NULLIF($4, ''), first_name),
		    last_name = COALESCE(NULLIF($5, ''), last_name),
		    address = COALESCE(NULLIF($6, ''), address),
		    updated_at = now()
		WHERE user_id = $1
	`, id, req.Username, req.PhoneNumber, req.FirstName, req.LastName, req.Address)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *postgresUserRepo) SubmitKYC(ctx context.Context, userID string, req models.SubmitKYCRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var applicationID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO kyc_applications(user_id, status)
		VALUES ($1, 'pending')
		RETURNING kyc_application_id
	`, userID).Scan(&applicationID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO kyc_documents(kyc_application_id, document_type, document_url)
		VALUES ($1, $2, $3)
	`, applicationID, req.DocumentType, req.DocumentURL)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users SET kyc_status = 'pending', updated_at = now() WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *postgresUserRepo) GetKYCStatus(ctx context.Context, userID string) (string, error) {
	var status string
	err := r.db.QueryRowContext(ctx, `SELECT kyc_status::text FROM users WHERE user_id = $1`, userID).Scan(&status)
	return status, err
}

func (r *postgresUserRepo) SaveRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens(user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *postgresUserRepo) GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*models.User, error) {
	var userID string
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id::text
		FROM refresh_tokens
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > now()
		ORDER BY created_at DESC
		LIMIT 1
	`, tokenHash).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, userID)
}

func RefreshTokenHash(token string) string {
	return auth.HashToken(token)
}
