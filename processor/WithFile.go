package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WithFile struct {
	processor Processor
}

func (w *WithFile) Process(update tgbotapi.Update, ctx context.Context) bool {
	if update.Message.Document == nil {
		return false
	}

	return w.processor.Process(update, ctx)
}

