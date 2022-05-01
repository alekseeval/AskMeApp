package stepByStepScenarios

import (
	"AskMeApp/internal/interfaces"
	"AskMeApp/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StepByStepManager struct {
	questionRepo interfaces.QuestionsRepositoryInterface
	// TODO: предусмотреть mutex для мапы для работы в конкурентном режиме
	createQuestionSteps map[int32]*CreateQuestionsSequence
}

func newStepByStepManager(questionRepo interfaces.QuestionsRepositoryInterface) *StepByStepManager {
	return &StepByStepManager{
		questionRepo:        questionRepo,
		createQuestionSteps: make(map[int32]*CreateQuestionsSequence),
	}
}

func (manager *StepByStepManager) doUserHaveSequences(user *model.User) bool {
	if manager.createQuestionSteps[user.Id] != nil {
		return true
	}
	return false
}

func (manager *StepByStepManager) DoSequenceStep(user *model.User, update *tgbotapi.Update) error {
	// 1) Выполнение степа (сначала надо найти его в мапе)
	// 2) Проверка завершенности степа
	// 3) В случае незавершенности return nil
	// 4) Иначе сохренение сущности на выходе степа в БД, удаление степа из мэпа
	// 5) return err
	return nil
}

func (manager *StepByStepManager) AddCreateQuestionSequenceForUser(user *model.User) {
	// TODO: убедиться что именно так и добавляются значения в Map
	manager.createQuestionSteps[user.Id] = NewCreateQuestionsSequence(user)
}
