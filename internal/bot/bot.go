package bot

import (
	"log"
	"strings"
	"encoding/json"
	"strconv"

	"github.com/mattermost/mattermost-server/v6/model"
	
	"VkTestMattermostBot/internal/config"
	"VkTestMattermostBot/internal/database"
	"VkTestMattermostBot/internal/usecases"
)

type MattermostBot struct{
	Client	*model.Client4;
	WSclient 	*model.WebSocketClient;
	BotConfig	config.BotConfig;
}

type InfoToGenerateResponse struct{
	voteId int;
	chanelId string;
	creatorId string;
}


// Загрузка конфигов бота, инициализация и запуск event loop после подключения к web socket
func StartBot() {
	botConfig := config.LoadBotConfig()

	// Инициализация клиента
	client := model.NewAPIv4Client(botConfig.ServerURL)
	client.SetToken(botConfig.Token)

	// Проверка подключения
	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}

	botConfig.BotUserID = user.Id  // сохранение id бота, чтобы не отвечать на свои же сообщения
	botConfig.BotUserName = user.Username  // сохранение username бота, чтобы понимать, что к нему обращаются 

	// Подключение к WebSocket
	wsClient, err := model.NewWebSocketClient4(botConfig.WebSocketURL, botConfig.Token)
	if err != nil {
		log.Fatalf("Ошибка WebSocket: %v", err)
	}
	defer wsClient.Close()

	wsClient.Listen()

	// Запуск обработчика событий в горутине
	go handleEvents(wsClient, client, botConfig)

	// Бесконечный цикл для поддержания работы программы
	select {}
}

// основной цикл работы бота
func handleEvents(wsClient *model.WebSocketClient, client *model.Client4, botConfig config.BotConfig) {
	for event := range wsClient.EventChannel {
		processEvent(event, client, botConfig)
	}
}

func processEvent(event *model.WebSocketEvent, client *model.Client4, botConfig config.BotConfig) {
	// Обрабатываем только новые сообщения
	if event.EventType() != model.WebsocketEventPosted {
		log.Println("Получено не новое сообщение")
		return
	}

	// Десериализация сообщения
	var post model.Post
	if err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post); err != nil {
		log.Printf("Ошибка разбора сообщения: %v", err)
		return
	}
	
	log.Println("ChannelId ", post.ChannelId)

	if post.UserId == botConfig.BotUserID {
		log.Println("Получено сообщение от самого бота")
		return
	}

	// Обработка команд
	handleCommand(&post, client, botConfig)
}

func handleCommand(post *model.Post, client *model.Client4, botConfig config.BotConfig) {
	// Реагируем только на упоминания бота
	if !strings.Contains(post.Message, "@"+botConfig.BotUserName) {
		log.Println("В сообщении нет упоминания бота")
		return
	}
	
	// запуск необходимых методов для выполнения логики приложения
	flagDoneMainLogic, resultMainLogic, infoGenerateResp := mainLogic(post.Message, botConfig, post.UserId, post.ChannelId)
	log.Println("MainLogic ", flagDoneMainLogic, resultMainLogic)

	// создание сообщения, отвечающего пользователю на его запрос
	reply := &model.Post{
		ChannelId: post.ChannelId,
		Message:   generateResponse(post.Message, botConfig, infoGenerateResp),
	}

	if _, _, err := client.CreatePost(reply); err != nil {
		log.Printf("Ошибка отправки ответа: %v", err)
	}
}

// генерация ответов в зависимости от сообщения пользователя и результатов выполнения логики приложения
func generateResponse(message string, botConfig config.BotConfig, infoGenerateResp InfoToGenerateResponse) string {
	message = strings.TrimSpace(message)

	switch {
		case strings.Contains(message, "help"):  
			// получено сообщение с help
			return BotAnswers["help"]
		case strings.TrimSpace(strings.Replace(strings.Replace(message, botConfig.BotUserName, "", 1), "@", "", 1)) == "":  
			// получено пустое сообщение
			return BotAnswers["help"]
		case strings.Contains(message, "create"):  // получена команда на создание нового голосования
			return "Создано голосование с id - " + strconv.Itoa(infoGenerateResp.voteId) + "\n\n" + BotAnswers["create"]
		default:
			return "" + message
	}
}

// запуск необходимых функций в соответствии с полученным сообщением от пользователя
func mainLogic(message string, botConfig config.BotConfig, userMatterMostId string, chanelId string) (bool, []database.VoteModel, InfoToGenerateResponse){
	var result []database.VoteModel
	log.Println(message, botConfig.BotUserName, userMatterMostId)

	switch {
	case strings.Contains(message, "create"):
		newVoteId := usecases.CreateVote(userMatterMostId, chanelId)
		return true, result, InfoToGenerateResponse{voteId: newVoteId}
	}


	return true, result, InfoToGenerateResponse{}
}
