package telegramBot

import (
	"AskMeApp/internal/interfaces"
	"AskMeApp/internal/model"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func HandleBotMessages(botClient *BotClient, userRepo interfaces.UserRepositoryInterface) {
	for update := range botClient.Updates {

		if update.Message != nil {

			user, err := VerifyOrRegisterUser(update.Message.Chat.ID, update.Message.From, userRepo)
			if err != nil {
				err = botClient.SendTextMessage("Что-то пошло не так во время авторизации: \n"+err.Error(), update.Message.Chat.ID)
				if err != nil {
					log.Panic("Жопа наступила, не удалось получить или создать юзера,"+
						" а потом еще и сообщение не отправилось", err)
				}
			}

			switch update.Message.Command() {
			case "start":
				err = botClient.SendTextMessage("Это была команда /start", user.TgChatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "help":
				err = botClient.SendTextMessage("Это была команда /help", user.TgChatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "question":
				err = botClient.SendTextMessage("Это была команда /question", user.TgChatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "changecategory":
				err = botClient.SendTextMessage("Это была команда /changecategory", user.TgChatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "stop":
				if user.TgUserName != "al_andrew" {
					continue
				}
				err = botClient.SendTextMessage("Приложение завершило свою работу", user.TgChatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
				botClient.Stop()
				return
			}
		}
	}
}

func VerifyOrRegisterUser(tgChatId int64, tgUserInfo *TgBotApi.User, repository interfaces.UserRepositoryInterface) (*model.User, error) {
	user, err := repository.GetByChatId(tgChatId)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user = &model.User{
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
