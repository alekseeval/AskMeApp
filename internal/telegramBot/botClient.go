package telegramBot

import (
	"AskMeApp/internal"
	"context"
	"errors"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type BotClient struct {
	bot     *TgBotApi.BotAPI
	updates TgBotApi.UpdatesChannel

	cancelFunc context.CancelFunc

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

func (bot *BotClient) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	bot.cancelFunc = cancelFunc
	var wg sync.WaitGroup
	wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			break
		case update := <-bot.updates:
			wg.Add(1)
			bot.handleUpdate(&wg, &update)
			continue
		}
		break
	}
	log.Print("Waiting for all processes..")
	wg.Wait()
}

func (bot *BotClient) Shutdown() error {
	if bot.cancelFunc != nil {
		bot.cancelFunc()
		log.Print("Shutdown..")
	} else {
		return errors.New("bot isn't running yet")
	}
	return nil
}

func (bot *BotClient) handleUpdate(wg *sync.WaitGroup, update *TgBotApi.Update) {

	defer wg.Done()

	if update.Message != nil {

		user, err := VerifyOrRegisterUser(update.Message.Chat.ID, update.Message.From, bot.userRepository)
		if err != nil {
			err = bot.SendTextMessage("Что-то пошло не так во время авторизации: \n"+err.Error(), update.Message.Chat.ID)
			if err != nil {
				log.Panic("Жопа наступила, не удалось получить или создать юзера,"+
					" а потом еще и сообщение не отправилось", err)
			}
		}

		switch update.Message.Command() {
		case "start":
			err = bot.SendTextMessage("Это была команда /start", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "help":
			err = bot.SendTextMessage("Это была команда /help", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "question":
			err = bot.SendTextMessage("Это была команда /question", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "changecategory":
			err = bot.SendTextMessage("Это была команда /changecategory", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "shutdown":
			if user.TgUserName != "al_andrew" {
				log.Print("Нарушитель пытался завершить работу приложения", user.TgUserName)
				return
			}
			err = bot.SendTextMessage("Приложение завершило свою работу", user.TgChatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
			err = bot.Shutdown()
			if err != nil {
				log.Panic("Запущенный бот не запущен", err)
			}
		}
	}
}
