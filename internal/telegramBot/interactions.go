package telegramBot

import (
	"AskMeApp/internal"
	"sync"
)

const NilStep int = 0

const (
	// TODO: реализовать добавление задачи по шагам
	NewQuestionEndStep int = iota*3 + 1
	NewQuestion2Step
	NewQuestion3Step
	NewQuestion4Step
)

const (
	ChangeCategoryEndStep int = iota*3 + 2
	ChangeCategoryChooseStep
)

type userState struct {
	CurrentCategory internal.Category
	SequenceStep    int

	unfilledQuestion *internal.Question
	unfilledCategory *internal.Category

	mutex sync.Mutex
}

func NewUserState(currentCategory internal.Category) *userState {
	return userState{
		CurrentCategory: currentCategory,
		SequenceStep:    NilStep,
		mutex:           sync.Mutex{},
	}
}

func (bot *BotClient) ProcessUserStep(user *internal.User, stepState *userState) error {
	// TODO: Написать switch на каждый Step из констант с обработкой сценария, изменить переданный state
	return nil
}
