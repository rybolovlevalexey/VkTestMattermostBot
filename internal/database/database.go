package database

import()


// получение списка id всех голосований
func GetVotesIds() []int{
	var ids []int

	return ids
}

// получение списка названий всех голосованийй
func GetVotesNames() []string{
	var names []string

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
	var id int

	return id
}

// голосование пользователя за определённый вариант в определённом голосовании
func CastVote(userId int, voteId int, variant string) bool{
	var resultFlag bool = false

	return resultFlag
}

// остановка голосования
func FinishVote(voteId int) bool{
	var resultFlag bool = false

	return resultFlag
}

// удаление голосования
func DeleteVote(voteId int) bool{
	var resultFlag bool = false

	return resultFlag
}