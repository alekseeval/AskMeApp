package main

import (
	"AskMeApp/internal/telegramBot/client"
	"AskMeApp/repo"
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
		log.Panic("Не удалось роинициалихировать БД ", err)
	}
	questionRepository := repo.NewQuestionRepository(db)
	//userRepository := repo.NewUserRepository(db)

	token := os.Getenv("ASK_ME_APP_TG_TOKEN")
	log.Println(token)
	botClient, err := client.NewBotClient(token)
	if err != nil {
		log.Panic(err)
	}
	categories, err := questionRepository.GetAllCategories()
	if err != nil {
		log.Panic("Не получены категории из БД", err)
	}
	err = botClient.SendTextMessage("Просто привет, просто как дела?", 405658316)
	fmt.Println(categories)
	err = botClient.SendInlineCategories("Выберие категорию:", categories, 405658316)
	fmt.Println(categories)
	if err != nil {
		log.Panic("Что-то не так с отправкой сообщения", err)
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
