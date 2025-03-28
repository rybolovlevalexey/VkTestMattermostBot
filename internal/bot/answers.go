package bot


var BotAnswers = map[string]string{
	"help": "Этот бот создан для реализации функциональности голосования.\n\n" + 
			"Все сообщения связанные с функциональностью голосования должны начинаться с @название_бота\n" +
			"Все команды связанные с голосованием нужно отправлять в следующем формате: <@название_бота> <команда>\n" +
			"После этого следовать инструкциям из сообщения от бота\n\n" +
			"Список возможных команд для использования функциональности бота:\n" + 
			"1. create - создание нового голосования\n" + 
			"2. cast - отдача своего голоса в конкретном голосовании\n" + 
			"3. check - просмотр результатов конкретного голосования или всех голосований\n" + 
			"4. stop - остановка конкретного голосования\n" + 
			"5. delete - удаление конкретного голосования\n",
	"create": `Вам необходимо закончить наполнение голосования, чтобы она начало работать, для этого используйте следующие команды:
			   1. @название_бота votename <id голосования> <название голосования>
			   2. @название_бота votedesc <id голосования> <описание голосования>
			   3. @название_бота votevariants <id голосования> <варианты ответов через ';'>
			   4. @название_бота voteoneanswer <id голосования> <введите 'Y', если в голосовании может быть только один вариант ответа, введите 'N' в противном случае>`,
	"1": "",
	"2": "",
	"3": "",
}