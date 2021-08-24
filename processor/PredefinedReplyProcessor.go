package processor

import (
	"context"
	"github.com/AlexeyDsov/go-tgbotapi-help/easy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PredefinedReplyProcessor struct {
	easyApi *easy.EasyApi
	replyMessage string
}

func NewPredefinedReplyProcessor(easyApi *easy.EasyApi, replyMessage string) *PredefinedReplyProcessor {
	return &PredefinedReplyProcessor{easyApi: easyApi, replyMessage: replyMessage}
}

func (p *PredefinedReplyProcessor) Process(update tgbotapi.Update, ctx context.Context) bool {
	p.easyApi.SendMessage(update, p.replyMessage)

	return true
}
