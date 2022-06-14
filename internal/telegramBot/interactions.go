package telegramBot

import (
	"AskMeApp/internal"
	"sync"
)

const NilStep int = 0

const (
	QuestionEndStep int = iota*3 + 1
	Question2Step
	Question3Step
	Question4Step
)

const (
	CategoryEndStep int = iota*3 + 2
	Category2Step
)

type userState struct {
	CurrentCategory internal.Category
	SequenceStep    int

	unfilledQuestion *internal.Question
	unfilledCategory *internal.Category

	mutex *sync.Mutex
}

func (bot BotClient) ProcessUserStep(user *internal.User, stepState *userState) error {
	// TODO: Написать switch на каждый Step из констант с обработкой сценария, изменить переданный state
	return nil
}
