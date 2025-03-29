package usecases

import (
	"VkTestMattermostBot/internal/database"
)

// создание нового голосования
func CreateVote(userId string, chanelId string) int{
	var voteId = database.AddVote(database.VoteModel{ChanelId: chanelId, CreatorId: userId})

	return voteId 
}


// установить название голосования
func SetVoteName(userId string, voteId int, voteName string) bool{
	vote := database.GetVoteInfoById(voteId)
	
	if vote.CreatorId != userId{
		return false
	}

	database.UpdateVoteName(voteId, voteName)

	return true
}

// установить описание голосования
func SetVoteDesc(userId string, voteId int, voteDesc string) bool{
	vote := database.GetVoteInfoById(voteId)
	
	if vote.CreatorId != userId{
		return false
	}

	database.UpdateVoteDesc(voteId, voteDesc)

	return true
}

// установить варианты ответа голосования
func SetVoteVariants(userId string, voteId int, voteVariants []string) bool{
	vote := database.GetVoteInfoById(voteId)
	
	if vote.CreatorId != userId{
		return false
	}

	database.UpdateVoteVariants(voteId, voteVariants)

	return true
}

// установить голосование с один вариантом ответа или несколькими
func SetVoteIsOneAnswer(userId string, voteId int, isOneAnswerVote bool) bool{
	vote := database.GetVoteInfoById(voteId)
	
	if vote.CreatorId != userId{
		return false
	}

	database.UpdateVoteIsOneAnswer(voteId, isOneAnswerVote)

	return true
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