package database

import (
	"fmt"
	"log"

	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool"
)

var DbConnection *tarantool.Connection


// получение списка id всех голосований
func GetVotesIds() []int{
	var ids []int
	core.AppLogger.Println("Запрос в БД на получение списка всех Id")

	// Создаем SelectRequest
	req := tarantool.NewSelectRequest("vote").
		Key([]interface{}{}) // пустой ключ для выбора всех записей

	// Выполняем запрос
	resp, err := DbConnection.Do(req).Get()
	if err != nil {
		log.Fatalf("Select request failed: %v", err)
	}

	// Извлекаем ID
	for _, tuple := range resp.Data {
		// Проверяем тип данных
		fields, ok := tuple.([]interface{})
		if !ok || len(fields) == 0 {
			continue
		}
		
		// Первое поле - ID (uint64)
		if id, ok := fields[0].(uint64); ok {
			ids = append(ids, int(id))
		}
	}

	core.AppLogger.Println("Запрос в БД на получение списка всех Id выполнен успешно")

	return ids
}


// получение списка названий всех голосований
func GetVotesNames() []string{
	var names []string
	core.AppLogger.Println("Запрос в БД на получение названий всех голосований")

	// Создаем SelectRequest
	req := tarantool.NewSelectRequest("vote").
		Key([]interface{}{}) // пустой ключ для выбора всех записей

	// Выполняем запрос
	resp, _ := DbConnection.Do(req).Get()

	// Извлекаем поле namr
	for _, tuple := range resp.Data {
		// Проверяем тип данных
		fields, ok := tuple.([]interface{})
		if !ok || len(fields) == 0 {
			continue
		}
		
		// Второе поле - name (string)
		if name, ok := fields[1].(string); ok {
			names = append(names, name)
		}
	}

	core.AppLogger.Println("Запрос в БД на получение названий всех голосований выполнен успешно")

	return names
}


// получение информации о голосовании по id
func GetVoteInfoById(voteId int) VoteModel{
	var resultVote VoteModel
	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Id")

	req := tarantool.NewSelectRequest("vote").Index("primary").Key([]interface{}{voteId})
	resp, _ := DbConnection.Do(req).Get()
	
	if len(resp.Data) == 0{
		return VoteModel{Id: -1,}
	}

	resTuple := resp.Data[0].([]interface{})
	variants := make(map[string][]string)

	if vars, ok := resTuple[3].(map[interface{}][]interface{}); ok {
		for key, val := range vars{
			strKey := fmt.Sprintf("%v", key)
			strValues := make([]string, len(val))

			for i, v := range val{
				strValues[i] = fmt.Sprintf("%v", v)
			}
			
			variants[strKey] = strValues
		}
	}

	resultVote = VoteModel{
		Id: int(resTuple[0].(uint64)),
		Name: resTuple[1].(string),
		Description: resTuple[2].(string),
		Variants: variants,
		IsActive: resTuple[4].(bool),
		ChanelId: resTuple[5].(string),
	}

	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Id выполнен успешно")

	return resultVote
}


// получение информации о голосовании по названию
func GetVoteInfoByName(voteName string) VoteModel{
	var resultVote VoteModel
	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Name")

	req := tarantool.NewSelectRequest("vote").Index("name_index").Key([]interface{}{voteName})
	resp, _ := DbConnection.Do(req).Get()
	
	if len(resp.Data) == 0{
		return VoteModel{Id: -1,}
	}

	resTuple := resp.Data[0].([]interface{})
	variants := make(map[string][]string)

	if vars, ok := resTuple[3].(map[interface{}][]interface{}); ok {
		for key, val := range vars{
			strKey := fmt.Sprintf("%v", key)
			strValues := make([]string, len(val))

			for i, v := range val{
				strValues[i] = fmt.Sprintf("%v", v)
			}
			
			variants[strKey] = strValues
		}
	}

	resultVote = VoteModel{
		Id: int(resTuple[0].(uint64)),
		Name: resTuple[1].(string),
		Description: resTuple[2].(string),
		Variants: variants,
		IsActive: resTuple[4].(bool),
		ChanelId: resTuple[5].(string),
	}

	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Name выполнен успешно")

	return resultVote
}


// создание нового голосования
func AddVote(vote VoteModel) int{
	core.AppLogger.Println("Запрос в БД на создание нового голосования")
	var curId int

	for _, elem := range GetVotesIds(){
		if elem > curId{
			curId = elem
		}
	}
	
	curId += 1  // автоинкремент поля id
	vote.IsActive = true // изначально любое голосование активно

	// Вставка нового голосования в базу данных
	resp, _ := DbConnection.Call("box.space.vote:insert", []interface{}{
		[]interface{}{
			curId, // id
			vote.Name, // name
			vote.Description, // description
			vote.Variants, // variants
			vote.IsActive, // is_active
			vote.ChanelId, // chanel_id
			vote.CreatorId,
			vote.OneAnswerOpinion,
		},
	})
	core.AppLogger.Printf("Insert response (id %d)- Code: %d, Data: %v\n", curId, resp.Code, resp.Data)

	core.AppLogger.Printf("Запрос в БД на создание нового голосования id %v выполнен успешно\n", curId)

	return curId
}


// голосование пользователя за определённый вариант в определённом голосовании
func CastVote(userId int, voteId int, variant string) bool{
	var resultFlag bool = false

	return resultFlag
}


// остановка голосования
func FinishVote(voteId int) bool{
	var resultFlag bool = false
	core.AppLogger.Printf("Запрос в БД на завершение голосования id %v\n", voteId)

	req := tarantool.NewUpdateRequest("vote").
	Index("primary").
	Key([]interface{}{voteId}).
	Operations(tarantool.NewOperations().Assign(5, false))

	resp, _ := DbConnection.Do(req).Get()
	log.Println(resp.SQLInfo)
	if len(resp.Data) > 0{
		core.AppLogger.Printf("Запрос в БД на завершение голосования id %v не выполнен успешно\n", voteId)
		resultFlag = true
	}

	core.AppLogger.Printf("Запрос в БД на завершение голосования id %v выполнен успешно\n", voteId)

	return resultFlag
}


// удаление голосования
func DeleteVote(voteId int) bool{
	var resultFlag bool = false
	core.AppLogger.Printf("Запрос в БД на завершение голосования id %v\n", voteId)

	req := tarantool.NewDeleteRequest("vote").Index("primary").Key([]interface{}{voteId})

	// Выполняем запрос
	resp, _ := DbConnection.Do(req).Get()
	
	if len(resp.Data) > 0{
		core.AppLogger.Printf("Запрос в БД на завершение голосования id %v не выполнен успешно\n", voteId)
		resultFlag = true
	}

	core.AppLogger.Printf("Запрос в БД на удаление голосования id %v выполнен успешно\n", voteId)

	return resultFlag
}


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
        "vote",
        map[string]bool{"if_not_exists": true},
	})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	// Зададим типы полей
	core.AppLogger.Println("Запрос в БД на определение типов полей в таблице vote")
    resp, err = DbConnection.Call("box.space.vote:format", [][]map[string]string{
        {
            {"name": "id", "type": "unsigned"},
            {"name": "name", "type": "string"},
            {"name": "description", "type": "string"},
			{"name": "variants", "type": "map"},
			{"name": "is_active", "type": "boolean"},
			{"name": "chanel_id", "type": "string"},
			{"name": "creator_id", "type": "string"},
			{"name": "one_answer_opinion", "type": "boolean"},
        }})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	// Создадим необходиые индексы
	core.AppLogger.Println("Запрос в БД на создание первичного индекса по полю id")
    resp, err = DbConnection.Call("box.space.vote:create_index", []interface{}{
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
    resp, _ = DbConnection.Call("box.space.vote:create_index", []interface{}{
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
        "chanels",
        map[string]bool{"if_not_exists": true},
	})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)
	
    // определение полей
	core.AppLogger.Println("Запрос в БД на определение типов полей в таблице chanels")
    resp, err = DbConnection.Call("box.space.chanels:format", [][]map[string]string{
        {
            {"name": "chanel_id", "type": "string"},  // id канала в mattermost
            {"name": "votes_list", "type": "array"},  // список id голосований в канале
            {"name": "creating_vote_now", "type": "boolean"},  // флаг - ведётся ли сейчас создание какого-либо голосования
        }})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

	// создание индексов
	core.AppLogger.Println("Запрос в БД на создание первичного индекса таблицы chanels по полю chanel_id")
    resp, err = DbConnection.Call("box.space.chanels:create_index", []interface{}{
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