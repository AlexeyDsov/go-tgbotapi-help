package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WithFileMimeTypes struct {
	processor Processor
	mimeTypes []string
}

func (w *WithFileMimeTypes) Process(update tgbotapi.Update, ctx context.Context) bool {
	if update.Message.Document == nil {
		return false
	}

	mimeType := update.Message.Document.MimeType
	for _, expectedMimeType := range w.mimeTypes {
		if expectedMimeType == mimeType {
			return w.processor.Process(update, ctx)
		}
	}

	return false
}

