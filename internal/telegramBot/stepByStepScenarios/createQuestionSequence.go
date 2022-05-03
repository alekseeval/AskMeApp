package stepByStepScenarios

import (
	"AskMeApp/internal/interfaces"
	"AskMeApp/internal/model"
	"AskMeApp/internal/telegramBot/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	waitForCategory = iota
	waitForTitle
	waitForAnswer
	saveQuestion
	end
)

type CreateQuestionsSequence struct {
	step     int
	question *model.Question
}

func NewCreateQuestionsSequence(user *model.User) *CreateQuestionsSequence {
	return &CreateQuestionsSequence{
		step: waitForCategory,
		question: &model.Question{
			Author: user,
		},
	}
}

func (sequence *CreateQuestionsSequence) doStep(botClient *client.BotClient, update *tgbotapi.Update, questionsRepo interfaces.QuestionsRepositoryInterface) error {
	switch sequence.step {
	case waitForCategory:
		err := botClient.SendTextMessage("Select category:", sequence.question.Author.TgChatId)
		categories, err := questionsRepo.GetAllCategories()
		// TODO: бросить inline кнопки с категориями в чат
		log.Println(categories)
		if err != nil {
			return err
		}
	case waitForTitle:
		// TODO: правильно записать категорию через callback и обработать повторный запрос категории, если пользователь ввел текст
		sequence.question.Category = &model.Category{
			Id:    1,
			Title: "Все вопросы",
		}
		err := botClient.SendTextMessage("Enter title:", sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}
	case waitForAnswer:
		sequence.question.Title = update.Message.Text
		err := botClient.SendTextMessage("Enter answer:", sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}
	case saveQuestion:
		sequence.question.Answer = update.Message.Text
		_, err := questionsRepo.AddQuestion(sequence.question)
		if err != nil {
			return err
		}
		err = botClient.SendTextMessage("Question successfully saved", sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}
	default:
		log.Println("Smth went wrong with sequence step -", sequence.step)
	}
	sequence.step++
	return nil
}

func (sequence *CreateQuestionsSequence) isDone() bool {
	return sequence.step == end
}
