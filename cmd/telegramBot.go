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

	//question := &internal.Question{
	//	Id:    int32(1),
	//	Title: "Name all SOLID principles",
	//	Answer: "S -\tSingle responsibility principle\n" +
	//		"O -\tOpen/Closed principle\n" +
	//		"L -\tLiskov substitute principle\n" +
	//		"I -\tInterface segregation principle\n" +
	//		"D -\tDependency inversion principle",
	//	Category: &internal.Category{
	//		Id:    1,
	//		Title: "Все вопросы",
	//	},
	//	Author: &internal.User{
	//		Id: 1,
	//	},
	//}
	//err = questionsRepo.EditQuestion(question)
	//if err != nil {
	//	log.Panic("На начале умер", err)
	//}

	botClient, err := telegramBot.NewBotClient(userRepo, questionsRepo, os.Getenv("ASK_ME_APP_TG_TOKEN"))
	if err != nil {
		log.Panic("Не удалось проинициализировать бота --- ", err)
	}
	botClient.Run()
	//botClient.Shutdown()
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
