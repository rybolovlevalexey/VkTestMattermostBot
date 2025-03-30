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

	appLogger.Println("Запуск инизиализации базы данных")
	database.InitDataBase(database.ArgsInitDataBase{InitVote: true, InitChanels: true, InitVoteVariants: true, InitUsersTable: true})
	appLogger.Println("Инициализация базы данных выполнена успешно")

	
	// тестовое использование методов для работы с БД
	//----------
	/*
	log.Println(database.GetVotesIds())
	database.AddVote(database.VoteModel{
		Name: "новое крутое голосование",
		Variants: map[string][]string{"very cool": []string{}, "not very cool": []string{}, },
		ChanelId: "x123",
		CreatorId: "p123",
		OneAnswerOpinion: true,
	})
	log.Println(database.DeleteVote(1))
	log.Println(database.GetVotesNames())
	log.Println(database.FinishVote(2))
	log.Println(database.GetVotesIds())

	log.Println(database.GetVoteInfoById(4))
	log.Println(database.GetVoteInfoByName("новое крутое голосование"))
	log.Println(database.GetVotesIds())
	*/

	// database.AddNewVoteInChanel("chanelIDddddddddd", 133)
	// log.Println(database.GetAllIds())
	// database.AddNewVariant(1, "крутой вариант")
	// database.AddNewVariant(1, "не менее крутой вариант")
	// log.Println(database.GetAllIds())
	// database.AddUserCast(1, "alexeyID2", "крутой вариант")
	// log.Println(database.GetVoteVariant(1))
	// core.AppLogger.Println(database.GetAllVoteIdsInChanel("1kwwosu6gbxfh6y7xnmemtsh"))
	//----------

	// Создание и запуск бота
	bot.StartBot()
}
