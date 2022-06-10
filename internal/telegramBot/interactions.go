package telegramBot

import (
	"AskMeApp/internal"
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
	err = bot.SendStringMessageInChat("❔"+question.Title, user.TgChatId)
	return err
}
