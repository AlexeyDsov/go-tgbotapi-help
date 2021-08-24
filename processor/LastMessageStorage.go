package processor

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LastMessageStorage interface {
	GetByUser(user *tgbotapi.User) string
	GetByUpdate(update *tgbotapi.Update) string
	Store(user *tgbotapi.User, text string)
	StoreByUpdate(user *tgbotapi.Update, text string)
}

type LastMessageRepo struct {
	messages map[string]string
}

func NewLastMessageRepo() *LastMessageRepo {
	return &LastMessageRepo{messages: map[string]string{}}
}

func (r LastMessageRepo) GetByUser(user *tgbotapi.User) string {
	userName := r.username(user)
	if userName == "" {
		return ""
	}

	if command, found := r.messages[userName]; found {
		return command
	} else {
		return ""
	}
}

func (r LastMessageRepo) GetByUpdate(update *tgbotapi.Update) string {
	panic("implement me")
}

func (r LastMessageRepo) Store(user *tgbotapi.User, text string) {
	userName := r.username(user)
	if len(userName) > 0 {
		if len(text) > 0 {
			r.messages[userName] = text
		} else {
			delete(r.messages, userName)
		}
	}
}

func (r LastMessageRepo) StoreByUpdate(update *tgbotapi.Update, text string) {
	user := r.userByUpdate(update)
	if user == nil {
		return
	}

	r.Store(user, text)
}

func (r LastMessageRepo) userByUpdate(update *tgbotapi.Update) *tgbotapi.User {
	if message := update.Message; message != nil {
		return r.userByMessage(message)
	} else if cq := update.CallbackQuery; cq != nil {
		if message := cq.Message; message != nil {
			return r.userByMessage(message)
		}
	}
	return nil
}

func (r LastMessageRepo) userByMessage(message *tgbotapi.Message) *tgbotapi.User {
	return message.From
}

func (l LastMessageRepo) username(user *tgbotapi.User) string {
	if user == nil {
		return ""
	}
	return user.UserName
}

type LastMessageProcessor struct {
	storage   LastMessageStorage
	processor Processor
}

func (p LastMessageProcessor) Process(update tgbotapi.Update, ctx context.Context) bool {
	byUser := p.storage.GetByUser(update.Message.From)
	if len(byUser) > 0 {
		return p.processor.Process(update, context.WithValue(ctx, "command", byUser))
	}
	return false
}