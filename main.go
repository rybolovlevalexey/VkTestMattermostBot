package main

import (
	"log"

	"VkTestMattermostBot/internal/bot"
	"VkTestMattermostBot/internal/config"
	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool"
	// "github.com/chilts/sid"
)

func main() {
	// Создание логгера
	loggerConfig := config.LoadLoggerConfig()
	appLogger, err := core.LoggerFactory(loggerConfig.LogFilePath)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	appLogger.Println("Успешно загружены конфиги и инициализировван логгер")
	
	// создание подключения к базе данных
	dbConfig := config.LoadDBConfig()
	opts := tarantool.Opts{
		User: dbConfig.User, 
		Pass: dbConfig.Password,
	}
    conn, err := tarantool.Connect(dbConfig.Host + ":" + dbConfig.Port, opts)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
	appLogger.Printf("Подключение к базе данных %s:%s выполнено успешно\n", dbConfig.Host, dbConfig.Port)

	// Создание и запуск бота
	bot.StartBot()
}
