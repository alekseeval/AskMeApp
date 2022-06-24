package telegramBot

import (
	"AskMeApp/internal"
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
	"time"
)

const (
	randomQuestionCommand = "question"
	helpCommand           = "help"
	startCommand          = "start"
	changeCategoryCommand = "changecategory"
	addQuestionCommand    = "newquestion"

	randomQuestionCommandText = "❔Ask me"
	changeCategoryCommandText = "🔄 Select category"
	addQuestionCommandText    = "➕ Add new question"
	cancelAllStepsCommandText = "❌ Cancel"
)

var baseCategory = internal.Category{
	Id:    1,
	Title: "All",
}

type BotClient struct {
	botApi  *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel

	cancelFunc context.CancelFunc
	wg         sync.WaitGroup

	usersStates map[int64]*userState
	statesMutex sync.Mutex

	userRepository     internal.UserRepositoryInterface
	questionRepository internal.QuestionsRepositoryInterface
}

func NewBotClient(userRepository internal.UserRepositoryInterface, questionRepository internal.QuestionsRepositoryInterface, botToken string) (*BotClient, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updatesConfig)
	select {
	case <-updates:
		break
	case <-time.After(time.Second * 2):
		break
	}
	log.Print("Connected to updates..")
	updates.Clear()
	return &BotClient{
		botApi:  bot,
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

func (botClient *BotClient) Shutdown(timeout time.Duration) error {
	if botClient.cancelFunc == nil {
		return errors.New("botClient isn't running yet")
	}
	botClient.cancelFunc()
	botClient.cancelFunc = nil
	log.Print("Waiting for all processes..")
	c := make(chan struct{})
	go func() {
		defer close(c)
		botClient.wg.Wait()
	}()
	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return errors.New("some of bot workers doesn't stopped")
	}
}

func (botClient *BotClient) handleUpdate(update *tgbotapi.Update) {

	defer botClient.wg.Done()

	user, err := IdentifyOrRegisterUser(update.SentFrom(), botClient.userRepository)
	if err != nil {
		log.Panic("Что-то пошло не так во время авторизации пользователя: ", err)
	}
	botClient.statesMutex.Lock()
	userState, ok := botClient.usersStates[user.TgChatId]
	if ok {
		if userState.SequenceStep != nilStep {
			userState.mutex.Lock()
			defer userState.mutex.Unlock()
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Что-то пошло не так при выполнении шага цепочки действий пользователя: ", err)
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

		switch update.Message.Command() {
		case startCommand:
			msg := tgbotapi.NewMessage(user.TgChatId, "Welcome to AskMeApp!")
			msg.ReplyMarkup = MainKeyboard
			_, err := botClient.botApi.Send(msg)
			if err != nil {
				log.Panic("Не удалось установить кастомную клавиатуру после команды /start: ", err)
			}
		case helpCommand:
			msg := tgbotapi.NewMessage(user.TgChatId, "Приложение все еще находится в разработке, поэтому описание не доступно. Ожидайте релиза в ближайшее время")
			_, err = botClient.botApi.Send(msg)
			if err != nil {
				log.Panic("Не удалось отправить сообщение: ", err)
			}
		case randomQuestionCommand:
			err = botClient.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic("Что-то пошло не так при выдаче пользователю случайного вопроса: ", err)
			}
		case changeCategoryCommand:
			userState.SequenceStep = ChangeCategoryInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Что-то пошло не так при вызове команды смены категории пользователя: ", err)
			}
		case addQuestionCommand:
			userState.SequenceStep = NewQuestionInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Что-то пошло не так при вызове команды создания нового вопроса: ", err)
			}
		}

		switch update.Message.Text {
		case randomQuestionCommandText:
			err = botClient.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic("Что-то пошло не так при выдаче пользователю случайного вопроса: ", err)
			}
		case changeCategoryCommandText:
			userState.SequenceStep = ChangeCategoryInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Что-то пошло не так при вызове команды смены категории пользователя: ", err)
			}
		case addQuestionCommandText:
			userState.SequenceStep = NewQuestionInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Что-то пошло не так при вызове команды создания нового вопроса: ", err)
			}
		}
	}
}
