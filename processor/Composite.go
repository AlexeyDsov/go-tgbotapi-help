package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Composite []Processor

func (c Composite) Process(update tgbotapi.Update, ctx context.Context) bool {
	result := false
	for _, processor := range c {
		result = processor.Process(update, ctx) || result
	}

	return result
}