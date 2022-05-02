package stepByStepScenarios

import (
	"AskMeApp/internal/model"
	"AskMeApp/internal/telegramBot/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	start = iota
	fillCategory
	fillTitle
	fillAnswer
	end
)

type CreateQuestionsSequence struct {
	step     int
	Question *model.Question
}

func NewCreateQuestionsSequence(user *model.User) *CreateQuestionsSequence {
	return &CreateQuestionsSequence{
		step: start,
		Question: &model.Question{
			Author: user,
		},
	}
}

func (sequence *CreateQuestionsSequence) doStep(botClient *client.BotClient, update *tgbotapi.Update) {
	switch sequence.step {
	case start:
		// TODO: отправить сообщение с запросом категории вопроса
	case fillCategory:
		// TODO: записать категорию, спросить Title
	case fillTitle:
		// TODO: записать Title спросить Answer
	case fillAnswer:
		// TODO: записать Answer
	default:
		log.Println("Smth went wrong with sequence step -", sequence.step)
	}
	sequence.step++
}

func (sequence *CreateQuestionsSequence) isDone() bool {
	return sequence.step == end
}

func (sequence *CreateQuestionsSequence) getEntity() *model.Question {
	return sequence.Question
}
