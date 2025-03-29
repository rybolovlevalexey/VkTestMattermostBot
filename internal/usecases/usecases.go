package usecases

import (
	"VkTestMattermostBot/internal/database"
)

// создание нового голосования
func CreateVote(userId string, chanelId string) int{
	var voteId = database.AddVote(database.VoteModel{ChanelId: chanelId, CreatorId: userId})

	return voteId 
}

// голосование пользователя за определённый вариант в конкретном голосовании
func UserCastVoteByVoteId(userId int, voteId int, variant string){

}

// посмотреть информацию по конкретному голосованию
func ViewCurrentVoteResult(voteId int){

}

// посмотреть все возможные голосования
func ViewAllVotesResults(){

}

// остановка конкретного голосования
func StopCurrentVote(voteId int) bool{
	var resultFlag = false

	return resultFlag
}

// удаление конкретного голосования
func DeleteCurrentVote(voteId int) bool{
	var resultFlag = false

	return resultFlag
}