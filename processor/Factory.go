package processor

import (
	"github.com/AlexeyDsov/go-tgbotapi-help/easy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Factory struct {
	api                *tgbotapi.BotAPI
	lastMessageStorage LastMessageStorage
	easyApi            *easy.EasyApi
}

func NewFactory(api *tgbotapi.BotAPI, storage LastMessageStorage) *Factory {
	return &Factory{
		api:                api,
		lastMessageStorage: storage,
		easyApi:            easy.NewEasyApi(api, storage),
	}
}

func (f Factory) EasyApi() *easy.EasyApi {
	return f.easyApi
}

func (f Factory) Api() *tgbotapi.BotAPI {
	return f.api
}

func (f Factory) LastMessageStorage() LastMessageStorage {
	return f.lastMessageStorage
}

func (f Factory) Chain(processors ...Processor) Chain {
	return Chain(processors)
}

func (f Factory) Private(processor Processor) Private {
	return Private{processor: processor}
}

func (f Factory) Message(processor Processor) WithMessage {
	return WithMessage{processor: processor}
}

func (f Factory) CallbackQuery(processor Processor) CallbackQuery {
	return CallbackQuery{processor}
}

func (f Factory) SimpleCommandMap(mapping map[string]Processor) SimpleCommandMap {
	return SimpleCommandMap(mapping)
}

func (f Factory) SimpleCommand(command string, processor Processor) OneWordCommand {
	return OneWordCommand{CommandWordResolverExactWord{command}, processor}
}

func (f Factory) OneWordCommand(resolver CommandWordResolver, processor Processor) OneWordCommand {
	return OneWordCommand{resolver, processor}
}

func (f Factory) LastMessageProcessor(processor Processor) Processor {
	return &LastMessageProcessor{f.lastMessageStorage, processor}
}

func (f *Factory) WithFileMimeType(mimeTpe string, processor Processor) Processor {
	return &WithFileMimeTypes{processor: processor, mimeTypes: []string{mimeTpe}}
}

func (f *Factory) PredefinedReply(replyMessage string) Processor {
	return NewPredefinedReplyProcessor(f.easyApi, replyMessage)
}

func (f *Factory) EmptyPredefined() Processor {
	return &EmptyPredefinedProcessor{}
}