# Этап сборки: используем официальный образ Go
FROM golang:1.24.1 as builder
WORKDIR /app
# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download
# Копируем весь исходный код
COPY . .
# Собираем бинарный файл для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-terminal-service .

# Финальный образ на основе Alpine Linux
FROM alpine:latest
WORKDIR /root/
# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/auth-terminal-service .
# Копируем файл конфигурации (если нужен)
COPY .env .
# Открываем порт, указанный в .env (8080)
EXPOSE 8080
CMD ["./auth-terminal-service"]
