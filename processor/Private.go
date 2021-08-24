package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Private struct {
	processor Processor
}

func (p Private) Process(update tgbotapi.Update, ctx context.Context) bool {
	message := update.Message
	if message == nil {
		return false
	}

	chat := message.Chat
	if chat == nil {
		return false
	}

	if chat.Type != "private" {
		return false
	}

	return p.processor.Process(update, ctx)
}