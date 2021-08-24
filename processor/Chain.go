package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Chain []Processor

func (c Chain) Process(update tgbotapi.Update, ctx context.Context) bool {
	for _, processor := range c {
		if processor.Process(update, ctx) {
			return true
		}
	}

	return false
}