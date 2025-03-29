package database

import (
	"fmt"

	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool"
)


// получение списка id всех голосований
func GetVotesIds() []int{
	var ids []int
	core.AppLogger.Println("Запрос в БД на получение списка всех Id")

	// Создаем SelectRequest
	req := tarantool.NewSelectRequest(tableNames[0]).
		Key([]interface{}{}) // пустой ключ для выбора всех записей

	// Выполняем запрос
	resp, err := DbConnection.Do(req).Get()
	if err != nil {
		core.AppLogger.Fatalf("Select request failed: %v", err)
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
	req := tarantool.NewSelectRequest(tableNames[0]).
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

	req := tarantool.NewSelectRequest(tableNames[0]).Index("primary").Key([]interface{}{voteId})
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
		CreatorId: resTuple[6].(string),
		OneAnswerOpinion: resTuple[7].(bool),
		IsFillingFinished: resTuple[8].(bool),
	}

	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Id выполнен успешно")

	return resultVote
}


// получение информации о голосовании по названию
func GetVoteInfoByName(voteName string) VoteModel{
	var resultVote VoteModel
	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Name")

	req := tarantool.NewSelectRequest(tableNames[0]).Index("name_index").Key([]interface{}{voteName})
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
		CreatorId: resTuple[6].(string),
		OneAnswerOpinion: resTuple[7].(bool),
		IsFillingFinished: resTuple[8].(bool),
	}

	core.AppLogger.Println("Запрос в БД на получение названий информации о голосовании по Name выполнен успешно")

	return resultVote
}


// создание нового голосования
func AddVote(vote VoteModel) int{
	core.AppLogger.Println("Запрос в БД на создание нового голосования")
	var curId int

	// создание нового канала, для данного голосования, если канал ещё не был добавлен
	resChanelIdInTable := ChanelIdInTable(vote.ChanelId)
	core.AppLogger.Println("resChanelIdInTable ", resChanelIdInTable)

	if !resChanelIdInTable{
		AddNewChanel(vote.ChanelId)
	}

	for _, elem := range GetVotesIds(){
		if elem > curId{
			curId = elem
		}
	}
	
	curId += 1  // автоинкремент поля id
	vote.IsActive = true // изначально любое голосование активно
	vote.IsFillingFinished = false  // но при этом не готово к использованию, потому что не заполнено
	if vote.Variants == nil{
		vote.Variants = make(map[string][]string)
	}

	// вставка нового голосования
	insertReq := tarantool.NewInsertRequest(tableNames[0]).
    Tuple([]interface{}{
        curId,             // id (unsigned)
        vote.Name,         // name (string)
        vote.Description,  // description (string)
        vote.Variants,     // variants (map)
        vote.IsActive,     // is_active (boolean)
        vote.ChanelId,     // chanel_id (string)
        vote.CreatorId,    // creator_id (unsigned/string)
        vote.OneAnswerOpinion, // one_answer_opinion (boolean)
		vote.IsFillingFinished, // is_filling_finished (boolean)
    })
	
	resp, _ := DbConnection.Do(insertReq).Get()
	core.AppLogger.Println(resp.Error)
	core.AppLogger.Println(vote.Variants)
	core.AppLogger.Printf("Insert response (id %d)- Code: %d, Data: %v\n", curId, resp.Code, resp.Data)

	resNewVoteInChanel := AddNewVoteInChanel(vote.ChanelId, curId)
	core.AppLogger.Println("resNewVoteInChanel", resNewVoteInChanel)

	return curId
}


// голосование пользователя за определённый вариант в определённом голосовании
func CastVote(userId int, voteId int, variant string) bool{
	var resultFlag bool = false
	// реализовать метод голосования пользователя за определённый вариант
	return resultFlag
}


// остановка голосования
func FinishVote(voteId int) bool{
	var resultFlag bool = false
	core.AppLogger.Printf("Запрос в БД на завершение голосования id %v\n", voteId)

	req := tarantool.NewUpdateRequest(tableNames[0]).
	Index("primary").
	Key([]interface{}{voteId}).
	Operations(tarantool.NewOperations().Assign(5, false))

	resp, _ := DbConnection.Do(req).Get()
	core.AppLogger.Println(resp.SQLInfo)
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

	req := tarantool.NewDeleteRequest(tableNames[0]).Index("primary").Key([]interface{}{voteId})

	// Выполняем запрос
	resp, _ := DbConnection.Do(req).Get()
	
	if len(resp.Data) > 0{
		core.AppLogger.Printf("Запрос в БД на завершение голосования id %v не выполнен успешно\n", voteId)
		resultFlag = true
	}

	core.AppLogger.Printf("Запрос в БД на удаление голосования id %v выполнен успешно\n", voteId)

	return resultFlag
}


// изменение названия голосования
func UpdateVoteName(voteId int, voteName string){
	core.AppLogger.Printf("Запрос в БД на обновление имени голосования id %v\n", voteId)

	req := tarantool.NewUpdateRequest(tableNames[0]).
	Index("primary").
	Key([]interface{}{voteId}).
	Operations(tarantool.NewOperations().Assign(1, voteName))

	DbConnection.Do(req).Get()
}


// 