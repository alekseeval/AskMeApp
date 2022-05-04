package client

import (
	"AskMeApp/internal/model"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
)

const OptimalSymbolsNumberInRow int = 36

func (bot *BotClient) SendTextMessage(msgText string, chatId int64) error {
	msg := TgBotApi.NewMessage(chatId, msgText)
	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (bot *BotClient) SendInlineCategories(messageExplanation string, categories []*model.Category, chatId int64) error {
	inlineButtons := formatCategoriesToInline(categories)
	inlineMarkup := TgBotApi.NewInlineKeyboardMarkup(inlineButtons...)
	msg := TgBotApi.NewMessage(chatId, messageExplanation)
	msg.ReplyMarkup = inlineMarkup
	_, err := bot.bot.Send(msg)
	return err
}

func formatCategoriesToInline(categories []*model.Category) [][]TgBotApi.InlineKeyboardButton {
	categoriesCopy := append([]*model.Category(nil), categories...)
	sort.Slice(categoriesCopy, func(i, j int) bool {
		return len(categoriesCopy[i].Title) < len(categoriesCopy[j].Title)
	})

	weights := make([]int, 0)
	for _, c := range categoriesCopy {
		weight := OptimalSymbolsNumberInRow / len(c.Title)
		if weight > 4 {
			weight = 4
		}
		if weight == 0 {
			weight = 1
		}
		weights = append(weights, weight)
	}

	inlineButtons := make([][]TgBotApi.InlineKeyboardButton, 1)
	weightSum := 0
	for i := range categoriesCopy {
		weightSum += weights[i]
		if weightSum > 4 {
		}
	}

	inlineButtons := make([][]TgBotApi.InlineKeyboardButton, 1)
	curNumberInRow := 4
	i := 0
	j := curNumberInRow
	for _, category := range categoriesCopy {
		if len(category.Title) < OptimalSymbolsNumberInRow/curNumberInRow {

		}
		//inlineButtons[i/numberInRow] = append(inlineButtons[i/numberInRow], TgBotApi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
	}
	return
}
