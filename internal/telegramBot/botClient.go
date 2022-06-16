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
	randomQuestionCommandText = "Gimme question"
	randomQuestionCommand     = "question"
	helpCommand               = "help"
	startCommand              = "start"
)

type BotClient struct {
	bot     *TgBotApi.BotAPI
	updates TgBotApi.UpdatesChannel

	cancelFunc context.CancelFunc
	wg         sync.WaitGroup

	usersStates map[int64]userState
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
	}, nil
}

func (bot *BotClient) Run() error {
	if bot.cancelFunc != nil {
		return errors.New("bot already running")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	bot.cancelFunc = cancelFunc
	bot.wg = sync.WaitGroup{}
	bot.wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			break
		case update := <-bot.updates:
			bot.wg.Add(1)
			go bot.handleUpdate(&update)
			continue
		}
		break
	}
	bot.wg.Done()
	return nil
}

func (bot *BotClient) Shutdown() error {
	if bot.cancelFunc == nil {
		return errors.New("bot isn't running yet")
	}
	bot.cancelFunc()
	bot.cancelFunc = nil

	log.Print("Waiting for all processes..")
	bot.wg.Wait()
	return nil
}

func (bot *BotClient) handleUpdate(update *TgBotApi.Update) {

	defer bot.wg.Done()

	user, err := VerifyOrRegisterUser(update.SentFrom(), bot.userRepository)
	state, ok := bot.usersStates[user.TgChatId]
	if ok && state.SequenceStep != NilStep {
		state.mutex.Lock()
		defer state.mutex.Unlock()
		err = bot.ProcessUserStep(user, &state)
		if err != nil {
			log.Panic(err)
		}
		return
	} else {
		// TODO:
		//state = NewUserState()
	}
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if update.Message != nil {
		if err != nil {
			err = bot.SendStringMessageInChat("Что-то пошло не так во время авторизации: \n"+err.Error(), update.Message.Chat.ID)
			if err != nil {
				log.Panic("Жопа наступила, не удалось получить или создать юзера,"+
					" а потом еще и сообщение не отправилось", err)
			}
		}

		switch update.Message.Command() {
		case startCommand:
			err = bot.setCustomKeyboardToUser(user)
			if err != nil {
				log.Panic("Не удалось установить клавиатуру", err)
			}
		case helpCommand:
			err = bot.SendStringMessageInChat("Это была команда /help", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case randomQuestionCommand:
			err = bot.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic(err)
			}
		}

		switch update.Message.Text {
		case randomQuestionCommandText:
			err = bot.SendRandomQuestionToUser(user)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
