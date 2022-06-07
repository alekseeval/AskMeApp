package telegramBot

import (
	"AskMeApp/internal"
	TgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type BotClient struct {
	bot     *TgBotApi.BotAPI
	updates TgBotApi.UpdatesChannel
	wg      *sync.WaitGroup

	userRepository     internal.UserRepositoryInterface
	questionRepository internal.QuestionsRepositoryInterface
}

func NewBotClient(userRepository internal.UserRepositoryInterface, questionRepository internal.QuestionsRepositoryInterface, botToken string) (*BotClient, error) {
	bot, err := TgBotApi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	updatesConfig := TgBotApi.NewUpdate(0)
	updatesConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updatesConfig)
	return &BotClient{
		bot:     bot,
		updates: updates,

		userRepository:     userRepository,
		questionRepository: questionRepository,
	}, nil
}

// TODO: Заменить WaitGroup на Context
func (bot *BotClient) Run() {
	bot.wg = &sync.WaitGroup{}
	bot.wg.Add(1)
	go bot.handleBotUpdates()
}

func (bot *BotClient) Stop() {
	bot.wg.Done()
}
