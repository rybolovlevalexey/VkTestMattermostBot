package main

import (
	"log"
	"context"
	"time"

	"VkTestMattermostBot/internal/bot"
	"VkTestMattermostBot/internal/config"
	"VkTestMattermostBot/internal/core"
	"VkTestMattermostBot/internal/database"

	// "github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/v2"
)

func main() {
	// Создание логгера
	appLogger, err := core.LoggerFactory()
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	appLogger.Println("Успешно инициализировван логгер")
	core.AppLogger = appLogger
	
	// создание подключения к базе данных
	dbConfig := config.LoadDBConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  dbConfig.Host + ":" + dbConfig.Port,
		User:     dbConfig.User,
		Password: dbConfig.Password,
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}
    conn, err := tarantool.Connect(ctx, dialer, opts)
    if err != nil {
		log.Println("Connection refused:", err)
		return
	}

	database.DbConnection = conn
	appLogger.Printf("Подключение к базе данных %s:%s выполнено успешно\n", dbConfig.Host, dbConfig.Port)

	appLogger.Println("Запуск инизиализации базы данных")
	database.InitDataBase(database.ArgsInitDataBase{InitVote: true, InitChanels: true, InitVoteVariants: true, InitUsersTable: true})
	appLogger.Println("Инициализация базы данных выполнена успешно")

	// Создание и запуск бота
	bot.StartBot()
}
