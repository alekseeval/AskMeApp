package dbManager

type TelegramUserData struct {
	Id       int
	Username string
	ChatId   int64
	UserId   int
}

type User struct {
	Id         int
	FirstName  string
	MiddleName string
	SecondName string
}

type Question struct {
	Id     int
	Title  string
	Answer string
}

type QuestionCategory struct {
	Id    int
	Title string
}
