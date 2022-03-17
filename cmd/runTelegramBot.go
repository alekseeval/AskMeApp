package main

import (
	"AskMeApp/internal"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	botClient := internal.NewTelegramBotClient(os.Getenv("TELEGRAM_API_TOKEN"))

	defaultText := "Андрей снова будет засыпать весь день, зато мы с ним научились новым приколам :)\n\n" +
		"Напиши команду /start"

	// Чтение канала изменений бота
	for update := range botClient.UpdatesChan {

		if update.Message != nil {
			// Создание базового сообщения
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, defaultText)

			// Проверка ожидаемых от пользователя команд
			switch update.Message.Text {
			case "/start":
				msg.Text = "Теперь тут кастомная клавиатура, жесть🙊🎇"
			case "Верни мне клавиатуру!!!":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				msg.Text = "Ты можешь сделать это самостоятельно, но так и быть, держи\n\n" +
					"P.S. Команда /start все еще в деле ;)"
			case "Пришли мне котячий стикер :3":
				msg.Text = "Андрей слишком хотел спать и не успел с этим разобраться :(((\n\n" +
					"UPD: Разобрался как отправить в качестве стикера файл, но не понял как выбрать из уже готовых.."
			}

			// Отправка сообщения
			_, err := botClient.Bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
			continue
		}
		if update.CallbackQuery != nil {
			// TODO: Разобраться что происходит в InlineKeyboard
			continue
		}
	}
}
