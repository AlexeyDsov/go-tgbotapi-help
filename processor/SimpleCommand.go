package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SimpleCommand struct {
	command   string
	processor Processor
}

func (c SimpleCommand) Process(update tgbotapi.Update, ctx context.Context) bool {
	command, other := EasyCommandSplit{}.Split(update, ctx)

	if command != c.command {
		return false
	}

	return c.processor.Process(update, context.WithValue(ctx, "command", other))
}

type SimpleCommandCtx struct {
	context.Context
}

func (s SimpleCommandCtx) Command() string {
	if command, ok := s.Value("command").(string); ok {
		return command
	} else {
		return ""
	}
}