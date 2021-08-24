package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EmptyPredefinedProcessor struct {

}

func (e EmptyPredefinedProcessor) Process(update tgbotapi.Update, ctx context.Context) bool {
	return true
}
