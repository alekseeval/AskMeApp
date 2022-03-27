package main

import (
	"AskMeApp/internal/handlers"
	"os"
)

func main() {
	handlers.HandleBotMessages(os.Getenv("ASK_ME_APP_TG_TOKEN"))
}
