package telegramBot

import (
	"AskMeApp/internal"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
