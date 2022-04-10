package interfaces

import "AskMeApp/internal/model"

// UserRepositoryInterface is CRUD repository interface for model.User
type UserRepositoryInterface interface {
	GetByChatId(telegramChatId int64) (*model.User, error)
	Add(user *model.User) error
	Delete(user *model.User) error
	Edit(user *model.User) error
}
