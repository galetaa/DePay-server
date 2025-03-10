package main

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Чтение секретного ключа из переменной окружения
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// TerminalRegistrationRequest описывает входные данные для регистрации терминала.
type TerminalRegistrationRequest struct {
	SerialNumber string `json:"serial_number"`
	SecretKey    string `json:"secret_key"`
}

// AuthResponse описывает ответ с JWT-токеном и временем его жизни.
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// generateJWT создает JWT, добавляя в payload идентификатор терминала и время истечения.
func generateJWT(terminalID string) (string, int64, error) {
	// Устанавливаем время истечения токена (например, через 24 часа)
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"terminal_id": terminalID,
		"exp":         expirationTime,
	}
	// Создаем токен с использованием алгоритма HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, expirationTime, err
}

// terminalRegisterHandler обрабатывает регистрацию терминала.
// Принимает JSON с serial_number и secret_key, генерирует JWT и возвращает его.
func terminalRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req TerminalRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Здесь можно добавить проверку serial_number и secret_key через базу данных или иной механизм.
	// Для демонстрации будем принимать все запросы.

	token, exp, err := generateJWT(req.SerialNumber)
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

// terminalRefreshTokenHandler обновляет JWT.
// Получает текущий токен, проверяет его валидность, а затем генерирует новый токен.
func terminalRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Расшифровка и валидация существующего токена.
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
	terminalID, ok := claims["terminal_id"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newToken, exp, err := generateJWT(terminalID)
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

func main() {
	// Проверка, что JWT_SECRET задан
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET is not set")
	}

	// Создаем роутер с помощью Gorilla Mux
	r := mux.NewRouter()
	r.HandleFunc("/terminal/register", terminalRegisterHandler).Methods("POST")
	r.HandleFunc("/terminal/refresh-token", terminalRefreshTokenHandler).Methods("POST")

	// Определяем порт, получая его из переменной окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Auth Terminal Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
