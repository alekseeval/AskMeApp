package model

type User struct {
	Id         int32  `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	TgChatId   int64  `json:"tgChatId"`
	TgUserName string `json:"tgUserName"`
}
