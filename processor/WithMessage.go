package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WithMessage struct {
	processor Processor
}

func (p WithMessage) Process(update tgbotapi.Update, ctx context.Context) bool {
	if update.Message != nil {
		return p.processor.Process(update, ctx)
	}

	return false
}