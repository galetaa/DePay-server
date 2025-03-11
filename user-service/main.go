package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// =====================
// Модели и типы данных
// =====================

// Пользователь (users)
type User struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	FName        string    `json:"fname"`
	SName        string    `json:"sname"`
	PasswordHash string    `json:"-"`
	KYCVerified  bool      `json:"kyc_verified"`
}

// Запрос регистрации
type UserRegistrationRequest struct {
	Email    string `json:"email"`
	FName    string `json:"fname"`
	SName    string `json:"sname"`
	Password string `json:"password"`
}

// Запрос логина
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Ответ при регистрации/логине
type UserResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	FName       string    `json:"fname"`
	SName       string    `json:"sname"`
	KYCVerified bool      `json:"kyc_verified"`
}

// Ответ с JWT
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Запрос для обновления токена
type RefreshTokenRequest struct {
	Token string `json:"token"`
}

// Запрос KYC
type KYCRequest struct {
	UserID       string `json:"user_id"`
	DocumentType string `json:"document_type"`
	DocumentData string `json:"document_data"`
}

// Ответ KYC
type KYCResponse struct {
	UserID    string `json:"user_id"`
	KYCStatus string `json:"kyc_status"`
	Message   string `json:"message"`
}

// Кошелёк (wallets)
type Wallet struct {
	WalletID   uuid.UUID       `json:"wallet_id"`
	UserID     uuid.UUID       `json:"user_id"`
	WalletName string          `json:"wallet_name"`
	Blockchain string          `json:"blockchain"`
	Addresses  json.RawMessage `json:"addresses"` // Храним JSON
	Tokens     json.RawMessage `json:"tokens"`    // Храним JSON
}

// Ответ для экспорта кошельков
type WalletResponse struct {
	WalletID   uuid.UUID       `json:"wallet_id"`
	UserID     uuid.UUID       `json:"user_id"`
	WalletName string          `json:"wallet_name"`
	Blockchain string          `json:"blockchain"`
	Addresses  json.RawMessage `json:"addresses"`
	Tokens     json.RawMessage `json:"tokens"`
}

// Ответ баланса кошелька
type WalletBalanceResponse struct {
	Address    string `json:"address"`
	Blockchain string `json:"blockchain"`
	Balance    string `json:"balance"`
}

// =====================
// Глобальные переменные
// =====================

var (
	db        *sql.DB
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

// =====================
// Инициализация базы данных
// =====================

func initDB() error {
	pgConn := os.Getenv("PG_CONN")
	if pgConn == "" {
		return errors.New("PG_CONN is not set")
	}
	var err error
	db, err = sql.Open("postgres", pgConn)
	if err != nil {
		return err
	}
	// Проверяем соединение
	if err = db.Ping(); err != nil {
		return err
	}

	// Создаем таблицу пользователей
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		user_id UUID PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		fname TEXT NOT NULL,
		sname TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		kyc_verified BOOLEAN DEFAULT FALSE
	);`
	_, err = db.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Создаем таблицу кошельков
	walletsTable := `
	CREATE TABLE IF NOT EXISTS wallets (
		wallet_id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		wallet_name TEXT NOT NULL,
		blockchain TEXT NOT NULL,
		addresses JSONB,
		tokens JSONB,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);`
	_, err = db.Exec(walletsTable)
	if err != nil {
		return fmt.Errorf("failed to create wallets table: %v", err)
	}

	return nil
}

// =====================
// JWT Функции
// =====================

// Генерация JWT для пользователя
func generateJWT(userID uuid.UUID, email string, kycVerified bool) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id":      userID.String(),
		"email":        email,
		"kyc_verified": kycVerified,
		"exp":          expirationTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, expirationTime, err
}

// =====================
// Обработчики HTTP
// =====================

// Регистрация пользователя
func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	// Создаем нового пользователя
	newUserID := uuid.New()
	_, err = db.Exec(
		"INSERT INTO users (user_id, email, fname, sname, password_hash) VALUES ($1, $2, $3, $4, $5)",
		newUserID, req.Email, req.FName, req.SName, string(hashedPassword),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		UserID:      newUserID,
		Email:       req.Email,
		FName:       req.FName,
		SName:       req.SName,
		KYCVerified: false,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Логин пользователя
func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Получаем пользователя по email
	var user User
	err := db.QueryRow("SELECT user_id, email, fname, sname, password_hash, kyc_verified FROM users WHERE email = $1", req.Email).
		Scan(&user.UserID, &user.Email, &user.FName, &user.SName, &user.PasswordHash, &user.KYCVerified)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Сравниваем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Генерируем JWT
	token, exp, err := generateJWT(user.UserID, user.Email, user.KYCVerified)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := AuthResponse{
		Token:     token,
		ExpiresIn: exp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Обновление JWT
func userRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Расшифровываем токен
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	email, _ := claims["email"].(string)
	kycVerified, _ := claims["kyc_verified"].(bool)

	// Генерируем новый токен
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	newToken, exp, err := generateJWT(uid, email, kycVerified)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	resp := AuthResponse{
		Token:     newToken,
		ExpiresIn: exp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// KYC-заглушка — всегда возвращает, что пользователь верифицирован
func userKYCHandler(w http.ResponseWriter, r *http.Request) {
	var req KYCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Можно обновить пользователя в БД, установив kyc_verified = true
	_, err := db.Exec("UPDATE users SET kyc_verified = TRUE WHERE user_id = $1", req.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating KYC status: %v", err), http.StatusInternalServerError)
		return
	}

	resp := KYCResponse{
		UserID:    req.UserID,
		KYCStatus: "verified",
		Message:   "KYC request processed successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Вспомогательная функция для извлечения user_id из JWT в Authorization header
func extractUserIDFromToken(r *http.Request) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return uuid.Nil, errors.New("Authorization header missing")
	}
	var tokenString string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if err != nil {
		return uuid.Nil, errors.New("Invalid Authorization header")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("Invalid token claims")
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}
	return uuid.Parse(userIDStr)
}

// Экспорт кошельков пользователя
func walletsExportHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из JWT
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query("SELECT wallet_id, user_id, wallet_name, blockchain, addresses, tokens FROM wallets WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error fetching wallets", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wallets []WalletResponse
	for rows.Next() {
		var wallet WalletResponse
		if err := rows.Scan(&wallet.WalletID, &wallet.UserID, &wallet.WalletName, &wallet.Blockchain, &wallet.Addresses, &wallet.Tokens); err != nil {
			http.Error(w, "Error scanning wallet", http.StatusInternalServerError)
			return
		}
		wallets = append(wallets, wallet)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallets)
}

// Получение баланса кошелька по адресу
func walletBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из JWT
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	address := vars["address"]

	// Ищем кошелёк, принадлежащий пользователю, содержащий указанный адрес в JSON-поле addresses
	var blockchain string
	var addressesJSON []byte
	err = db.QueryRow("SELECT blockchain, addresses FROM wallets WHERE user_id = $1", userID).Scan(&blockchain, &addressesJSON)
	if err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Предполагаем, что addresses хранится как JSON-объект вида {"addr1": "balance1", ...}
	var addresses map[string]string
	if err := json.Unmarshal(addressesJSON, &addresses); err != nil {
		http.Error(w, "Error parsing addresses", http.StatusInternalServerError)
		return
	}

	balance, exists := addresses[address]
	if !exists {
		http.Error(w, "Address not found", http.StatusNotFound)
		return
	}

	resp := WalletBalanceResponse{
		Address:    address,
		Blockchain: blockchain,
		Balance:    balance,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// =====================
// Main функция
// =====================

func main() {
	// Инициализируем базу данных
	if err := initDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Создаем роутер
	r := mux.NewRouter()
	// Эндпоинты пользователей
	r.HandleFunc("/user/register", userRegisterHandler).Methods("POST")
	r.HandleFunc("/user/login", userLoginHandler).Methods("POST")
	r.HandleFunc("/user/refresh-token", userRefreshTokenHandler).Methods("POST")
	r.HandleFunc("/user/kyc", userKYCHandler).Methods("POST")
	// Эндпоинты кошельков
	r.HandleFunc("/wallets/export", walletsExportHandler).Methods("GET")
	r.HandleFunc("/wallet/{address}/balance", walletBalanceHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("User Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
