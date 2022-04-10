package repo

import "AskMeApp/internal/model"

// UserRepository is structure for manage work with DB
// TODO: реализовать основные методы из интерфейса interfaces.UserRepository
type UserRepository struct {
	Users []model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Users: make([]model.User, 5),
	}
}

func (r UserRepository) GetByChatId(telegramChatId int64) (*model.User, error) {
	return nil, nil
}

func (r UserRepository) Add(user *model.User) error {
	return nil
}

func (r UserRepository) Delete(user *model.User) error {
	return nil
}

func (r UserRepository) Edit(user *model.User) error {
	return nil
}
