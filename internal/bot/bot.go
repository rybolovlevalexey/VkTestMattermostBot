package bot

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"

	"VkTestMattermostBot/internal/config"
	"VkTestMattermostBot/internal/core"
	"VkTestMattermostBot/internal/database"
	"VkTestMattermostBot/internal/usecases"
)

var BotCommands = []string{"create", "votename", "votedesc", "votevariants", "voteoneanswer", 
						   "votestart", "cast", "check", "", "", "", "", "", }
var BotCommandsWithId = []string{"votename", "votedesc", "votevariants", "voteoneanswer", 
							     "votestart", "cast", "check"}


type MattermostBot struct{
	Client	*model.Client4;
	WSclient 	*model.WebSocketClient;
	BotConfig	config.BotConfig;
}

type InfoToGenerateResponse struct{
	voteId int;
	chanelId string;
	creatorId string;
	updatingNameDone bool;
	updatingDescDone bool;
	updatingVarinatsDone bool;
	updatingIsOneAnswerDone bool;
	updatingVoteStart bool;
	UserCastVoteByVoteIdDone bool;
	ViewVoteInfoDone map[string]bool;  // done - выполнено ли успешно, not_exist - не существует, from_another_chanel - из другого канала
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
		// log.Println("Получено не новое сообщение")
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
		// log.Println("Получено сообщение от самого бота")
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
	resultMainLogic, infoGenerateResp := mainLogic(post.Message, botConfig, post.UserId, post.ChannelId)
	log.Println("MainLogic ", resultMainLogic)

	// создание сообщения, отвечающего пользователю на его запрос
	reply := &model.Post{
		ChannelId: post.ChannelId,
		Message:   generateResponse(post.Message, botConfig, infoGenerateResp, resultMainLogic),
	}

	if _, _, err := client.CreatePost(reply); err != nil {
		log.Printf("Ошибка отправки ответа: %v", err)
	}
}

// генерация ответов в зависимости от сообщения пользователя и результатов выполнения логики приложения
func generateResponse(message string, 
					  botConfig config.BotConfig, 
					  infoGenerateResp InfoToGenerateResponse, 
					  resultMainLogic database.VoteModel) string {
	message = strings.TrimSpace(message)
	flagCommandInMessage := false

	for _, command := range BotCommands{
		if strings.Contains(message, command){
			flagCommandInMessage = true
			break
		}
	}

	if !flagCommandInMessage{  // сообщении нет команды
		return "В отправленном сообщении нет команды или название команды напечатно некорректно"
	}

	if strings.Contains(message, "votestart"){
		curVoteId, _ := strconv.Atoi(strings.TrimSpace(strings.Split(message, " ")[2]))
		curVote := database.GetVoteInfoById(curVoteId)

		core.AppLogger.Println("curVote.Name curVote.Variants ", curVote.Name, curVote.Variants)
		if curVote.Name == "" || len(curVote.Variants) < 2{
			return "Голосование с id - " + strconv.Itoa(curVoteId) + " не готово к старту (не заполнена обязательная информация)"
		}
	}

	switch {
		case strings.Contains(message, "help"):  
			// получено сообщение с help
			return BotAnswers["help"]
		case strings.TrimSpace(strings.Replace(strings.Replace(message, botConfig.BotUserName, "", 1), "@", "", 1)) == "":  
			// получено пустое сообщение
			return BotAnswers["help"]
		case strings.Contains(message, "create"):  // получена команда на создание нового голосования
			return "Создано голосование с id - " + strconv.Itoa(infoGenerateResp.voteId) + "\n\n" + BotAnswers["create"]
		case strings.Contains(message, "votename"):  // получена команда на установку названия голосования
			if infoGenerateResp.updatingNameDone{
				return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " установлено название"
			}
			return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " не было установлено название (ошибка прав доступа)"
		case strings.Contains(message, "votedesc"):
			if infoGenerateResp.updatingDescDone{
				return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " установлено описание"
			}
			return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " не установлено описание (ошибка прав доступа)"
		case strings.Contains(message, "votevariants"):
			if infoGenerateResp.updatingVarinatsDone{
				return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " установлены варианты ответа"
			}
			return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " не установлены варианты ответа (ошибка прав доступа)"
		case strings.Contains(message, "voteoneanswer"):
			if infoGenerateResp.updatingIsOneAnswerDone{
				return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " установлено является ли голосование с одним вариантом ответа или с несколькими"
			}
			return "У голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " не установлено является ли голосование с одним вариантом ответа или с несколькими (ошибка прав доступа)"
		case strings.Contains(message, "votestart"):
			if infoGenerateResp.updatingVoteStart{
				return "Голосование с id - " + strconv.Itoa(infoGenerateResp.voteId) + " начато"
			}
			return "Голосование с id - " + strconv.Itoa(infoGenerateResp.voteId) + " не начато (ошибка прав доступа)"
		case strings.Contains(message, "cast"):
			if infoGenerateResp.UserCastVoteByVoteIdDone{
				return "Вы успешно проголосовали в голосовании с id - " + strconv.Itoa(infoGenerateResp.voteId)
			}
			return "Вы не смогли проголосовать в голосовании с id - " + strconv.Itoa(infoGenerateResp.voteId) + ". Скорее всего такого голосования в вашем канале не существует или вы уже проголосовали в указанном голосовании. "
		case strings.Contains(message, "check"):
			if infoGenerateResp.ViewVoteInfoDone["done"]{
				st := "Информация о голосовании с id - " + strconv.Itoa(infoGenerateResp.voteId) + ":\n" + 
					  "Название - " + resultMainLogic.Name + "\n"
				if len(resultMainLogic.Description) > 0{
					st += "Описание - " + resultMainLogic.Description + "\n"
				}
				st += "Варианты ответа и количество пользователей за них проголосовавших:\n"
				for key, value := range resultMainLogic.Variants{
					st += key + " - " + strconv.Itoa(len(value)) + "\n"
				}
				if resultMainLogic.IsActive{
					st += "Голосование ещё не окончено\n"
				} else {
					st += "Голосование уже завершено\n"
				}
				return st
			}
			if infoGenerateResp.ViewVoteInfoDone["not_exist"]{
				return "Голосования с id - " + strconv.Itoa(infoGenerateResp.voteId) + " не существует"
			}
			return "В вашем канале нет голосования с id - " + strconv.Itoa(infoGenerateResp.voteId)
		default:
			return "" + message
	}
}

// запуск необходимых функций в соответствии с полученным сообщением от пользователя
func mainLogic(message string, 
			   botConfig config.BotConfig, 
			   userMatterMostId string, 
			   chanelId string) (database.VoteModel, InfoToGenerateResponse){
	var result database.VoteModel
	log.Println(message, botConfig.BotUserName, userMatterMostId)

	flagCommandInMessage := false

	for _, command := range BotCommands{
		if strings.Contains(message, command){
			flagCommandInMessage = true
			break
		}
	}

	if !flagCommandInMessage{
		return result, InfoToGenerateResponse{}
	}

	flagCommandWithIdInMessage := false
	for _, command := range BotCommandsWithId{
		if strings.Contains(message, command){
			flagCommandWithIdInMessage = true
			break
		}
	}

	var messageSplited []string
	var voteIdStr string
	var voteId int

	if flagCommandWithIdInMessage{
		// @название_бота команда id_голосования <доп информация>
		messageSplited = strings.Split(message, " ")
		voteIdStr = messageSplited[2]
		voteId, _ = strconv.Atoi(voteIdStr)
	}
	

	switch {
	case strings.Contains(message, "create"):
		newVoteId := usecases.CreateVote(userMatterMostId, chanelId)
		core.AppLogger.Println("create", newVoteId)
		return result, InfoToGenerateResponse{voteId: newVoteId}
	
	case strings.Contains(message, "votename"):
		voteName := strings.Join(messageSplited[3:], " ")
		core.AppLogger.Println("votename", voteName, voteId)

		resSetVoteName := usecases.SetVoteName(userMatterMostId, voteId, voteName)

		return result, InfoToGenerateResponse{updatingNameDone: resSetVoteName, voteId: voteId}
	
	case strings.Contains(message, "votedesc"):
		voteDesc := strings.Join(messageSplited[3:], " ")
		core.AppLogger.Println("votedesc", voteDesc, voteId)

		resSetVoteDesc := usecases.SetVoteDesc(userMatterMostId, voteId, voteDesc)
		return result, InfoToGenerateResponse{updatingDescDone: resSetVoteDesc, voteId: voteId}
	
	case strings.Contains(message, "votevariants"):
		voteVariants := strings.Join(messageSplited[3:], " ")
		core.AppLogger.Println("votevariants", voteVariants, voteId)

		voteVariantsList := make([]string, 0)
		for _, elem := range strings.Split(voteVariants, ";"){
			voteVariantsList = append(voteVariantsList, strings.TrimSpace(elem))
		}

		resSetVoteVariants := usecases.SetVoteVariants(userMatterMostId, voteId, voteVariantsList)

		return result, InfoToGenerateResponse{updatingVarinatsDone: resSetVoteVariants, voteId: voteId}
	
	case strings.Contains(message, "voteoneanswer"):
		voteOneAnswer := strings.Join(messageSplited[3:], " ")
		core.AppLogger.Println("voteoneanswer", voteOneAnswer, voteId)
		
		var voteOneAnswerBool bool
		if strings.TrimSpace(voteOneAnswer) == "Y"{
			voteOneAnswerBool = true
		} else {
			voteOneAnswerBool = false
		}
			
		resSetVoteIsOneAnswer := usecases.SetVoteIsOneAnswer(userMatterMostId, voteId, voteOneAnswerBool)
		return result, InfoToGenerateResponse{updatingIsOneAnswerDone: resSetVoteIsOneAnswer, voteId: voteId}
	
	case strings.Contains(message, "votestart"):
		core.AppLogger.Println("votestart", voteId)
		resVoteStart := usecases.StartVote(userMatterMostId, voteId)
		return result, InfoToGenerateResponse{updatingVoteStart: resVoteStart, voteId: voteId}

	case strings.Contains(message, "cast"):
		variants := strings.Split(strings.Join(messageSplited[3:], " "), ";")
		resUserCastVoteByVoteId := usecases.UserCastVoteByVoteId(userMatterMostId, voteId, chanelId, variants)
		return result, InfoToGenerateResponse{UserCastVoteByVoteIdDone: resUserCastVoteByVoteId, voteId:  voteId}
	
	case strings.Contains(message, "check"):
		resViewCurrentVoteResult := usecases.ViewCurrentVoteResult(voteId, chanelId)
		
		if resViewCurrentVoteResult.Id >= 0{
			return resViewCurrentVoteResult, InfoToGenerateResponse{
				ViewVoteInfoDone: map[string]bool{"done": true, "not_exist": false, "from_another_chanel": false}, 
				voteId:  voteId,
			}
		}
		if resViewCurrentVoteResult.Id == -1{
			return result, InfoToGenerateResponse{
				ViewVoteInfoDone: map[string]bool{"done": false, "not_exist": true, "from_another_chanel": false}, 
				voteId:  voteId,
			}
		}
		if resViewCurrentVoteResult.Id == -2{
			return result, InfoToGenerateResponse{
				ViewVoteInfoDone: map[string]bool{"done": false, "not_exist": false, "from_another_chanel": true}, 
				voteId:  voteId,
			}
		}
	}


	return result, InfoToGenerateResponse{}
}
