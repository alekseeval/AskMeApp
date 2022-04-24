package interfaces

import "AskMeApp/internal/model"

// UserRepositoryInterface is CRUD repository interface for model.User
type UserRepositoryInterface interface {
	Add(user *model.User) (*model.User, error)
	Delete(user *model.User) error
	Edit(user *model.User) error

	GetByChatId(telegramChatId int64) (*model.User, error)
}

// QuestionsRepositoryInterface is CRUD repository interface for model.Question
type QuestionsRepositoryInterface interface {
	Add(question *model.Question) (*model.Question, error)
	Delete(question *model.Question) error
	Edit(question *model.Question) error

	GetAll() ([]model.Question, error)
}
