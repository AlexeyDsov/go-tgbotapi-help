package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAddedToChat struct {
	proccessor Processor
	bot *tgbotapi.BotAPI
}

func NewBotAddedToChat(proccessor Processor, bot *tgbotapi.BotAPI) *BotAddedToChat {
	return &BotAddedToChat{proccessor: proccessor, bot: bot}
}

func (a *BotAddedToChat) Process(update tgbotapi.Update, ctx context.Context) bool {
	if update.Message.NewChatMembers == nil || len(update.Message.NewChatMembers) == 0 {
		return false
	}

	for _, user := range update.Message.NewChatMembers {
		if user.ID == a.bot.Self.ID {
			return a.proccessor.Process(update, ctx)
		}
	}

	return false
}
