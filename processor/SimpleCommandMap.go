package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type SimpleCommandMap map[string]Processor

func (p SimpleCommandMap) Process(update tgbotapi.Update, ctx context.Context) bool {
	parts, found := newSimpleCommandParts(update)
	if !found {
		return false
	}

	command := parts.command()

	for expectedCommand, processor := range p {
		if strings.ToLower(expectedCommand) == strings.ToLower(command) {
			return processor.Process(update, context.WithValue(ctx, "command_args", parts.message()))
		}
	}

	return false
}

type simpleCommandParts []string

func newSimpleCommandParts(update tgbotapi.Update) (simpleCommandParts, bool) {
	text := update.Message.Text
	if !strings.HasPrefix(text, "/") {
		return nil, false
	}

	return strings.SplitN(text, " ", 2), true
}

func (s simpleCommandParts) command() string {
	return s[0][1:]
}

func (s simpleCommandParts) message() string {
	if len(s) > 1 {
		return s[1]
	}
	return ""
}