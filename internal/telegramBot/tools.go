package telegramBot

import (
	"AskMeApp/internal"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"time"
)

func VerifyOrRegisterUser(tgChatId int64, tgUserInfo *TgBotApi.User, repository internal.UserRepositoryInterface) (*internal.User, error) {
	user, err := repository.GetByChatId(tgChatId)
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
		TgChatId:   tgChatId,
	}
	user, err = repository.Add(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (bot *BotClient) GetTotallyRandomQuestion() (question *internal.Question, err error) {
	allQuestions, err := bot.questionRepository.GetAllQuestions()
	if err != nil {
		return question, err
	}
	if len(allQuestions) == 0 {
		return question, internal.NewZeroQuestionsError()
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(randSource)

	randomNumber := randomizer.Intn(len(allQuestions))
	question = allQuestions[randomNumber]
	return question, err
}
