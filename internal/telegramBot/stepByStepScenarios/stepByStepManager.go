package stepByStepScenarios

import (
	"AskMeApp/internal/interfaces"
	"AskMeApp/internal/model"
	"AskMeApp/internal/telegramBot/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

// StepByStepManager is struct for managing sequences of user actions
type StepByStepManager struct {
	questionRepo            interfaces.QuestionsRepositoryInterface
	createQuestionSequences map[int32]*CreateQuestionsSequence
	qLock                   sync.Mutex // mutex for save createQuestionSequences map
}

func newStepByStepManager(questionRepo interfaces.QuestionsRepositoryInterface) *StepByStepManager {
	return &StepByStepManager{
		questionRepo:            questionRepo,
		createQuestionSequences: make(map[int32]*CreateQuestionsSequence),
	}
}

// DoUserHaveSequences - checks all possible user actions (now only sequence of creation new model.Question)
func (manager *StepByStepManager) DoUserHaveSequences(user *model.User) bool {
	manager.qLock.Lock()
	_, found := manager.createQuestionSequences[user.Id]
	manager.qLock.Unlock()
	return found
}

// StartNewCreateQuestionSequence - starts new sequence for creation new model.Question
func (manager *StepByStepManager) StartNewCreateQuestionSequence(botClient *client.BotClient, update *tgbotapi.Update, user *model.User) error {
	manager.qLock.Lock()
	sequence := NewCreateQuestionsSequence(user)
	err := sequence.doStep(botClient, update)
	if err != nil {
		return err
	}
	manager.createQuestionSequences[user.Id] = sequence
	manager.qLock.Unlock()
	return nil
}

// DoUserHaveSequences - checks all possible user actions (now only sequence of creation new model.Question)
func (manager *StepByStepManager) dropSequenceForUser(user *model.User) {
	manager.qLock.Lock()
	_, found := manager.createQuestionSequences[user.Id]
	if found {
		delete(manager.createQuestionSequences, user.Id)
	}
	manager.qLock.Unlock()
}

func (manager *StepByStepManager) ExecuteUserSequence(botClient *client.BotClient, user *model.User, update *tgbotapi.Update) error {
	manager.qLock.Lock()
	sequence, found := manager.createQuestionSequences[user.Id]
	if found {
		err := sequence.doStep(botClient, update)
		if err != nil {
			return err
		}
		if sequence.isDone() {
			question := sequence.getEntity()
			_, err := manager.questionRepo.AddQuestion(question)
			delete(manager.createQuestionSequences, user.Id)
			manager.qLock.Unlock()
			return err
		}
		manager.qLock.Unlock()
		return nil
	}
	manager.qLock.Unlock()
	// From here may be any other sequence's execution (same as the previous one for CreateQuestionsSequence)
	return nil
}
