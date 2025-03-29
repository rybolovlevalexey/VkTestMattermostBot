package database

import (
	"fmt"
	"log"

	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool"
)

var DbConnection *tarantool.Connection
var tableNames = []string{"vote", "chanels"}


// структура для описание, какие таблицы нужно инициализировать, а какие нет
type ArgsInitDataBase struct{
	InitVote bool;
	InitChanels bool;
}

// инициализация базы данных: создание таблицы, задание типов полей, создание первичного индекса
func InitDataBase(args ArgsInitDataBase){
	if args.InitVote{
		initVoteTable()
	}

	if args.InitChanels{
		initChanelsTable()
	}
}


// инициализация таблицы vote
func initVoteTable(){
	// Создадим таблицу vote с информацией о голосованиях
	core.AppLogger.Println("Запрос в БД на создание таблицы vote, если она ещё не существует")
    resp, err := DbConnection.Call("box.schema.space.create", []interface{}{
        tableNames[0],
        map[string]bool{"if_not_exists": true},
	})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	// Зададим типы полей
	core.AppLogger.Println("Запрос в БД на определение типов полей в таблице vote")
    resp, err = DbConnection.Call(fmt.Sprintf("box.space.%s:format", tableNames[0]), 
	[][]map[string]string{
        {
            {"name": "id", "type": "unsigned"},
            {"name": "name", "type": "string"},
            {"name": "description", "type": "string"},
			{"name": "variants", "type": "map"},
			{"name": "is_active", "type": "boolean"},
			{"name": "chanel_id", "type": "string"},
			{"name": "creator_id", "type": "string"},
			{"name": "one_answer_opinion", "type": "boolean"},
			{"name": "is_filling_finished", "type": "boolean"},
        }})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	// Создадим необходиые индексы
	core.AppLogger.Println("Запрос в БД на создание первичного индекса по полю id")
    resp, err = DbConnection.Call(fmt.Sprintf("box.space.%s:create_index", tableNames[0]), 
	[]interface{}{
        "primary",
        map[string]interface{}{
            "parts":         []string{"id"},
            "if_not_exists": true,
	}})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	core.AppLogger.Println("Запрос в БД на создание индекса по полю Name")
    resp, _ = DbConnection.Call(fmt.Sprintf("box.space.%s:create_index", tableNames[0]),
	[]interface{}{
        "name_index",
        map[string]interface{}{
            "parts":         []string{"name"},
            "if_not_exists": true,
			"unique": false,
	}})
	log.Println(resp.Data)
}


// инициализация таблицы chanels
func initChanelsTable(){
	// создание таблицы
	core.AppLogger.Println("Запрос в БД на создание таблицы chanels, если она ещё не существует")
    resp, err := DbConnection.Call("box.schema.space.create", []interface{}{
        tableNames[1],
        map[string]bool{"if_not_exists": true},
	})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)
	
    // определение полей
	core.AppLogger.Println("Запрос в БД на определение типов полей в таблице chanels")
    resp, err = DbConnection.Call(fmt.Sprintf("box.space.%s:format", tableNames[1]), 
	[][]map[string]string{
        {
            {"name": "chanel_id", "type": "string"},  // id канала в mattermost
            {"name": "votes_list", "type": "array"},  // список id голосований в канале
        }})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	// создание индексов
	core.AppLogger.Println("Запрос в БД на создание первичного индекса таблицы chanels по полю chanel_id")
    resp, err = DbConnection.Call(fmt.Sprintf("box.space.%s:create_index", tableNames[1]), []interface{}{
        "primary",
        map[string]interface{}{
            "parts":         []string{"chanel_id"},
            "if_not_exists": true,
	}})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)
}