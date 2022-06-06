package repo

import (
	"AskMeApp/internal/model"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

// UserRepository is structure for manage work with DB
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Add(user *model.User) (*model.User, error) {
	tx, err := repo.db.Begin()
	sqlStatement := `
		INSERT INTO users (first_name, last_name)
		VALUES ($1, $2)
		RETURNING id`
	err = tx.QueryRow(sqlStatement, user.FirstName, user.LastName).Scan(&user.Id)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}
	sqlStatement = `
		INSERT INTO telegram_users (username, chat_id, user_id)
		VALUES ($1, $2, $3)`
	_, err = tx.Exec(sqlStatement, user.TgUserName, user.TgChatId, user.Id)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}
	err = tx.Commit()
	return user, err
}

func (repo *UserRepository) Delete(user *model.User) error {
	tx, err := repo.db.Begin()
	_, err = tx.Exec(`DELETE FROM telegram_users WHERE user_id=$1;`, user.Id)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}
		return err
	}
	_, err = tx.Exec(`DELETE FROM users WHERE id=$1`, user.Id)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}
		return err
	}
	err = tx.Commit()
	return err
}

func (repo *UserRepository) Edit(user *model.User) error {
	if user.Id <= 0 {
		return errors.New("expected model.User entity with correct id field")
	}
	tx, err := repo.db.Begin()
	sqlStatement := `
		UPDATE users
		SET first_name=$2, last_name=$3
		WHERE id=$1`
	_, err = tx.Exec(sqlStatement, user.Id, user.FirstName, user.LastName)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}
		return err
	}
	if user.TgChatId != 0 {
		sqlStatement = `
			UPDATE telegram_users
			SET username=$2, chat_id=$3
			WHERE user_id=$1`
		_, err = tx.Exec(sqlStatement, user.Id, user.TgUserName, user.TgChatId)
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return txErr
			}
			return err
		}
	}
	err = tx.Commit()
	return err
}

func (repo *UserRepository) GetByChatId(telegramChatId int64) (*model.User, error) {
	query := `SELECT u.id, u.first_name, u.last_name, tg.username, tg.chat_id
			  FROM users u LEFT JOIN telegram_users tg ON u.id=tg.user_id
			  WHERE tg.chat_id=$1`
	row := repo.db.QueryRow(query, telegramChatId)

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
