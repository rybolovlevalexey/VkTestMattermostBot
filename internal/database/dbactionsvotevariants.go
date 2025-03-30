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

	insertReq := tarantool.NewInsertRequest(tableNames[2]).
    Tuple([]interface{}{
        voteVariant.Id,
		voteVariant.VoteId,
		voteVariant.VariantName,
		voteVariant.UsersIdsCastVariant,
    })
	
	resp, _ := DbConnection.Do(insertReq).Get()
	core.AppLogger.Println(resp.Data)
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