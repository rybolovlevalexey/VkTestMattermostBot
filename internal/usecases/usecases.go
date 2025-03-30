package usecases

import (
	// "VkTestMattermostBot/internal/core"
	"VkTestMattermostBot/internal/core"
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
	
	if vote.CreatorId != userId || !vote.IsActive{
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
	
	if vote.CreatorId != userId || !vote.IsActive{
		return false
	}

	database.UpdateVoteIsOneAnswer(voteId, isOneAnswerVote)

	return true
}

// начать голосование
func StartVote(userId string, voteId int) bool{
	vote := database.GetVoteInfoById(voteId)

	if vote.CreatorId != userId || vote.Name == "" || len(vote.Variants) < 2{
		return false
	}

	database.UpdateVoteReadyToStart(voteId)

	return true
}

// голосование пользователя за определённый вариант в конкретном голосовании
// variants - может состоять и из одной строки
func UserCastVoteByVoteId(userId string, voteId int, chanelId string, variants []string) bool{
	vote := database.GetVoteInfoById(voteId)
	// проверка, что есть такое голосование, что оно принадлежит данному каналу и 
	// что голосование запущено (наполнено контентом и не остановлено)
	if vote.Id == -1 || vote.ChanelId != chanelId || !vote.IsFillingFinished || !vote.IsActive{
		return false
	}

	// проверка, что пользователь ещё не проголосовал в этом голосовании
	for _, elem := range database.GetUserVotesDoneCast(userId){
		if elem == voteId{
			return false
		}
	}

	for _, variant := range variants{
		flagDone := database.AddUserCast(voteId, userId, variant)
		if !flagDone{
			return false
		}
	}

	flagDone := database.AddInUserListNewVoteDoneCast(userId, voteId)

	return flagDone
}

// посмотреть информацию по конкретному голосованию
func ViewCurrentVoteResult(voteId int, chanelId string) database.VoteModel{
	flagVoteIdInDB := false
	
	// проверка, что голосование с таким id существует
	for _, elem := range database.GetAllIds(){
		if elem == voteId{
			flagVoteIdInDB = true
			break
		}
	}
	if !flagVoteIdInDB{ // не существует голосования с таким id
		return database.VoteModel{Id: -1,}
	}

	voteResult := database.GetVoteInfoById(voteId)
	core.AppLogger.Println(voteResult)
	// проверка что голосование принадлежит данному каналу, чтобы не показывать голосования созданные в других каналах
	if voteResult.ChanelId != chanelId{
		return database.VoteModel{Id: -2,}
	}

	return voteResult
}

// посмотреть все возможные голосования
func ViewAllVotesResults(chanelId string) []database.VoteModel{
	idsInChanelList := database.GetAllVoteIdsInChanel(chanelId)
	result := []database.VoteModel{} 
	
	for _, elem := range idsInChanelList{
		result = append(result, ViewCurrentVoteResult(elem, chanelId))
	}

	return result
}

// остановка конкретного голосования
func StopCurrentVote(voteId int, chanelId string) bool{
	curVote := database.GetVoteInfoById(voteId)
	if curVote.Id == -1{  // проверка на существование голосования
		return false
	}
	if curVote.ChanelId != chanelId{ // проверка на то, что она из этого канала
		return false
	}
	database.FinishVote(voteId)

	return true
}

// удаление конкретного голосования
func DeleteCurrentVote(voteId int, chanelId string, userId string) bool{
	curVote := database.GetVoteInfoById(voteId)
	if curVote.Id == -1{  // проверка на существование голосования
		return false
	}
	if curVote.ChanelId != chanelId{ // проверка на то, что она из этого канала
		return false
	}
	if curVote.CreatorId != userId{  // проверка, что запрос на удаление делает пользователь
		return false
	}
	// удаление голосования
	database.DeleteVote(voteId)

	// ПОДУМАТЬ: нужно ли сделать следующие функции
	// удаление всех вариантов ответа данного голосования из таблицы vote_variants
	// удаление данного голосования у всех пользователей
	// удаление удаление данного голосования из канала

	return true
}