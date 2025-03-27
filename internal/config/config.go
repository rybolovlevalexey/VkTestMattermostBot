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
}

func LoadBotConfig() BotConfig{
    godotenv.Load()

	botConfig := BotConfig{
		ServerURL: os.Getenv("SERVER_URL"),
		WebSocketURL: os.Getenv("WEBSOCKET_URL"),
		Token: os.Getenv("BOT_TOKEN"),
		BotUserID: os.Getenv("BOT_USER_ID"),
	}

	return botConfig
}

type LoggerConfig struct{
	LogFilePath string;  // путь до файла, в который сохранять логи
}

func LoadLoggerConfig() LoggerConfig{
    godotenv.Load()

	loggerConfig := LoggerConfig{
		LogFilePath: os.Getenv("LOG_FILE_PATH"),
	}

	return loggerConfig
}
