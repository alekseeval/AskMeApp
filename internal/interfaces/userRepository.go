package interfaces

import "AskMeApp/internal/model"

// UserRepository is CRUD repository interface for model.User
type UserRepository interface {
	GetByChatId(telegramChatId int64) (model.User, error)
	Add(user *model.User) error
	Delete(user *model.User) error
	Edit(user *model.User) error
}
