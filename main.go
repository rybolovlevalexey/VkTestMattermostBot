package main

import (
	"log"

	"VkTestMattermostBot/internal/bot"
	"VkTestMattermostBot/internal/config"
	"VkTestMattermostBot/internal/core"
	"VkTestMattermostBot/internal/database"

	"github.com/tarantool/go-tarantool"
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
	opts := tarantool.Opts{
		User: dbConfig.User, 
		Pass: dbConfig.Password,
	}
    conn, err := tarantool.Connect(dbConfig.Host + ":" + dbConfig.Port, opts)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
	database.DbConnection = conn
	appLogger.Printf("Подключение к базе данных %s:%s выполнено успешно\n", dbConfig.Host, dbConfig.Port)

	appLogger.Println("Запуск инизиализации базы данных")
	database.InitDataBase(database.ArgsInitDataBase{InitVote: true, InitChanels: true, InitVoteVariants: true, InitUsersTable: true})
	appLogger.Println("Инициализация базы данных выполнена успешно")

	// Создание и запуск бота
	bot.StartBot()
}
