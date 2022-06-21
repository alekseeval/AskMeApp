package telegramBot

import (
	"AskMeApp/internal"
	"context"
	"errors"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

const (
	randomQuestionCommand = "question"
	helpCommand           = "help"
	startCommand          = "start"
	changeCategoryCommand = "changecategory"

	randomQuestionCommandText = "❓Ask me"
	changeCategoryCommandText = "🔄 Change questions category"
)

var baseCategory = internal.Category{
	Id:    1,
	Title: "All",
}

type BotClient struct {
	bot     *TgBotApi.BotAPI
	updates TgBotApi.UpdatesChannel

	cancelFunc context.CancelFunc
	wg         sync.WaitGroup

	usersStates map[int64]*userState
	statesMutex sync.Mutex
	// TODO: встроить map[internal.User.Id]->*userState
	// 	 Хватать Mutex в userState и отпускать через defer в начале обработки каждого Update
	//	 Инициализировать запуск сценария с нужного шага при необходимости (скорее всего команда /newQuestion)

	userRepository     internal.UserRepositoryInterface
	questionRepository internal.QuestionsRepositoryInterface
}

func NewBotClient(userRepository internal.UserRepositoryInterface, questionRepository internal.QuestionsRepositoryInterface, botToken string) (*BotClient, error) {
	bot, err := TgBotApi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	updatesConfig := TgBotApi.NewUpdate(0)
	updatesConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updatesConfig)
	return &BotClient{
		bot:     bot,
		updates: updates,

		userRepository:     userRepository,
		questionRepository: questionRepository,

		usersStates: map[int64]*userState{},
	}, nil
}

func (botClient *BotClient) Run() error {
	if botClient.cancelFunc != nil {
		return errors.New("botClient already running")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	botClient.cancelFunc = cancelFunc
	botClient.wg = sync.WaitGroup{}
	botClient.wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			break
		case update := <-botClient.updates:
			botClient.wg.Add(1)
			go botClient.handleUpdate(&update)
			continue
		}
		break
	}
	botClient.wg.Done()
	return nil
}

func (botClient *BotClient) Shutdown() error {
	if botClient.cancelFunc == nil {
		return errors.New("botClient isn't running yet")
	}
	botClient.cancelFunc()
	botClient.cancelFunc = nil

	log.Print("Waiting for all processes..")
	botClient.wg.Wait()
	return nil
}

func (botClient *BotClient) handleUpdate(update *TgBotApi.Update) {

	defer botClient.wg.Done()

	user, err := VerifyOrRegisterUser(update.SentFrom(), botClient.userRepository)
	botClient.statesMutex.Lock()
	userState, ok := botClient.usersStates[user.TgChatId]
	if ok {
		if userState.SequenceStep != NilStep {
			userState.mutex.Lock()
			defer userState.mutex.Unlock()
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic(err)
			}
			botClient.statesMutex.Unlock()
			return
		}
	} else {
		userState = NewUserState(baseCategory)
		botClient.usersStates[user.TgChatId] = userState
	}
	botClient.statesMutex.Unlock()
	userState.mutex.Lock()
	defer userState.mutex.Unlock()

	if update.Message != nil {
		if err != nil {
			err = botClient.SendStringMessageInChat("Что-то пошло не так во время авторизации: \n"+err.Error(), update.Message.Chat.ID)
			if err != nil {
				log.Panic("Жопа наступила, не удалось получить или создать юзера,"+
					" а потом еще и сообщение не отправилось", err)
			}
		}

		switch update.Message.Command() {
		case startCommand:
			err = botClient.setCustomKeyboardToUser(user)
			if err != nil {
				log.Panic("Не удалось установить клавиатуру", err)
			}
		case helpCommand:
			err = botClient.SendStringMessageInChat("Это была команда /help", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case randomQuestionCommand:
			err = botClient.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic(err)
			}
		case changeCategoryCommand:
			userState.SequenceStep = ChangeCategoryInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic(err)
			}
		}

		switch update.Message.Text {
		case randomQuestionCommandText:
			err = botClient.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic(err)
			}
		case changeCategoryCommandText:
			userState.SequenceStep = ChangeCategoryInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
