package client

import TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (bot *BotClient) SendTextMessage(msgText string, chatId int64) error {
	msg := TgBotApi.NewMessage(chatId, msgText)
	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
