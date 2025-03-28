package main

import (
	"log"

	"VkTestMattermostBot/internal/bot"
	"VkTestMattermostBot/internal/config"
	"VkTestMattermostBot/internal/core"
	"VkTestMattermostBot/internal/database"

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

	// appLogger.Println("Запуск инизиализации базы данных")
	// database.InitDataBase()
	// appLogger.Println("Инициализация базы данных выполнена успешно")


	// тестовое использование методов для работы с БД
	//----------
	// database.AddVote(database.VoteModel{})
	// log.Println(database.GetVotesIds())
	/*
	database.AddVote(database.VoteModel{
		Name: "новое голосование",
		Variants: map[string][]string{"cool": []string{}, "not cool": []string{}, },
		ChanelId: "x123",
	})
	*/
	// log.Println(database.DeleteVote(3))
	// log.Println(database.GetVotesNames())
	log.Println(database.FinishVote(2))
	//----------

	// Создание и запуск бота
	bot.StartBot()
}
