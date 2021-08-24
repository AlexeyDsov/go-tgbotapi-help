package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotRemovedFromChat struct {
	proccessor Processor
	bot *tgbotapi.BotAPI
}

func NewBotRemovedFromChat(proccessor Processor, bot *tgbotapi.BotAPI) *BotRemovedFromChat {
	return &BotRemovedFromChat{proccessor: proccessor, bot: bot}
}

func (a *BotRemovedFromChat) Process(update tgbotapi.Update, ctx context.Context) bool {
	if update.Message.LeftChatMember == nil {
		return false
	}

	if update.Message.LeftChatMember.ID == a.bot.Self.ID {
		return a.proccessor.Process(update, ctx)
	}

	return false
}
