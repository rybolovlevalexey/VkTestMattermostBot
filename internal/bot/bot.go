package bot

import (
	"log"
	"strings"
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
	
	"vk_back_dev_test/internal/config"
)

type MattermostBot struct{
	Client	*model.Client4;
	WSclient 	*model.WebSocketClient;
	BotConfig	config.BotConfig;
}

/*
func CreateBot() MattermostBot{
	// Конфигурация бота
	botСonfig := config.LoadBotConfig()

	// клиент Mattermost
	client := model.NewAPIv4Client(botСonfig.ServerURL)
	client.SetToken(botСonfig.Token)

	// Проверяем подключение
	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	botСonfig.BotUserID = user.Id
	log.Printf("Бот запущен как @%s", user.Username)

	// подключение к WevSocket
	wsClient, err := model.NewWebSocketClient4(botСonfig.WebSocketURL, botСonfig.Token)
	if err != nil {
		log.Fatalf("Ошибка WebSocket: %v", err)
	}
	defer wsClient.Close()

	stopChan := make(chan bool) // поток для остановки горутины

	go func(){
		wsClient.Listen()
		log.Println("WebSocket подключен. Ожидание сообщений")
	}()
	
	return MattermostBot{
		Client: client,
		WSclient: wsClient,
		BotConfig: botСonfig,
		stopChan: stopChan,
	}
}

func (mattermostBot MattermostBot) StartEventLoop(){
	// Обрабатываем события
	for event := range mattermostBot.WSclient.EventChannel {
		if event.EventType() != model.WebsocketEventPosted {
			continue
		}

		post := &model.Post{}
		err := json.Unmarshal([]byte(event.GetData()["post"].(string)), post)
		if err != nil {
			log.Printf("Ошибка разбора поста: %v", err)
			continue
		}

		if post.UserId == mattermostBot.BotConfig.BotUserID {
			continue
		}

		if strings.Contains(post.Message, "@my_go_bot") {
			reply := &model.Post{
				ChannelId: post.ChannelId,
				Message:   "Привет! Я получил ваше сообщение: _" + post.Message + "_",
			}

			if _, _, err := mattermostBot.Client.CreatePost(reply); err != nil {
				log.Printf("Ошибка отправки: %v", err)
			} else {
				log.Printf("Ответил на сообщение в канале %s", post.ChannelId)
			}
		}
	}
}

func (mattermost MattermostBot) StopEventLoop(){
	close(mattermost.stopChan)
	mattermost.WSclient.Close()
}
*/

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
	botConfig.BotUserID = user.Id
	log.Printf("Бот запущен как @%s", user.Username)

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

func handleEvents(wsClient *model.WebSocketClient, client *model.Client4, botConfig config.BotConfig) {
	fmt.Println("here1")
	for event := range wsClient.EventChannel {
		processEvent(event, client, botConfig)
	}
}

func processEvent(event *model.WebSocketEvent, client *model.Client4, botConfig config.BotConfig) {
	fmt.Println("here2")
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

	// Игнорируем сообщения от самого бота
	if post.UserId == botConfig.BotUserID {
		log.Println("Получено сообщение от самого бота")
		return
	}

	// Обработка команд
	handleCommand(&post, client, botConfig)
}

func handleCommand(post *model.Post, client *model.Client4, botConfig config.BotConfig) {
	fmt.Println("here3")
	// Реагируем только на упоминания бота
	
	if !strings.Contains(post.Message, "@"+botConfig.BotUserID) {
		// TODO: вот это условие поправить надо
		log.Println("В сообщении нет упоминания бота")
		return
	}
	

	reply := &model.Post{
		ChannelId: post.ChannelId,
		Message:   generateResponse(post.Message),
	}

	if _, _, err := client.CreatePost(reply); err != nil {
		log.Printf("Ошибка отправки ответа: %v", err)
	}
}

func generateResponse(message string) string {
	// Здесь можно реализовать логику обработки команд
	switch {
	case strings.Contains(message, "привет"):
		return "Привет! Чем могу помочь?"
	case strings.Contains(message, "help"):
		return "Доступные команды:\n/help - справка\n/ping - проверка работы"
	default:
		return "Я получил ваше сообщение: " + message
	}
}