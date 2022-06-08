package telegramBot

import (
	"AskMeApp/internal"
	"context"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BotClient struct {
	bot     *TgBotApi.BotAPI
	updates TgBotApi.UpdatesChannel

	ctx context.Context

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
	ctx := context.Background()
	return &BotClient{
		bot:     bot,
		updates: updates,
		ctx:     ctx,

		userRepository:     userRepository,
		questionRepository: questionRepository,
	}, nil
}

func (bot *BotClient) Run() {

	for update := range bot.updates {

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
					continue
				}
				err = bot.SendTextMessage("Приложение завершило свою работу", user.TgChatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
				bot.Shutdown()
				return
			}
		}
	}
}

func (bot *BotClient) Shutdown() {
	bot.ctx
}
