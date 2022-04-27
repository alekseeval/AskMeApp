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
	AddQuestion(question *model.Question) (*model.Question, error)
	DeleteQuestion(question *model.Question) error
	EditQuestion(question *model.Question) error
	GetAllQuestions() ([]*model.Question, error)

	CategoriesRepositoryInterface
}

// CategoriesRepositoryInterface is CRUD repository interface for model.Category
type CategoriesRepositoryInterface interface {
	AddCategory(category *model.Category) (*model.Category, error)
	DeleteCategory(category *model.Category) error
	EditCategory(category *model.Category) error

	GetAllCategories() ([]*model.Category, error)
}
