package database

import (
	"VkTestMattermostBot/internal/core"

	"github.com/tarantool/go-tarantool/v2"
)


// получение списка id голосований в которых пользователь принял участие
func GetUserVotesDoneCast(userId string) []int{
	var res []int

	req := tarantool.NewSelectRequest(tableNames[3]).Index("primary").Key([]interface{}{userId})
	resp, _ := DbConnection.Do(req).Get()

	if resp == nil{
		return []int{}
	}

	core.AppLogger.Println(resp)

	return res
}


// добавление у пользователя нового голосования, в котором он поучаствовал
func AddInUserListNewVoteDoneCast(userId string, voteId int) bool{
	reqGet := tarantool.NewSelectRequest(tableNames[3]).Index("primary").Key([]interface{}{userId})
	respGet, _ := DbConnection.Do(reqGet).Get()
	if respGet == nil{
		req := tarantool.NewInsertRequest(tableNames[3]).Tuple([]interface{}{
			userId,
			[]int{voteId},
		})
		resp, _ := DbConnection.Do(req).Get()
		if resp == nil {
			return false
		}
	} else {
		newVoteList := []int{}
		if len(respGet) == 0{

		} else{
			newVoteListInterface := respGet[0].([]interface{})[1].([]interface{})
			for _, elem := range newVoteListInterface{
				newVoteList = append(newVoteList, int(elem.(uint64)))
			}
		}
		
		for _, elem := range newVoteList{
			if elem == voteId{
				return false
			}
		}

		newVoteList = append(newVoteList, voteId)
		req := tarantool.NewReplaceRequest(tableNames[3]).Tuple([]interface{}{
			userId,
			newVoteList,
		})
		resp, _ := DbConnection.Do(req).Get()
		if resp == nil{
			return false
		}
	}

	return true
}