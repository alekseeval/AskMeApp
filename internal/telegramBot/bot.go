package telegramBot

import (
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotClient struct {
	Bot     *TgBotApi.BotAPI
	Updates TgBotApi.UpdatesChannel
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
		Bot:     bot,
		Updates: updates,
	}, nil
}

func (c *BotClient) SendMessage(msgText string, chatId int64) error {
	msg := TgBotApi.NewMessage(chatId, msgText)
	_, err := c.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
