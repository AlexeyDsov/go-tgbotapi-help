package easy

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type lastMessageStorage interface {
	GetByUser(user *tgbotapi.User) string
	GetByUpdate(update *tgbotapi.Update) string
	Store(user *tgbotapi.User, text string)
	StoreByUpdate(user *tgbotapi.Update, text string)
}

type EasyApi struct {
	*EasyCallback
	*EasyDownloadFile
	*EasyApiKeyboard
	api                *tgbotapi.BotAPI
	lastMessageStorage lastMessageStorage
	NoReplace          bool
	MarkdownMode       string
}

func NewEasyApi(api *tgbotapi.BotAPI, lastMessageStorage lastMessageStorage) *EasyApi {
	return &EasyApi{
		EasyCallback:       &EasyCallback{api},
		EasyDownloadFile:   &EasyDownloadFile{api},
		EasyApiKeyboard:    &EasyApiKeyboard{api, lastMessageStorage},
		api:                api,
		lastMessageStorage: lastMessageStorage,
		MarkdownMode:       tgbotapi.ModeMarkdownV2,
	}
}

func (a *EasyApi) Api() *tgbotapi.BotAPI {
	return a.api
}

func (a *EasyApi) SendEasyMessage(easyMessage *Message) (*tgbotapi.Message, error) {
	text := easyMessage.text
	if easyMessage.oldReplace {
		//@todo не должно быть таких замен в этом коде
		text = strings.Replace(text, "(", "\\(", -1)
		text = strings.Replace(text, ".", "\\.", -1)
		text = strings.Replace(text, ")", "\\)", -1)
		text = strings.Replace(text, "_", "\\_", -1)
		text = strings.Replace(text, "-", "\\-", -1)
	}

	message, err := a.NewMessage(*easyMessage.update, text)
	if err != nil {
		return nil, err
	}
	if len(easyMessage.ParseMode) > 0 {
		message.ParseMode = easyMessage.ParseMode
	}

	messageResult, err := a.SendChattableCommand(*easyMessage.update, message, easyMessage.command)

	return &messageResult, err
}

func (a *EasyApi) Send(chattable tgbotapi.Chattable) (tgbotapi.Message, error) {
	return a.api.Send(chattable)
}

func (a *EasyApi) SendMessage(update tgbotapi.Update, text string) (tgbotapi.Message, error) {
	return a.SendMessageCommand(update, text, "")
}

func (a *EasyApi) SendMessageCommand(update tgbotapi.Update, text string, command string) (tgbotapi.Message, error) {
	easyMessage :=
		NewMessage(&update, text).
			WithCommand(command).
			WithOldReplace().
			WithParseMk2()

	resultMessage, err := a.SendEasyMessage(easyMessage)

	if resultMessage == nil {
		resultMessage = &tgbotapi.Message{}
	}

	return *resultMessage, err
	////@todo не должно быть таких замен в этом коде
	//text = strings.Replace(text, "(", "\\(", -1)
	//text = strings.Replace(text, ".", "\\.", -1)
	//text = strings.Replace(text, ")", "\\)", -1)
	//
	//newMessage, err := a.NewMessageMk2(update, text)
	//if err != nil {
	//	return tgbotapi.Message{}, err
	//}
	//
	//return a.SendChattableCommand(update, newMessage, command)
}

func (a *EasyApi) NewMessageMk2(update tgbotapi.Update, text string) (*tgbotapi.MessageConfig, error) {
	message, err := a.NewMessage(update, text)
	if err == nil {
		message.ParseMode = tgbotapi.ModeMarkdownV2
	}

	return message, err
}

func (a *EasyApi) NewMessage(update tgbotapi.Update, text string) (*tgbotapi.MessageConfig, error) {
	message := a.resolverMessage(update)

	if message == nil {
		return nil, errors.New("no message in update")
	}
	if message.Chat == nil {
		return nil, errors.New("no chat in message")
	}

	newMessage := tgbotapi.NewMessage(message.Chat.ID, text)
	return &newMessage, nil
}

func (a *EasyApi) SendChattable(update tgbotapi.Update, chattable tgbotapi.Chattable) (tgbotapi.Message, error) {
	return a.SendChattableCommand(update, chattable, "")
}

func (a *EasyApi) SendChattableCommand(
	update tgbotapi.Update,
	chattable tgbotapi.Chattable,
	command string,
) (tgbotapi.Message, error) {
	sentMessage, err := a.api.Send(chattable)
	if err == nil {
		isPrivate := sentMessage.Chat.IsPrivate()
		if fromUser := a.resolveFromUser(update); fromUser != nil && isPrivate {
			a.lastMessageStorage.Store(fromUser, command)
		}
	}

	return sentMessage, err
}

func (a *EasyApi) resolverMessage(update tgbotapi.Update) *tgbotapi.Message {
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		return update.CallbackQuery.Message
	} else if update.Message != nil {
		return update.Message
	}
	return nil
}

func (a *EasyApi) resolveFromUser(update tgbotapi.Update) *tgbotapi.User {
	if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
		return update.CallbackQuery.From
	} else if update.Message != nil && update.Message.From != nil {
		return update.Message.From
	}
	return nil
}
