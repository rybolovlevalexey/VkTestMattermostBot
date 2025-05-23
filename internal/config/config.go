package config

import (
	"os"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	ServerURL   string // URL сервера
	WebSocketURL string // URL для WebSocket
	Token       string  // токен бота
	BotUserID	string // ID бота
	BotUserName string  // никнейм бота
}

func LoadBotConfig() BotConfig{
    godotenv.Load()

	botConfig := BotConfig{
		ServerURL: os.Getenv("SERVER_URL"),
		WebSocketURL: os.Getenv("WEBSOCKET_URL"),
		Token: os.Getenv("BOT_TOKEN"),
	}

	return botConfig
}

type DataBaseConfig struct{
	User string;
	Password string;
	Host string;
	Port string;
}

func LoadDBConfig() DataBaseConfig{
	godotenv.Load()

	dbConfig := DataBaseConfig{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
	}
	return dbConfig
}