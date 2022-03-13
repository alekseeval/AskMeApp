package telegram_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramBotClient struct {
	Bot         *tgbotapi.BotAPI
	UpdatesChan tgbotapi.UpdatesChannel
}

func NewTelegramBotClient(ApiToken string) *TelegramBotClient {
	bot, err := tgbotapi.NewBotAPI(ApiToken)
	if err != nil {
		log.Panic(err)
	}

	// TODO: Разобраться с конфигом в API
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updatesChan := bot.GetUpdatesChan(updateConfig)

	return &TelegramBotClient{
		Bot:         bot,
		UpdatesChan: updatesChan,
	}
}

func (c *TelegramBotClient) SendTextMessage(message fmt.Stringer, chatId int64) error {
	msgText := message.String()
	msg := tgbotapi.NewMessage(chatId, msgText)
	// TODO: Разобраться с отправкой сообщений и возможными параметрами
	msg.ParseMode = "MarkDown"
	_, err := c.Bot.Send(msg)
	return err
}
