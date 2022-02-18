package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type ChatType struct {
	check     func(chat *tgbotapi.Chat) bool
	processor Processor
}

func NewPrivateChat(processor Processor) *ChatType {
	return &ChatType{check: chatPrivateCheck, processor: processor}
}

func NewAnyGroupChat(processor Processor) *ChatType {
	return &ChatType{check: chatAnyGroupCheck, processor: processor}
}

func NewChannelChat(processor Processor) *ChatType {
	return &ChatType{check: chatChannelCheck, processor: processor}
}

func chatPrivateCheck(chat *tgbotapi.Chat) bool {
	return chat.IsPrivate()
}

func chatAnyGroupCheck(chat *tgbotapi.Chat) bool {
	return chat.IsGroup() || chat.IsSuperGroup()
}

func chatChannelCheck(chat *tgbotapi.Chat) bool {
	return chat.IsChannel()
}

func (p *ChatType) Process(update tgbotapi.Update, ctx context.Context) bool {
	message := update.Message
	if message == nil {
		return false
	}

	chat := message.Chat
	if chat == nil {
		return false
	}

	log.Printf("Chat type: %s", chat.Type)

	if !p.check(chat) {
		return false
	}

	return p.processor.Process(update, ctx)
}
