package telegramBot

import (
	"AskMeApp/internal"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sort"
)

const (
	OptimalSymbolsNumberInRow float32 = 36
	maxButtonsInLineNumber    float32 = 4
)

func (bot *BotClient) SendStringMessageInChat(msgText string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, msgText)
	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (bot *BotClient) SendAllCategoriesInline(messageExplanation string, categories []*internal.Category, chatId int64) error {
	inlineButtons := formatCategoriesToInline(categories)
	inlineMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	msg := tgbotapi.NewMessage(chatId, messageExplanation)
	msg.ReplyMarkup = inlineMarkup
	_, err := bot.bot.Send(msg)
	return err
}

func formatCategoriesToInline(categories []*internal.Category) [][]tgbotapi.InlineKeyboardButton {
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

	inlineButtons := make([][]tgbotapi.InlineKeyboardButton, 1)
	var rowNumber int
	var weightSum float32
	for i, category := range categoriesCopy {
		weightSum += weights[i]
		if weightSum > maxButtonsInLineNumber {
			newRow := make([]tgbotapi.InlineKeyboardButton, 0)
			if len(inlineButtons[rowNumber]) == 0 {
				inlineButtons[rowNumber] = append(inlineButtons[rowNumber], tgbotapi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = 0
			} else {
				newRow = append(newRow, tgbotapi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = weights[i]
			}
			rowNumber += 1
		} else {
			inlineButtons[rowNumber] = append(inlineButtons[rowNumber], tgbotapi.NewInlineKeyboardButtonData(category.Title, fmt.Sprint(category.Id)))
		}
	}
	return inlineButtons
}

func (bot *BotClient) SendRandomQuestionToUser(user *internal.User) error {
	allQuestions, err := bot.questionRepository.GetAllQuestions()
	if err != nil {
		return err
	}
	if len(allQuestions) == 0 {
		err = bot.SendStringMessageInChat("На данный момент ваша База знаний пуста", user.TgChatId)
		if err != nil {
			return err
		}
		return nil
	}
	question := GetRandomQuestion(allQuestions)

	msg := tgbotapi.NewMessage(user.TgChatId, "*Theme:* __"+question.Category.Title+
		"__\n\n*Question:\n*_"+question.Title+"_")
	msg.ParseMode = "MarkdownV2"
	_, err = bot.bot.Send(msg)
	return err
}

func (bot *BotClient) setCustomKeyboardToUser(user *internal.User) error {
	keyBoardMarkup := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Gimme question"))
	replyKeyboard := tgbotapi.NewReplyKeyboard(keyBoardMarkup)
	msg := tgbotapi.NewMessage(user.TgChatId, "Welcome!")
	msg.ReplyMarkup = replyKeyboard
	_, err := bot.bot.Send(msg)
	return err
}
