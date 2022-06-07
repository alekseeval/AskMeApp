package telegramBot

import (
	"AskMeApp/internal"
	"fmt"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sort"
)

const OptimalSymbolsNumberInRow float32 = 36
const maxButtonsInLineNumber float32 = 4

func (bot BotClient) SendTextMessage(msgText string, chatId int64) error {
	msg := TgBotApi.NewMessage(chatId, msgText)
	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (bot BotClient) SendInlineCategories(messageExplanation string, categories []*internal.Category, chatId int64) error {
	inlineButtons := formatCategoriesToInline(categories)
	inlineMarkup := TgBotApi.NewInlineKeyboardMarkup(inlineButtons...)
	msg := TgBotApi.NewMessage(chatId, messageExplanation)
	msg.ReplyMarkup = inlineMarkup
	_, err := bot.bot.Send(msg)
	return err
}

func formatCategoriesToInline(categories []*internal.Category) [][]TgBotApi.InlineKeyboardButton {
	categoriesCopy := append([]*internal.Category(nil), categories...)
	sort.Slice(categoriesCopy, func(i, j int) bool {
		return len(categoriesCopy[i].Title) < len(categoriesCopy[j].Title)
	})

	weights := make([]float32, 0)
	for _, c := range categoriesCopy {
		weight := OptimalSymbolsNumberInRow / float32(len(c.Title))
		if weight > maxButtonsInLineNumber {
			weight = maxButtonsInLineNumber
		}
		if weight == maxButtonsInLineNumber {
			weight = maxButtonsInLineNumber
		}
		weights = append(weights, maxButtonsInLineNumber/weight)
	}
	log.Println(weights)

	inlineButtons := make([][]TgBotApi.InlineKeyboardButton, 1)
	var rowNumber int
	var weightSum float32
	for i, category := range categoriesCopy {
		weightSum += weights[i]
		if weightSum > maxButtonsInLineNumber {
			newRow := make([]TgBotApi.InlineKeyboardButton, 0)
			if len(inlineButtons[rowNumber]) == 0 {
				inlineButtons[rowNumber] = append(inlineButtons[rowNumber], TgBotApi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = 0
			} else {
				newRow = append(newRow, TgBotApi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = weights[i]
			}
			rowNumber += 1
		} else {
			inlineButtons[rowNumber] = append(inlineButtons[rowNumber], TgBotApi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
		}
	}
	return inlineButtons
}
