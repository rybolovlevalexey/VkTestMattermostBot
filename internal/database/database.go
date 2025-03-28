package database

import(
	"log"

	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool"
)

var DbConnection *tarantool.Connection


// получение списка id всех голосований
func GetVotesIds() []int{
	var ids []int

	// Создаем SelectRequest
	req := tarantool.NewSelectRequest("vote").
		Key([]interface{}{}) // пустой ключ для выбора всех записей

	// Выполняем запрос
	resp, err := DbConnection.Do(req).Get()
	if err != nil {
		log.Fatalf("Select request failed: %v", err)
	}
	//fmt.Println(resp.Data, resp.SQLInfo, resp.Code)

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

	return ids
}


// получение списка названий всех голосований
func GetVotesNames() []string{
	var names []string

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

	return names
}


// получение информации о голосовании по id
func GetVoteInfoById(voteId int) VoteModel{
	var resultVote VoteModel

	return resultVote
}


// получение информации о голосовании по названию
func GetVoteInfoByName(voteName string) VoteModel{
	var resultVote VoteModel

	return resultVote
}


// создание нового голосования
func AddVote(vote VoteModel) int{
	core.AppLogger.Println("Получен запрос на создание нового голосования")
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
		},
	})
	core.AppLogger.Printf("Insert response (id %d)- Code: %d, Data: %v\n", curId, resp.Code, resp.Data)

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

	req := tarantool.NewUpdateRequest("vote").
	Index("primary").
	Key([]interface{}{voteId}).
	Operations(tarantool.NewOperations().Assign(5, false))

	resp, _ := DbConnection.Do(req).Get()
	log.Println(resp.SQLInfo)
	if len(resp.Data) > 0{
		resultFlag = true
	}

	return resultFlag
}


// удаление голосования
func DeleteVote(voteId int) bool{
	var resultFlag bool = false

	req := tarantool.NewDeleteRequest("vote").Index("primary").Key([]interface{}{voteId})

	// Выполняем запрос
	resp, _ := DbConnection.Do(req).Get()
	
	if len(resp.Data) > 0{
		resultFlag = true
	}

	return resultFlag
}


// инициализация базы данных: создание таблицы, задание типов полей, создание первичного индекса
func InitDataBase(){
	// Создадим таблицу vote с информацией о голосованиях
    resp, err := DbConnection.Call("box.schema.space.create", []interface{}{
        "vote",
        map[string]bool{"if_not_exists": true},
	})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)
	
    // Зададим типы полей
    resp, err = DbConnection.Call("box.space.vote:format", [][]map[string]string{
        {
            {"name": "id", "type": "unsigned"},
            {"name": "name", "type": "string"},
            {"name": "description", "type": "string"},
			{"name": "variants", "type": "map"},
			{"name": "is_active", "type": "boolean"},
			{"name": "chanel_id", "type": "string"},
        }})
	if err != nil{
		panic(err)
	}
	log.Println(resp.Data)

    // Создадим первичный индекс
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
}