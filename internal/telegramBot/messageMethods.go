package telegramBot

import (
	"AskMeApp/internal"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
)

const (
	OptimalSymbolsNumberInRow float32 = 36
	maxButtonsInLineNumber    float32 = 4
)

func (botClient *BotClient) SendStringMessageInChat(msgText string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, msgText)
	_, err := botClient.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (botClient *BotClient) SendCategoriesToChooseInline(messageExplanation string, categories []*internal.Category, chatId int64) error {
	inlineButtons := formatCategoriesToInline(categories)
	inlineMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	msg := tgbotapi.NewMessage(chatId, messageExplanation)
	msg.ReplyMarkup = inlineMarkup
	_, err := botClient.bot.Send(msg)
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

	inlineButtons := make([][]tgbotapi.InlineKeyboardButton, 1)
	var rowNumber int
	var weightSum float32
	for i, category := range categoriesCopy {
		weightSum += weights[i]
		if weightSum > maxButtonsInLineNumber {
			newRow := make([]tgbotapi.InlineKeyboardButton, 0)
			if len(inlineButtons[rowNumber]) == 0 {
				inlineButtons[rowNumber] = append(inlineButtons[rowNumber], tgbotapi.NewInlineKeyboardButtonData(category.Title, "c"+fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = 0
			} else {
				newRow = append(newRow, tgbotapi.NewInlineKeyboardButtonData(category.Title, "c"+fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = weights[i]
			}
			rowNumber += 1
		} else {
			inlineButtons[rowNumber] = append(inlineButtons[rowNumber], tgbotapi.NewInlineKeyboardButtonData(category.Title, "c"+fmt.Sprint(category.Id)))
		}
	}
	return inlineButtons
}

func (botClient *BotClient) SendRandomQuestionToUser(user *internal.User) error {
	allQuestions, err := botClient.questionRepository.GetAllQuestions()
	if err != nil {
		return err
	}
	if len(allQuestions) == 0 {
		err = botClient.SendStringMessageInChat("На данный момент ваша База знаний пуста", user.TgChatId)
		if err != nil {
			return err
		}
		return nil
	}

	currentUserState, ok := botClient.usersStates[user.TgChatId]
	if !ok {
		return errors.New("user have no current state")
	}
	requestedQuestions := make([]*internal.Question, 0)
	for _, question := range allQuestions {
		for _, category := range question.Categories {
			if category.Id == currentUserState.CurrentCategory.Id {
				requestedQuestions = append(requestedQuestions, question)
				continue
			}
		}
	}
	if len(requestedQuestions) == 0 {
		msg := tgbotapi.NewMessage(user.TgChatId, "На данный момент вопросы по категории __"+currentUserState.CurrentCategory.Title+"__ отсутствуют")
		msg.ParseMode = "MarkdownV2"
		_, err = botClient.bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	}

	question := GetRandomQuestion(requestedQuestions)
	themesText := ""
	for _, category := range question.Categories {
		themesText += "\t__" + category.Title + "__"
	}
	msg := tgbotapi.NewMessage(user.TgChatId, themesText+
		"\n\n*Question:\n*_"+question.Title+"_")
	msg.ParseMode = "MarkdownV2"
	_, err = botClient.bot.Send(msg)
	return err
}

func (botClient *BotClient) setCustomKeyboardToUser(user *internal.User) error {
	keyBoardFirstRow := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(randomQuestionCommandText))
	keyBoardSecondRow := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(changeCategoryCommandText))

	replyKeyboard := tgbotapi.NewReplyKeyboard(keyBoardFirstRow, keyBoardSecondRow)
	msg := tgbotapi.NewMessage(user.TgChatId, "Welcome!")
	msg.ReplyMarkup = replyKeyboard
	_, err := botClient.bot.Send(msg)
	return err
}
