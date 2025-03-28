package main

import (
	"log"

	"vk_back_dev_test/internal/core"
	"vk_back_dev_test/internal/bot"
	"vk_back_dev_test/internal/config"

	// "github.com/tarantool/go-tarantool"
)

func main() {
	// Создание логгера
	loggerConfig := config.LoadLoggerConfig()
	appLogger, err := core.LoggerFactory(loggerConfig.LogFilePath)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	appLogger.Println("Успешно загружены конфиги и инициализировван логгер")
	

	//appLogger.Printf("Подключение к базе данных %s:%s выполнено успешно\n", dbConfig.Host, dbConfig.Port)

	// Создание и запуск бота
	bot.StartBot()
}
