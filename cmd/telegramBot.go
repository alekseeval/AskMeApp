package main

import (
	"AskMeApp/internal/repo"
	"AskMeApp/internal/telegramBot"
	"database/sql"
	"fmt"
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
	db, err := initDb()
	if err != nil {
		log.Panic("Не удалось проинициализировать БД")
	}
	userRepo := repo.NewUserRepository(db)
	questionsRepo := repo.NewQuestionRepository(db)

	botClient, err := telegramBot.NewBotClient(userRepo, questionsRepo, os.Getenv("ASK_ME_APP_TG_TOKEN"))
	if err != nil {
		log.Panic("Не удалось проинициализировать бота --- ", err)
	}
	botClient.Run()
	botClient.Stop()
}

func initDb() (*sql.DB, error) {
	psqlConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(5)
	return db, err
}
