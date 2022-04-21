package repo

import (
	"AskMeApp/internal/model"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// UserRepository is structure for manage work with DB
type UserRepository struct {
	dbConnection *sql.DB
}

func NewUserRepository(host string, port int, user string, password string, dbname string) (*UserRepository, error) {
	psqlConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnString)
	if err != nil {
		return nil, err
	}
	repository := UserRepository{
		dbConnection: db,
	}
	db.SetMaxIdleConns(5)
	return &repository, nil
}

func (r *UserRepository) Add(user *model.User) (*model.User, error) {
	var userId int32
	sqlStatement := `
		INSERT INTO users (first_name, last_name)
		VALUES ($1, $2)
		RETURNING id`
	err := r.dbConnection.QueryRow(sqlStatement, user.FirstName, user.LastName).Scan(&userId)
	if err != nil {
		log.Println(51, err)
		return nil, err
	}
	sqlStatement = `
		INSERT INTO telegram_users (username, chat_id, user_id)
		VALUES ($1, $2, $3)`
	_, err = r.dbConnection.Exec(sqlStatement, user.TgUserName, user.TgChatId, userId)
	if err != nil {
		log.Println(59, err)
		return nil, err
	}
	user.Id = userId
	return user, nil
}

func (r *UserRepository) Delete(user *model.User) error {
	_, err := r.dbConnection.Exec(`DELETE FROM telegram_users WHERE user_id=$1;`, user.Id)
	if err != nil {
		return err
	}
	_, err = r.dbConnection.Exec(`DELETE FROM users WHERE id=$1`, user.Id)
	return err
}

func (r *UserRepository) Edit(user *model.User) error {
	if user.Id <= 0 {
		return errors.New("expected model.User entity with correct id field")
	}
	sqlStatement := `
		UPDATE users
		SET first_name=$2, last_name=$3
		WHERE id=$1`
	_, err := r.dbConnection.Exec(sqlStatement, user.Id, user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	if user.TgChatId != 0 {
		sqlStatement = `
			UPDATE telegram_users
			SET username=$2, chat_id=$3
			WHERE user_id=$1`
		_, err = r.dbConnection.Exec(sqlStatement, user.Id, user.TgUserName, user.TgChatId)
	}
	return err
}

func (r *UserRepository) GetByChatId(telegramChatId int64) (*model.User, error) {
	query := `SELECT u.id, u.first_name, u.last_name, tg.username, tg.chat_id
			  FROM users u LEFT JOIN telegram_users tg ON u.id=tg.user_id
			  WHERE tg.chat_id=$1`
	row := r.dbConnection.QueryRow(query, telegramChatId)

	var id sql.NullInt32
	var firstName sql.NullString
	var lastName sql.NullString
	var tgChatId sql.NullInt64
	var tgUserName sql.NullString

	err := row.Scan(&id, &firstName, &lastName, &tgUserName, &tgChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	usr := model.User{
		Id:         id.Int32,
		FirstName:  firstName.String,
		LastName:   lastName.String,
		TgChatId:   tgChatId.Int64,
		TgUserName: tgUserName.String,
	}
	return &usr, err
}
