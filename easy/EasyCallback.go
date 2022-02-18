package easy

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type EasyCallback struct {
	api *tgbotapi.BotAPI
}

func (a *EasyCallback) AnswerEmptyCallback(update tgbotapi.Update) {
	a.AnswerCallback(update, "")
}

func (a *EasyCallback) AnswerCallback(update tgbotapi.Update, text string) {
	if update.CallbackQuery == nil {
		return
	}

	_, err := a.api.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, text))

	if err != nil {
		log.Printf("error while sending response: %s", err)
	}
}