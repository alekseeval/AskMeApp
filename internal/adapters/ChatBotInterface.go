package adapters

type ChatBotInterface interface {
	//MessageSender
}

type Creatable interface {
	NewBotClient(MessengerApiToken string) (ChatBotInterface, error)
}

//	TODO: Заполнить методами изменения inline-клавиатуры
type KeyboardChanger interface {
}

type MessageSender interface {
	SendMessage(text string) error
}
