package telegramBot

import (
	"AskMeApp/internal"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"time"
)

func VerifyOrRegisterUser(tgUserInfo *TgBotApi.User, repository internal.UserRepositoryInterface) (*internal.User, error) {
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
