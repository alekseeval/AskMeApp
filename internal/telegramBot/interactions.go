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
	newQuestionEndStep int = iota*(totalNumberOfScenarios+1) + 1
	newQuestionAskCategoryStep
	newQuestionAskAnswerStep
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
	if state.SequenceStep <= 0 {
		state.SequenceStep = nilStep
		state.unfilledNewQuestion = nil
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

	case NewQuestionStartStep:
		msg := tgbotapi.NewMessage(user.TgChatId, "Enter your question:")
		msg.ParseMode = "MarkdownV2"
		_, err := botClient.botApi.Send(msg)
		if err != nil {
			return userState, err
		}
		userState.unfilledNewQuestion = &internal.Question{
			Author: user,
		}
	case newQuestionAskAnswerStep:
		if update.Message.Text == "" {
			return userState, errors.New("empty question title received")
		}
		userState.unfilledNewQuestion.Title = update.Message.Text
		msg := tgbotapi.NewMessage(user.TgChatId, "Enter answer of your question:")
		msg.ParseMode = "MarkdownV2"
		_, err := botClient.botApi.Send(msg)
		if err != nil {
			return userState, err
		}
	case newQuestionAskCategoryStep:
		if update.Message.Text == "" {
			return userState, errors.New("empty question title received")
		}
		userState.unfilledNewQuestion.Answer = update.Message.Text
		allCategories, err := botClient.questionRepository.GetAllCategories()
		if err != nil {
			return userState, err
		}
		err = botClient.SendCategoriesToChooseInline(
			"Choose category of your question:", allCategories, user.TgChatId)
		if err != nil {
			return userState, err
		}
	case newQuestionEndStep:
		callbackData := update.CallbackData()
		if callbackData == "" || callbackData[0] != 'c' {
			allCategories, err := botClient.questionRepository.GetAllCategories()
			if err != nil {
				return userState, err
			}
			err = botClient.SendCategoriesToChooseInline(
				"Again. Choose category of your question:", allCategories, user.TgChatId)
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
		userState.unfilledNewQuestion.Categories = make([]*internal.Category, 0)
		userState.unfilledNewQuestion.Categories = append(userState.unfilledNewQuestion.Categories, category)
		if category.Id != baseCategory.Id {
			userState.unfilledNewQuestion.Categories = append(userState.unfilledNewQuestion.Categories, &baseCategory)
		}
		_, err = botClient.questionRepository.AddQuestion(userState.unfilledNewQuestion)
		if err != nil {
			return userState, err
		}
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Question successfully saved!")
		if _, err := botClient.botApi.Request(callback); err != nil {
			return nil, err
		}
	default:
		return userState, errors.New("unknown step")
	}
	userState = userState.increaseStep()
	return userState, nil
}
