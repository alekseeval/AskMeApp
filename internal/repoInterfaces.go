package internal

// UserRepositoryInterface is CRUD repository interface for models.User
type UserRepositoryInterface interface {
	Add(user *User) (*User, error)
	Delete(user *User) error
	Edit(user *User) error

	GetByChatId(telegramChatId int64) (*User, error)
}

// QuestionsRepositoryInterface is CRUD repository interface for models.Question
type QuestionsRepositoryInterface interface {
	AddQuestion(question *Question) (*Question, error)
	DeleteQuestion(question *Question) error
	EditQuestion(question *Question) error
	GetAllQuestions() ([]*Question, error)

	CategoriesRepositoryInterface
}

// CategoriesRepositoryInterface is CRUD repository interface for models.Category
type CategoriesRepositoryInterface interface {
	AddCategory(category *Category) (*Category, error)
	DeleteCategory(category *Category) error
	EditCategory(category *Category) error

	GetAllCategories() ([]*Category, error)
	GetCategoryById(id int32) (*Category, error)
}
