# Этап сборки: используем официальный образ Go
FROM golang:1.24.1 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

# Финальный образ на основе Alpine Linux
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/user-service .
COPY .env .
EXPOSE 8082
CMD ["./user-service"]
