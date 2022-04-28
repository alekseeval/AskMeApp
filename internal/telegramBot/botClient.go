package telegramBot

import (
	"AskMeApp/internal/interfaces"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type BotClient struct {
	bot     *TgBotApi.BotAPI
	updates TgBotApi.UpdatesChannel
	wg      *sync.WaitGroup
}

func NewBotClient(botToken string) (*BotClient, error) {
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
	}, nil
}

func (bot *BotClient) SendTextMessage(msgText string, chatId int64) error {
	msg := TgBotApi.NewMessage(chatId, msgText)
	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (bot *BotClient) Run(userRepository interfaces.UserRepositoryInterface, wg *sync.WaitGroup) {
	bot.wg = wg
	bot.wg.Add(1)
	go HandleBotUpdates(bot, userRepository)
}

func (bot *BotClient) Stop() {
	bot.wg.Done()
}
