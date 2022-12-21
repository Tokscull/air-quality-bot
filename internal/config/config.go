package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	TgBot        *TgBotConfig
	Postgresql   *PostgresqlConfig
	WAQIApiToken string
}

type TgBotConfig struct {
	Token string
	Url   string
	Port  string
}

type PostgresqlConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func GetConfig(isProd *bool) (*Config, error) {

	if !*isProd {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file")
		}
	}

	c := &Config{
		TgBot: &TgBotConfig{
			Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
			Url:   os.Getenv("TELEGRAM_BOT_WEBHOOK_URL"),
			Port:  os.Getenv("TELEGRAM_BOT_WEBHOOK_PORT"),
		},
		Postgresql: &PostgresqlConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DB"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		WAQIApiToken: os.Getenv("WAQI_API_TOKEN"),
	}

	return c, nil
}
