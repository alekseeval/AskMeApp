package stepByStepScenarios

import (
	"AskMeApp/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CreateQuestionsSequence struct {
	step     int
	Question *model.Question
	done     bool
}

func NewCreateQuestionsSequence(user *model.User) *CreateQuestionsSequence {
	return &CreateQuestionsSequence{
		step: 0,
		Question: &model.Question{
			Author: user,
		},
	}
}

func (sequence *CreateQuestionsSequence) doStep(update *tgbotapi.Update) {

}

func (sequence *CreateQuestionsSequence) isDone() bool {
	return sequence.done
}
