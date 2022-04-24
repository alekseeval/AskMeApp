package repo

import (
	"AskMeApp/internal/model"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (repo *QuestionRepository) Add(question *model.Question) (*model.Question, error) {
	if question.Author.Id <= 0 {
		return nil, errors.New("user is unregistered")
	}
	if question.Category.Id <= 0 {
		return nil, errors.New("unknown question category (haven't id)")
	}
	sqlStatement := `INSERT INTO questions (title, answer) VALUES ($1, $2) RETURNING id`
	row := repo.db.QueryRow(sqlStatement, question.Title, question.Answer)
	err := row.Scan(&question.Id)
	if err != nil {
		return nil, err
	}
	sqlStatement = `INSERT INTO questions2categories VALUES ($1, $2)`
	_, err = repo.db.Exec(sqlStatement, question.Id, question.Category.Id)
	if err != nil {
		return nil, err
	}
	sqlStatement = `INSERT INTO users2questions VALUES ($1, $2)`
	_, err = repo.db.Exec(sqlStatement, question.Author.Id, question.Id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (repo *QuestionRepository) Delete(question *model.Question) error {
	sqlStatement := `DELETE FROM questions2categories WHERE question_id=$1`
	_, err := repo.db.Exec(sqlStatement, question.Id)
	if err != nil {
		return err
	}

	sqlStatement = `DELETE FROM users2questions WHERE question_id=$1`
	_, err = repo.db.Exec(sqlStatement, question.Id)
	if err != nil {
		return err
	}

	sqlStatement = `DELETE FROM questions WHERE id=$1`
	_, err = repo.db.Exec(sqlStatement, question.Id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *QuestionRepository) Edit(question *model.Question) error {
	if question.Author.Id <= 0 {
		return errors.New("user is unregistered")
	}
	if question.Category.Id <= 0 {
		return errors.New("unknown question category (haven't id)")
	}
	sqlStatement := `
		UPDATE questions2categories
		SET category_id=$1
		WHERE question_id=$2`
	_, err := repo.db.Exec(sqlStatement, question.Category.Id, question.Id)
	if err != nil {
		return err
	}
	sqlStatement = `
		UPDATE users2questions
		SET user_id=$1
		WHERE question_id=$2`
	_, err = repo.db.Exec(sqlStatement, question.Author.Id, question.Id)
	if err != nil {
		return err
	}
	sqlStatement = `
		UPDATE questions
		SET title=$1, answer=$2
		WHERE id=$3`
	_, err = repo.db.Exec(sqlStatement, question.Title, question.Answer, question.Id)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Пересмотреть смысл получения всех данных о пользователях и категориях
//  или оптимизировать формирование сущностей без дублирования
func (repo *QuestionRepository) GetAll() ([]model.Question, error) {
	sqlStatement := `
			SELECT q.id, q.title, q.answer,
			       category_id, qc.title category_title,
			       u.id user_id, u.first_name user_first_name, u.last_name user_last_name, tu.chat_id user_chat_id, tu.username user_tg_username
			FROM questions q
			    LEFT JOIN questions2categories q2c on q.id = q2c.question_id
				LEFT JOIN users2questions u2q on q.id = u2q.question_id
				LEFT JOIN users u on u2q.user_id = u.id
				LEFT JOIN question_categories qc on q2c.category_id = qc.id
				LEFT JOIN telegram_users tu on u.id = tu.user_id`
	rows, err := repo.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	questions := make([]model.Question, 0)
	for rows.Next() {
		var id sql.NullInt32
		var title sql.NullString
		var answer sql.NullString
		var categoryId sql.NullInt32
		var categoryTitle sql.NullString
		var authorId sql.NullInt32
		var authorFirstName sql.NullString
		var authorLastName sql.NullString
		var authorTgChatId sql.NullInt64
		var authorTgUserName sql.NullString
		err = rows.Scan(&id, &title, &answer, &categoryId, &categoryTitle, &authorId, &authorFirstName, &authorLastName, &authorTgChatId, &authorTgUserName)
		if err != nil {
			return nil, err
		}
		q := model.Question{
			Id:     id.Int32,
			Title:  title.String,
			Answer: answer.String,
			Category: &model.Category{
				Id:    categoryId.Int32,
				Title: categoryTitle.String,
			},
			Author: &model.User{
				Id:         authorId.Int32,
				FirstName:  authorFirstName.String,
				LastName:   authorLastName.String,
				TgChatId:   authorTgChatId.Int64,
				TgUserName: authorTgUserName.String,
			},
		}
		questions = append(questions, q)
	}
	return questions, nil
}
