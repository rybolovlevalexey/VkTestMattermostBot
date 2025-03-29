package database

import()


type VoteModel struct{
	Id int;
	Name string;  // название голосования
	Description string;  // описание голосования(опционально)
	Variants map[string][]string;  // название варианта: список id пользователей проголосовавших за этот вариант
	IsActive bool;  // true - голосование продолжается, false - голосование завершено
	ChanelId string;  // id канала mattermost, в котором создано данное голосование
}

type UserModel struct{
	Id int;
	MattermostId string;  // id пользователя в Mattermost
	Username string;  // логин пользователя
	VotesInfo map[string]string;  // название голосования: вариант за который пользователь отдал свой голос
}
