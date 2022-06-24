package telegramBot

import (
	"AskMeApp/internal"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"sort"
	"time"
)

const (
	OptimalSymbolsNumberInRow float32 = 36
	maxButtonsInLineNumber    float32 = 4
)

var (
	randomizeButtonRow      = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(randomQuestionCommandText))
	changeCategoryButtonRow = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(changeCategoryCommandText))
	addQuestionButtonRow    = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(addQuestionCommandText))
	cancelButtonRow         = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(cancelAllStepsCommandText))

	MainKeyboard       = tgbotapi.NewReplyKeyboard(randomizeButtonRow, changeCategoryButtonRow, addQuestionButtonRow)
	KeyboardWithCancel = tgbotapi.NewReplyKeyboard(cancelButtonRow)
)

func IdentifyOrRegisterUser(tgUserInfo *tgbotapi.User, repository internal.UserRepositoryInterface) (*internal.User, error) {
	user, err := repository.GetByChatId(tgUserInfo.ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user = &internal.User{
		FirstName:  tgUserInfo.FirstName,
		LastName:   tgUserInfo.LastName,
		TgUserName: tgUserInfo.UserName,
		TgChatId:   tgUserInfo.ID,
	}
	user, err = repository.Add(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetRandomQuestion(questions []*internal.Question) (question *internal.Question) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(randSource)
	randomNumber := randomizer.Intn(len(questions))
	question = questions[randomNumber]
	return question
}

// TODO: требуется рефактор и переосмысливание
func formatCategoriesToInlineMarkup(categories []*internal.Category) tgbotapi.InlineKeyboardMarkup {
	categoriesCopy := append([]*internal.Category(nil), categories...)
	sort.Slice(categoriesCopy, func(i, j int) bool {
		return len(categoriesCopy[i].Title) < len(categoriesCopy[j].Title)
	})

	weights := make([]float32, 0)
	for _, c := range categoriesCopy {
		weight := OptimalSymbolsNumberInRow / float32(len(c.Title))
		if weight > maxButtonsInLineNumber {
			weight = maxButtonsInLineNumber
		}
		if weight == maxButtonsInLineNumber {
			weight = maxButtonsInLineNumber
		}
		weights = append(weights, maxButtonsInLineNumber/weight)
	}

	inlineButtons := make([][]tgbotapi.InlineKeyboardButton, 1)
	var rowNumber int
	var weightSum float32
	for i, category := range categoriesCopy {
		weightSum += weights[i]
		if weightSum > maxButtonsInLineNumber {
			newRow := make([]tgbotapi.InlineKeyboardButton, 0)
			if len(inlineButtons[rowNumber]) == 0 {
				inlineButtons[rowNumber] = append(inlineButtons[rowNumber], tgbotapi.NewInlineKeyboardButtonData(category.Title, "c"+fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = 0
			} else {
				newRow = append(newRow, tgbotapi.NewInlineKeyboardButtonData(category.Title, "c"+fmt.Sprint(category.Id)))
				inlineButtons = append(inlineButtons, newRow)
				weightSum = weights[i]
			}
			rowNumber += 1
		} else {
			inlineButtons[rowNumber] = append(inlineButtons[rowNumber], tgbotapi.NewInlineKeyboardButtonData(category.Title, "c"+fmt.Sprint(category.Id)))
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
}
