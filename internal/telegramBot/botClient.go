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

	randomQuestionCommandText = "âAsk me"
	changeCategoryCommandText = "ð Select category"
	addQuestionCommandText    = "â Add new question"
	cancelAllStepsCommandText = "â Cancel"
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
		log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð²Ð¾ Ð²ÑÐµÐ¼Ñ Ð°Ð²ÑÐ¾ÑÐ¸Ð·Ð°ÑÐ¸Ð¸ Ð¿Ð¾Ð»ÑÐ·Ð¾Ð²Ð°ÑÐµÐ»Ñ: ", err)
	}
	botClient.statesMutex.Lock()
	userState, ok := botClient.usersStates[user.TgChatId]
	if ok {
		botClient.statesMutex.Unlock()
		if userState.SequenceStep != nilStep {
			userState.mutex.Lock()
			defer userState.mutex.Unlock()
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ¿Ð¾Ð»Ð½ÐµÐ½Ð¸Ð¸ ÑÐ°Ð³Ð° ÑÐµÐ¿Ð¾ÑÐºÐ¸ Ð´ÐµÐ¹ÑÑÐ²Ð¸Ð¹ Ð¿Ð¾Ð»ÑÐ·Ð¾Ð²Ð°ÑÐµÐ»Ñ: ", err)
			}
			return
		}
	} else {
		userState = NewUserState(baseCategory)
		botClient.usersStates[user.TgChatId] = userState
		botClient.statesMutex.Unlock()
	}
	userState.mutex.Lock()
	defer userState.mutex.Unlock()

	if update.Message != nil {

		switch update.Message.Command() {
		case startCommand:
			msg := tgbotapi.NewMessage(user.TgChatId, "Welcome to AskMeApp!")
			msg.ReplyMarkup = MainKeyboard
			_, err := botClient.botApi.Send(msg)
			if err != nil {
				log.Panic("ÐÐµ ÑÐ´Ð°Ð»Ð¾ÑÑ ÑÑÑÐ°Ð½Ð¾Ð²Ð¸ÑÑ ÐºÐ°ÑÑÐ¾Ð¼Ð½ÑÑ ÐºÐ»Ð°Ð²Ð¸Ð°ÑÑÑÑ Ð¿Ð¾ÑÐ»Ðµ ÐºÐ¾Ð¼Ð°Ð½Ð´Ñ /start: ", err)
			}
		case helpCommand:
			msg := tgbotapi.NewMessage(user.TgChatId, "ÐÑÐ¸Ð»Ð¾Ð¶ÐµÐ½Ð¸Ðµ Ð²ÑÐµ ÐµÑÐµ Ð½Ð°ÑÐ¾Ð´Ð¸ÑÑÑ Ð² ÑÐ°Ð·ÑÐ°Ð±Ð¾ÑÐºÐµ, Ð¿Ð¾ÑÑÐ¾Ð¼Ñ Ð¾Ð¿Ð¸ÑÐ°Ð½Ð¸Ðµ Ð½Ðµ Ð´Ð¾ÑÑÑÐ¿Ð½Ð¾. ÐÐ¶Ð¸Ð´Ð°Ð¹ÑÐµ ÑÐµÐ»Ð¸Ð·Ð° Ð² Ð±Ð»Ð¸Ð¶Ð°Ð¹ÑÐµÐµ Ð²ÑÐµÐ¼Ñ")
			_, err = botClient.botApi.Send(msg)
			if err != nil {
				log.Panic("ÐÐµ ÑÐ´Ð°Ð»Ð¾ÑÑ Ð¾ÑÐ¿ÑÐ°Ð²Ð¸ÑÑ ÑÐ¾Ð¾Ð±ÑÐµÐ½Ð¸Ðµ: ", err)
			}
		case randomQuestionCommand:
			err = botClient.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ´Ð°ÑÐµ Ð¿Ð¾Ð»ÑÐ·Ð¾Ð²Ð°ÑÐµÐ»Ñ ÑÐ»ÑÑÐ°Ð¹Ð½Ð¾Ð³Ð¾ Ð²Ð¾Ð¿ÑÐ¾ÑÐ°: ", err)
			}
		case changeCategoryCommand:
			userState.SequenceStep = ChangeCategoryInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ·Ð¾Ð²Ðµ ÐºÐ¾Ð¼Ð°Ð½Ð´Ñ ÑÐ¼ÐµÐ½Ñ ÐºÐ°ÑÐµÐ³Ð¾ÑÐ¸Ð¸ Ð¿Ð¾Ð»ÑÐ·Ð¾Ð²Ð°ÑÐµÐ»Ñ: ", err)
			}
		case addQuestionCommand:
			userState.SequenceStep = NewQuestionInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ·Ð¾Ð²Ðµ ÐºÐ¾Ð¼Ð°Ð½Ð´Ñ ÑÐ¾Ð·Ð´Ð°Ð½Ð¸Ñ Ð½Ð¾Ð²Ð¾Ð³Ð¾ Ð²Ð¾Ð¿ÑÐ¾ÑÐ°: ", err)
			}
		}

		switch update.Message.Text {
		case randomQuestionCommandText:
			err = botClient.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ´Ð°ÑÐµ Ð¿Ð¾Ð»ÑÐ·Ð¾Ð²Ð°ÑÐµÐ»Ñ ÑÐ»ÑÑÐ°Ð¹Ð½Ð¾Ð³Ð¾ Ð²Ð¾Ð¿ÑÐ¾ÑÐ°: ", err)
			}
		case changeCategoryCommandText:
			userState.SequenceStep = ChangeCategoryInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ·Ð¾Ð²Ðµ ÐºÐ¾Ð¼Ð°Ð½Ð´Ñ ÑÐ¼ÐµÐ½Ñ ÐºÐ°ÑÐµÐ³Ð¾ÑÐ¸Ð¸ Ð¿Ð¾Ð»ÑÐ·Ð¾Ð²Ð°ÑÐµÐ»Ñ: ", err)
			}
		case addQuestionCommandText:
			userState.SequenceStep = NewQuestionInitStep
			userState, err = botClient.ProcessUserStep(user, userState, update)
			if err != nil {
				log.Panic("Ð§ÑÐ¾-ÑÐ¾ Ð¿Ð¾ÑÐ»Ð¾ Ð½Ðµ ÑÐ°Ðº Ð¿ÑÐ¸ Ð²ÑÐ·Ð¾Ð²Ðµ ÐºÐ¾Ð¼Ð°Ð½Ð´Ñ ÑÐ¾Ð·Ð´Ð°Ð½Ð¸Ñ Ð½Ð¾Ð²Ð¾Ð³Ð¾ Ð²Ð¾Ð¿ÑÐ¾ÑÐ°: ", err)
			}
		}
	}
}
