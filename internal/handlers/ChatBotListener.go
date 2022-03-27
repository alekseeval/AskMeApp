package handlers

import (
	"AskMeApp/internal/telegramBot"
	"log"
)

func HandleBotMessages(telegramApiToken string) {
	botClient, err := telegramBot.NewBotClient(telegramApiToken)
	if err != nil {
		log.Panic("Не удалось проинициализировать бота", err)
	}
	for update := range botClient.Updates {

		chatId := update.Message.Chat.ID

		switch update.Message.Command() {
		case "/start":
			err = botClient.SendMessage("Это была команда /start", chatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "/help":
			err = botClient.SendMessage("Это была команда /start", chatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "/question":
			err = botClient.SendMessage("Это была команда /start", chatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		case "/chagecategory":
			err = botClient.SendMessage("Это была команда /start", chatId)
			if err != nil {
				log.Panic("Не удалось отправить сообщение", err)
			}
		}
	}
}
