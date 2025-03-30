package database

import()


// модель, описывающая голосование
type VoteModel struct{
	Id int;
	Name string;  // название голосования
	Description string;  // описание голосования(опционально)
	Variants map[string][]string;  // название варианта: список id пользователей проголосовавших за этот вариант
	IsActive bool;  // true - голосование продолжается, false - голосование завершено
	ChanelId string;  // id канала mattermost, в котором создано данное голосование
	CreatorId string; // id пользователя mattermost, создавшего голосование - только он может менять название, описание, варианты ответа
	OneAnswerOpinion bool;  // true - голосование с одним вариантом ответа, false - голосование с несколькими вариантами ответа
	IsFillingFinished bool;  // закончено ли наполнение голосования контентом -> пользователи могут отдавать свои голоса
}


// модель описывающая пользователей участвующих в голосованиях
type UserModel struct{
	Id int;
	MattermostId string;  // id пользователя в Mattermost
	Username string;  // логин пользователя
	VotesInfo map[string]string;  // название голосования: вариант за который пользователь отдал свой голос
}


// модель для описания каналов, в которых могут быть голосования
type ChanelModel struct{
	ChanelId string;  // id канала в mattermost
	VotesList []int;  // список id голосований данного канала
}

// модель для описания варианта ответа в конкретном голосовании
type VoteVariantModel struct{
	Id int;  // id записи
	VoteId int; // id голосования
	VariantName string;  // название варианта
	UsersIdsCastVariant []string;  // список id пользователей которые проголосовали за данный вариант
}