package main

import (
	"log"

	"vk_back_dev_test/internal/core"
	"vk_back_dev_test/internal/bot"
	"vk_back_dev_test/internal/config"
)

func main() {
	loggerConfig := config.LoadLoggerConfig()

	// Создание логгера
	appLogger, err := core.LoggerFactory(loggerConfig.LogFilePath)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	appLogger.Println("Успешно загружены конфиги и инициализировван логгер")

	bot.StartBot()
}
