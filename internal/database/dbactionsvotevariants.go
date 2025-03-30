package database

import (
	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool"
)

// добавление нового варианта для данного голосования
func AddNewVariant(voteId int, variantName string){
	var ids = GetAllIds()
	curId := 0

	for _, elem := range ids{
		if elem > curId{
			curId = elem
		}
	}
	curId += 1 // автоинкремент
	voteVariant := VoteVariantModel{Id: curId, VoteId: voteId, VariantName: variantName, UsersIdsCastVariant: []string{}}
	core.AppLogger.Println(voteVariant)
	insertReq := tarantool.NewInsertRequest(tableNames[2]).
    Tuple([]interface{}{
        voteVariant.Id,
		voteVariant.VoteId,
		voteVariant.VariantName,
		voteVariant.UsersIdsCastVariant,
    })
	
	resp, _ := DbConnection.Do(insertReq).Get()
	core.AppLogger.Println("AddNewVariant ", resp.Data, resp.Code, resp.SQLInfo)
}


// получение списка всех id, чтобы сделать автоинкермент
func GetAllIds() []int{
	var result []int

	core.AppLogger.Println("Запрос в БД на получение списка всех Id вариантов ответа голосований")

	req := tarantool.NewSelectRequest(tableNames[2]).
		Key([]interface{}{}) // пустой ключ для выбора всех записей
	resp, _ := DbConnection.Do(req).Get()

	if resp.Data == nil{
		return []int{}
	}

	for _, tuple := range resp.Data {
		fields, ok := tuple.([]interface{})
		if !ok || len(fields) == 0 {
			continue
		}
		
		if id, ok := fields[0].(uint64); ok {
			result = append(result, int(id))
		}
	}

	core.AppLogger.Println("Запрос в БД на получение списка всех Id выполнен успешно")
	return result
}


// получение информации о всех вариантах ответа у конкретного голосования
func GetVoteVariant(voteId int) map[string][]string{
	req := tarantool.NewSelectRequest(tableNames[2]).Key([]interface{}{})
	resp, _ := DbConnection.Do(req).Get()
	
	if resp.Data == nil{
		return map[string][]string{}
	}

	res := make(map[string][]string)
	for _, line := range resp.Data{
		tuple := line.([]interface{})
		
		if int(tuple[1].(uint64)) != voteId {
			continue
		}
		variantName := tuple[2].(string)
		usersIds := tuple[3].([]interface{})
		usersIdsStringArray := make([]string, 0)
		for _, elem := range usersIds{
			usersIdsStringArray = append(usersIdsStringArray, elem.(string))
		}
		res[variantName] = usersIdsStringArray
	}

	return res
}


// добавление голоса конкретного пользователя в конкретном голосовании за определённый вариант ответа
func AddUserCast(voteId int, userId string, variantName string) bool{
	var result bool

	getReq := tarantool.NewSelectRequest(tableNames[2]).
    Key([]interface{}{voteId, variantName})

	getResp, _ := DbConnection.Do(getReq).Get()
	core.AppLogger.Println(getResp.Data)
	record := getResp.Data[0].([]interface{})

	// Обновляем массив
	voters := record[3].([]interface{})
	flagUserIdInVoters := false
	for _, elem := range voters{
		if elem == userId{
			flagUserIdInVoters = true
			break
		}
	}
	if !flagUserIdInVoters{
		voters = append(voters, userId)
	}

	// Записываем обратно
	updateReq := tarantool.NewReplaceRequest(tableNames[2]).
		Tuple([]interface{}{
			record[0],
			record[1],
			record[2],
			voters,
		})

	DbConnection.Do(updateReq).Get()

	return result
}