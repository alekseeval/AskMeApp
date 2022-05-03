package client

import (
	"AskMeApp/internal/model"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *BotClient) SendTextMessage(msgText string, chatId int64) error {
	msg := TgBotApi.NewMessage(chatId, msgText)
	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (bot *BotClient) SendInlineCategories(messageExplanation string, categories []*model.Category, numberInRow int, chatId int64) error {
	inlineButtons := make([][]TgBotApi.InlineKeyboardButton, (len(categories)-1)/numberInRow+1)
	for i, category := range categories {
		inlineButtons[i/numberInRow] = append(inlineButtons[i/numberInRow], TgBotApi.NewInlineKeyboardButtonData(category.Title, string(category.Id)))
	}
	inlineMarkup := TgBotApi.NewInlineKeyboardMarkup(inlineButtons...)
	msg := TgBotApi.NewMessage(chatId, messageExplanation)
	msg.ReplyMarkup = inlineMarkup
	_, err := bot.bot.Send(msg)
	return err
}
