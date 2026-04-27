package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config хранит общие настройки приложения
type Config struct {
	AppName     string `mapstructure:"APP_NAME"`
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	RedisAddr   string `mapstructure:"REDIS_ADDR"`
	RedisPass   string `mapstructure:"REDIS_PASSWORD"`
	RabbitMQURL string `mapstructure:"RABBITMQ_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

// Cfg – глобальная переменная для доступа к конфигурации
var Cfg Config

// InitConfig загружает конфигурацию из файла или переменных окружения
func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // можно добавить дополнительные пути
	viper.AutomaticEnv()
	viper.SetDefault("REDIS_ADDR", "redis:6379")
	viper.SetDefault("RABBITMQ_URL", "amqp://myuser:mypassword@rabbitmq:5672/")

	// Если конфигурационный файл не найден, используем только env-переменные
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found, using environment variables")
	}

	err := viper.Unmarshal(&Cfg)
	if err != nil {
		return err
	}
	return nil
}
