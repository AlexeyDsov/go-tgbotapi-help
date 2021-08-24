package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type EasyCommandSplit struct{}

func (e EasyCommandSplit) Split(update tgbotapi.Update, ctx context.Context) (command string, other string) {
	return e.SplitByText(e.text(update, ctx))
}

func (e EasyCommandSplit) SplitByText(text string) (string, string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) == 0 {
		return "", ""
	} else if len(parts) == 1 {
		return parts[0], ""
	} else {
		return parts[0], parts[1]
	}
}

func (e EasyCommandSplit) text(update tgbotapi.Update, ctx context.Context) string {
	if command, found := ctx.Value("command").(string); found {
		return command
	} else if message := update.Message; message != nil {
		return message.Text
	} else if cq := update.CallbackQuery; cq != nil {
		return cq.Data
	}
	return ""
}