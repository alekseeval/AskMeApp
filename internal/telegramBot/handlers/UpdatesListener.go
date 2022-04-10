package handlers

import (
	"AskMeApp/internal/telegramBot"
	"log"
)

func HandleBotMessages(botClient *telegramBot.BotClient) {
	for update := range botClient.Updates {

		if update.Message != nil {
			chatId := update.Message.Chat.ID

			switch update.Message.Command() {
			case "start":
				err := botClient.SendMessage("Это была команда /start", chatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "help":
				err := botClient.SendMessage("Это была команда /help", chatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "question":
				err := botClient.SendMessage("Это была команда /question", chatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			case "changecategory":
				err := botClient.SendMessage("Это была команда /changecategory", chatId)
				if err != nil {
					log.Panic("Не удалось отправить сообщение", err)
				}
			}
		}
	}
}
