package telegramBot

import (
	"AskMeApp/internal"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (botClient *BotClient) SendCategoriesToChooseInline(explanatoryMessage string, categories []*internal.Category, chatId int64) error {
	inlineButtons := formatCategoriesToInlineButtons(categories)
	inlineMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	msg := tgbotapi.NewMessage(chatId, explanatoryMessage)
	msg.ReplyMarkup = inlineMarkup
	_, err := botClient.botApi.Send(msg)
	return err
}

func (botClient *BotClient) SendRandomQuestionToUser(user *internal.User) error {
	allQuestions, err := botClient.questionRepository.GetAllQuestions()
	if err != nil {
		return err
	}
	if len(allQuestions) == 0 {
		msg := tgbotapi.NewMessage(user.TgChatId, "На данный момент ваша База знаний пуста")
		_, err = botClient.botApi.Send(msg)
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
		_, err = botClient.botApi.Send(msg)
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
	_, err = botClient.botApi.Send(msg)
	return err
}

func (botClient *BotClient) setCustomKeyboardToChat(tgChatId int64) error {
	keyBoardFirstRow := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(randomQuestionCommandText))
	keyBoardSecondRow := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(changeCategoryCommandText))

	replyKeyboard := tgbotapi.NewReplyKeyboard(keyBoardFirstRow, keyBoardSecondRow)
	msg := tgbotapi.NewMessage(tgChatId, "Welcome!")
	msg.ReplyMarkup = replyKeyboard
	_, err := botClient.botApi.Send(msg)
	return err
}
