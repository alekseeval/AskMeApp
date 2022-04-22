package model

type User struct {
	Id int32 `json:"id"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	TgChatId   int64  `json:"tgChatId"`
	TgUserName string `json:"tgUserName"`
}

type Category struct {
	Id    int32
	Title string
}

type Question struct {
	Id     int32  `json:"id"`
	Title  string `json:"title"`
	Answer string `json:"answer"`

	Category *Category `json:"category"`
	AuthorId *User     `json:"authorId"`
}
