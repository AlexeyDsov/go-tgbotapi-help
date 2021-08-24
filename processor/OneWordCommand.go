package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandWordResolver interface {
	Resolve(word string, ctx context.Context) ResolveResult
}

type CommandWordResolverExactWord struct {
	word string
}

func (c CommandWordResolverExactWord) Resolve(word string, ctx context.Context) ResolveResult {
	if word == c.word {
		return NewResolverResultOk(ctx)
	} else {
		return NewResolverResultNok()
	}
}

type ResolveResult struct {
	Ok bool
	ForbidNext bool
	Context context.Context
}

func NewResolverResultOk(ctx context.Context) ResolveResult {
	return ResolveResult{Ok: true, ForbidNext: false, Context: ctx}
}

func NewResolverResultNok() ResolveResult {
	return ResolveResult{Ok: false, ForbidNext: false, Context: nil}
}

func NewResolverResultStop() ResolveResult {
	return ResolveResult{Ok: false, ForbidNext: true, Context: nil}
}

type OneWordCommand struct {
	resolver  CommandWordResolver
	processor Processor
}

func (c OneWordCommand) Process(update tgbotapi.Update, ctx context.Context) bool {
	command, other := EasyCommandSplit{}.Split(update, ctx)

	if resolveResult := c.resolver.Resolve(command, ctx); resolveResult.Ok {
		return c.processor.Process(update, context.WithValue(resolveResult.Context, "command", other))
	} else if resolveResult.ForbidNext {
		return true
	}

	return false
}

