package main

import (
	"AskMeApp/internal/telegramBot"
	"AskMeApp/internal/telegramBot/handlers"
	"AskMeApp/repo"
	"log"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "askmeapp"
	password = "askmeapp"
	dbname   = "askmeapp"
)

func main() {
	botClient, err := telegramBot.NewBotClient(os.Getenv("ASK_ME_APP_TG_TOKEN"))
	if err != nil {
		log.Panic("Не удалось проинициализировать бота --- ", err)
	}
	userRepo, err := repo.NewUserRepository(host, port, user, password, dbname)
	if err != nil {
		log.Panic("Не удалось установить соединение с БД", err)
	}
	handlers.HandleBotMessages(botClient, userRepo)
}
