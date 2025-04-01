package logging

import (
	"go.uber.org/zap"
)

// Logger – глобальный логгер, который можно использовать во всех сервисах
var Logger *zap.Logger

// InitLogger инициализирует логгер в production режиме
func InitLogger() error {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

// Sync завершает работу логгера (важно при завершении приложения)
func Sync() {
	err := Logger.Sync()
	if err != nil {
		return
	}
}
