package telegramBot

import (
	"AskMeApp/internal"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"sync"
)

const nilStep int = 0
const totalNumberOfScenarios = 2

const (
	// TODO: реализовать добавление задачи по шагам
	newQuestionEndStep int = iota*(totalNumberOfScenarios+1) + 1
	newQuestion2Step
	newQuestion3Step
	NewQuestionStartStep
)

const (
	changeCategoryEndStep int = iota*(totalNumberOfScenarios+1) + 2
	ChangeCategoryInitStep
)

type userState struct {
	CurrentCategory internal.Category
	SequenceStep    int

	unfilledNewQuestion *internal.Question

	mutex sync.Mutex
}

func NewUserState(currentCategory internal.Category) *userState {
	return &userState{
		CurrentCategory: currentCategory,
		SequenceStep:    nilStep,
		mutex:           sync.Mutex{},
	}
}

func (state *userState) increaseStep() *userState {
	state.SequenceStep -= totalNumberOfScenarios + 1
	if state.SequenceStep < 0 {
		state.SequenceStep = nilStep
	}
	return state
}

func (botClient *BotClient) ProcessUserStep(user *internal.User, userState *userState, update *tgbotapi.Update) (*userState, error) {
	switch userState.SequenceStep {
	case ChangeCategoryInitStep:
		allCategories, err := botClient.questionRepository.GetAllCategories()
		if err != nil {
			return userState, err
		}
		err = botClient.SendCategoriesToChooseInline(
			"Выберите желаемую категорию вопросов:", allCategories, user.TgChatId)
		if err != nil {
			return userState, err
		}
	case changeCategoryEndStep:
		callbackData := update.CallbackData()
		if callbackData == "" || callbackData[0] != 'c' {
			allCategories, err := botClient.questionRepository.GetAllCategories()
			if err != nil {
				return userState, err
			}
			err = botClient.SendCategoriesToChooseInline(
				"Все-таки выберите желаемую категорию вопросов:", allCategories, user.TgChatId)
			if err != nil {
				return userState, err
			}
			return userState, nil
		}
		categoryId, err := strconv.ParseInt(callbackData[1:len(callbackData)], 10, 32)
		if err != nil {
			return userState, nil
		}
		category, err := botClient.questionRepository.GetCategoryById(int32(categoryId))
		if err != nil {
			return userState, err
		}
		userState.CurrentCategory = *category

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Категория успешно изменена")
		if _, err := botClient.botApi.Request(callback); err != nil {
			return nil, err
		}

		msg := tgbotapi.NewMessage(user.TgChatId, "Категория успешна изменена на: __"+category.Title+"__")
		msg.ParseMode = "MarkdownV2"
		_, err = botClient.botApi.Send(msg)
		if err != nil {
			return userState, err
		}

	//	TODO: Question steps
	default:
		return userState, errors.New("unknown step")
	}
	userState = userState.increaseStep()
	return userState, nil
}
