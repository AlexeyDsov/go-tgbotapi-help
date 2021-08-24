package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQuery struct {
	processor Processor
}

func (c CallbackQuery) Process(update tgbotapi.Update, ctx context.Context) bool {
	if update.CallbackQuery != nil {
		return c.processor.Process(update, ctx)
	}
	return false
}