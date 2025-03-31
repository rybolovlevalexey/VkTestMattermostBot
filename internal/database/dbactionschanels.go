package database

import (
	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool/v2"
)


// создание нового канала
func AddNewChanel(chanelId string) ChanelModel{
	var result ChanelModel
	core.AppLogger.Printf("Получен запрос на создание нового канала с id %s\n", chanelId)

	req := tarantool.NewInsertRequest(tableNames[1]).Tuple([]interface{}{
		chanelId,
		[]int{},
	})
	resp, _ := DbConnection.Do(req).Get()
	
	if len(resp) == 0{
		result = ChanelModel{ChanelId: "-1", VotesList: []int{}}
	} else {
		result = ChanelModel{ChanelId: chanelId, VotesList: []int{}}
	}

	return result
}

// добавление нового голосования в данный канал
func AddNewVoteInChanel(chanelId string, voteId int) bool {
	core.AppLogger.Printf("Получен запрос на добавление в канал %s голосования %d\n", chanelId, voteId)

	reqSelect := tarantool.NewSelectRequest(tableNames[1]).Index("primary").Key([]interface{}{chanelId})
	resp, _ := DbConnection.Do(reqSelect).Get()
	if resp == nil{
		return false
	}
	if len(resp) == 0{
		return false
	}
	
	votesIdList := resp[0].([]interface{})[1].([]interface{}) // получение исходного списка
	votesIdList = append(votesIdList, voteId)  // добавление нового id
	

	reqUpdate := tarantool.NewUpdateRequest(tableNames[1]).Index("primary").Key([]interface{}{chanelId}).Operations(
		tarantool.NewOperations().Assign(1, votesIdList))
	resp, _ = DbConnection.Do(reqUpdate).Get()
	core.AppLogger.Println(resp)
	return true
}


// проверка зарегистрирован ли данный канал в таблице
func ChanelIdInTable(chanelId string) bool{
	var result bool
	core.AppLogger.Printf("Получен запрос на проверку зарегистрирован ли канал с id %s\n", chanelId)

	req := tarantool.NewSelectRequest(tableNames[1]).Index("primary").Key([]interface{}{chanelId})
	resp, _ := DbConnection.Do(req).Get()

	if resp == nil{
		result = false
	} else {
		if len(resp) == 0{
			result = false
		} else {
			result = true
		}
	}

	return result
}


// получение списка всех id голосований, которые есть в данном канале
func GetAllVoteIdsInChanel(chanelId string) []int{
	req := tarantool.NewSelectRequest(tableNames[1]).Index("primary").Key([]interface{}{chanelId})
	resp, _ := DbConnection.Do(req).Get()

	if resp == nil{
		return []int{}
	}
	if len(resp) == 0{
		return []int{}
	}

	voteIdsList := []int{}
	for _, elem := range resp[0].([]interface{})[1].([]interface{}){
		voteIdsList = append(voteIdsList, int(elem.(uint64)))
	}

	return voteIdsList
}