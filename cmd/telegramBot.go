package main

import (
	"AskMeApp/internal/repo"
	"AskMeApp/internal/telegramBot"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		err := botClient.Run()
		if err != nil {
			log.Panic(err)
		}
	}()
	<-exitCh
	err = botClient.Shutdown(time.Second * 15)
	if err != nil {
		log.Panic("Не удалось завершить работу бота", err)
	}
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
