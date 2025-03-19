package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config хранит общие настройки приложения
type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    string `mapstructure:"PORT"`
	// Дополнительные переменные: база данных, ключи, URL нод и т.д.
}

// Cfg – глобальная переменная для доступа к конфигурации
var Cfg Config

// InitConfig загружает конфигурацию из файла или переменных окружения
func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // можно добавить дополнительные пути
	viper.AutomaticEnv()

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
