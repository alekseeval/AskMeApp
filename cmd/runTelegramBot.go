package main

import (
	"AskMeApp/internal/telegramBot"
	"AskMeApp/internal/telegramBot/handlers"
	"AskMeApp/repo"
	"log"
	"os"
)

func main() {
	botClient, err := telegramBot.NewBotClient(os.Getenv("ASK_ME_APP_TG_TOKEN"))
	if err != nil {
		log.Panic("Не удалось проинициализировать бота --- ", err)
	}
	userRepo := repo.NewUserRepository()
	handlers.HandleBotMessages(botClient, userRepo)
}
