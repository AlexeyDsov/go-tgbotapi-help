package processor

import (
	"context"
	"github.com/AlexeyDsov/go-tgbotapi-help/easy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type PredefinedReplyProcessor struct {
	easyApi *easy.EasyApi
	replyMessage string
}

func NewPredefinedReplyProcessor(easyApi *easy.EasyApi, replyMessage string) *PredefinedReplyProcessor {
	return &PredefinedReplyProcessor{easyApi: easyApi, replyMessage: replyMessage}
}

func (p *PredefinedReplyProcessor) Process(update tgbotapi.Update, ctx context.Context) bool {
	message, err := p.easyApi.NewMessageMk2(update, p.replyMessage)
	if err != nil {
		log.Printf("PredefinedReplyProcessor error: %s", err)
		return true
	}

	_, err = p.easyApi.SendChattable(update, message)
	if err != nil {
		log.Printf("PredefinedReplyProcessor error: %s", err)
	}

	return true
}
