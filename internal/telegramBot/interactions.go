package telegramBot

import (
	"AskMeApp/internal"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"sync"
)

// nilStep 					- нулевой шаг последовательности, означает, что никакой последовательности не выполняется
// totalNumberOfScenarios 	- кол-во сценариев бота
// stepRange 				- расстояние между двумя шагами одной последовательности
const (
	nilStep                int = 0
	totalNumberOfScenarios int = 2
	stepRange              int = totalNumberOfScenarios + 1 // 3
)

// Нумерация шагов последовательности действий по добавлению нового вопроса в базу знаний
// NewQuestionInitStep - шаг начала последовательности
const (
	newQuestionEndStep int = iota*stepRange + 1
	newQuestionAskCategoryStep
	newQuestionAskAnswerStep
	NewQuestionInitStep
)

// Нумерация шагов последовательности действий по смене категории вопросов
// ChangeCategoryInitStep - шаг начала последовательности
const (
	changeCategoryEndStep int = iota*stepRange + 2
	ChangeCategoryInitStep
)

type userState struct {
	CurrentCategory internal.Category
	SequenceStep    int

	unfilledNewQuestion *internal.Question

	mutex sync.Mutex
}

// NewUserState - метод-конструктор, возвращающий userState без начатой последовательности действий и
//  с базовой категорией currentCategory в качестве текущей
func NewUserState(currentCategory internal.Category) *userState {
	return &userState{
		CurrentCategory: currentCategory,
		SequenceStep:    nilStep,
		mutex:           sync.Mutex{},
	}
}

// increaseStep - переводит последовательность действий на следующий шаг
func (state *userState) increaseStep() *userState {
	state.SequenceStep -= stepRange
	if state.SequenceStep <= 0 {
		state.SequenceStep = nilStep
		state.unfilledNewQuestion = nil
	}
	return state
}

// dropStepsProgress - обнуляет прогресс по текущей последовательности действий, без сохранения результата
func (state *userState) dropStepsProgress() *userState {
	state.SequenceStep = nilStep
	state.unfilledNewQuestion = nil
	return state
}

func (botClient *BotClient) ProcessUserStep(user *internal.User, userState *userState, update *tgbotapi.Update) (*userState, error) {
	if update.Message != nil && update.Message.Text == cancelAllStepsCommandText {
		userState = userState.dropStepsProgress()
		msg := tgbotapi.NewMessage(user.TgChatId, "Action cancelled")
		msg.ReplyMarkup = MainKeyboard
		_, err := botClient.botApi.Send(msg)
		return userState, err
	}
	switch userState.SequenceStep {
	case ChangeCategoryInitStep:
		allCategories, err := botClient.questionRepository.GetAllCategories()
		if err != nil {
			return userState, err
		}
		msg := tgbotapi.NewMessage(user.TgChatId, "Выберите желаемую категорию вопросов:")
		msg.ReplyMarkup = formatCategoriesToInlineMarkup(allCategories)
		_, err = botClient.botApi.Send(msg)
		if err != nil {
			return userState, err
		}
		msg.Text = "Now is __" + userState.CurrentCategory.Title + "__"
		msg.ParseMode = "MarkdownV2"
		msg.ReplyMarkup = KeyboardWithCancel
		_, err = botClient.botApi.Send(msg)
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
			msg := tgbotapi.NewMessage(user.TgChatId, "Все-таки выберите желаемую категорию вопросов: ")
			msg.ReplyMarkup = formatCategoriesToInlineMarkup(allCategories)
			_, err = botClient.botApi.Send(msg)
			if err != nil {
				return userState, err
			}
			msg.Text = "Now is __" + userState.CurrentCategory.Title + "__"
			msg.ParseMode = "MarkdownV2"
			msg.ReplyMarkup = nil
			_, err = botClient.botApi.Send(msg)
			if err != nil {
				return userState, err
			}
			return userState, nil
		}
		categoryId, err := strconv.ParseInt(callbackData[1:], 10, 32)
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
		msg.ReplyMarkup = MainKeyboard
		_, err = botClient.botApi.Send(msg)
		if err != nil {
			return userState, err
		}

	case NewQuestionInitStep:
		msg := tgbotapi.NewMessage(user.TgChatId, "Enter your question:")
		msg.ParseMode = "MarkdownV2"
		msg.ReplyMarkup = KeyboardWithCancel
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
		msg := tgbotapi.NewMessage(user.TgChatId, "Choose category of your question:")
		msg.ReplyMarkup = formatCategoriesToInlineMarkup(allCategories)
		_, err = botClient.botApi.Send(msg)
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
			msg := tgbotapi.NewMessage(user.TgChatId, "Again. Choose category of your question:")
			msg.ReplyMarkup = formatCategoriesToInlineMarkup(allCategories)
			_, err = botClient.botApi.Send(msg)
			if err != nil {
				return userState, err
			}
			return userState, nil
		}
		categoryId, err := strconv.ParseInt(callbackData[1:], 10, 32)
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
		msg := tgbotapi.NewMessage(user.TgChatId, "Question successfully saved!")
		msg.ReplyMarkup = MainKeyboard
		_, err = botClient.botApi.Send(msg)
		if err != nil {
			return userState, err
		}
	default:
		return userState, errors.New("unknown step")
	}
	userState = userState.increaseStep()
	return userState, nil
}
