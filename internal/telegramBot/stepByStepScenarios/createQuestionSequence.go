package stepByStepScenarios

import (
	"AskMeApp/internal/interfaces"
	"AskMeApp/internal/model"
	"AskMeApp/internal/telegramBot/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

const (
	waitForCategory = iota
	waitForTitle
	waitForAnswer
	saveQuestion
	end
)

type CreateQuestionsSequence struct {
	step     int
	question *model.Question
}

func NewCreateQuestionsSequence(user *model.User) *CreateQuestionsSequence {
	return &CreateQuestionsSequence{
		step: waitForCategory,
		question: &model.Question{
			Author: user,
		},
	}
}

func (sequence *CreateQuestionsSequence) doStep(botClient *client.BotClient, update *tgbotapi.Update, questionsRepo interfaces.QuestionsRepositoryInterface) error {
	switch sequence.step {

	case waitForCategory:
		categories, err := questionsRepo.GetAllCategories()
		if err != nil {
			return err
		}
		err = botClient.SendInlineCategories("Select category", categories, sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}

	case waitForTitle:
		if update.CallbackQuery == nil {
			// TODO: Вывести пользователю ошибку
		}
		categoryId, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 32)
		if err != nil {
			botClient.SendTextMessage("Что-то пошло не так, попробуйте еще раз", sequence.question.Author.TgChatId)
			// TODO: Снова выполнить запрос категории
		}
		sequence.question.Category, err = questionsRepo.GetCategoryById(int32(categoryId))
		if err != nil {
			// TODO: Вывести ошибку пользователю
		}
		err = botClient.SendTextMessage("Enter title:", sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}

	case waitForAnswer:
		sequence.question.Title = update.Message.Text
		err := botClient.SendTextMessage("Enter answer:", sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}

	case saveQuestion:
		sequence.question.Answer = update.Message.Text
		_, err := questionsRepo.AddQuestion(sequence.question)
		if err != nil {
			return err
		}
		err = botClient.SendTextMessage("Question successfully saved", sequence.question.Author.TgChatId)
		if err != nil {
			return err
		}
	default:
		log.Println("Smth went wrong with sequence step -", sequence.step)
	}
	sequence.step++
	return nil
}

func (sequence *CreateQuestionsSequence) isDone() bool {
	return sequence.step == end
}
