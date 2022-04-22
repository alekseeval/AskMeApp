package repo

import (
	"AskMeApp/internal/model"
	"database/sql"
)

type QuestionRepository struct {
	db *sql.DB
}

func (repo *QuestionRepository) Add(question *model.Question) (*model.Question, error) {
	return nil, nil
}

func (repo *QuestionRepository) Delete(question *model.Question) error {
	return nil
}

func (repo *QuestionRepository) Edit(question *model.Question) error {
	return nil
}

func (repo *QuestionRepository) GetByAuthor(user *model.User) (*model.Question, error) {
	return nil, nil
}
