package easy

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/h2non/gentleman.v2"
	"log"
	"strings"
)

type lastMessageStorage interface {
	GetByUser(user *tgbotapi.User) string
	GetByUpdate(update *tgbotapi.Update) string
	Store(user *tgbotapi.User, text string)
	StoreByUpdate(user *tgbotapi.Update, text string)
}

type EasyApi struct {
	api                *tgbotapi.BotAPI
	lastMessageStorage lastMessageStorage
}

func NewEasyApi(api *tgbotapi.BotAPI, lastMessageStorage lastMessageStorage) *EasyApi {
	return &EasyApi{api: api, lastMessageStorage: lastMessageStorage}
}

func (a EasyApi) AnswerEmptyCallback(update tgbotapi.Update) {
	a.AnswerCallback(update, "")
}

func (a EasyApi) AnswerCallback(update tgbotapi.Update, text string) {
	if update.CallbackQuery == nil {
		return
	}

	_, err := a.api.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, text))

	if err != nil {
		log.Printf("error while sending response: %s", err)
	}
}

func (a EasyApi) SendMessage(update tgbotapi.Update, text string) (tgbotapi.Message, error) {
	return a.SendMessageCommand(update, text, "")
}

func (a EasyApi) SendMessageCommand(update tgbotapi.Update, text string, command string) (tgbotapi.Message, error) {
	var message *tgbotapi.Message
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		message = update.CallbackQuery.Message
	} else if update.Message != nil {
		message = update.Message
	}

	if message == nil {
		return tgbotapi.Message{}, errors.New("no message in update")
	}
	if message.Chat == nil {
		return tgbotapi.Message{}, errors.New("no chat in message")
	}

	text = strings.Replace(text, "(", "\\(", -1)
	text = strings.Replace(text, ".", "\\.", -1)
	text = strings.Replace(text, ")", "\\)", -1)
	newMessage := tgbotapi.NewMessage(message.Chat.ID, text)
	newMessage.ParseMode = tgbotapi.ModeMarkdownV2

	sentMessage, err := a.api.Send(newMessage)
	if err == nil {
		var fromUser *tgbotapi.User
		if update.CallbackQuery != nil && update.CallbackQuery.From != nil{
			fromUser = update.CallbackQuery.From
		} else if update.Message != nil && update.Message.From != nil {
			fromUser = update.Message.From
		}
		a.lastMessageStorage.Store(fromUser, command)
	}

	return sentMessage, err
}

func (a EasyApi) SendKeyboard(update tgbotapi.Update, keyboard EasyKeyboard, text string, tryInline bool) {
	markup := keyboard.InlineKeyboardMarkup()

	messages := a.keyboardMessage(update, markup, text, tryInline)

	for _, message := range messages {
		_, err := a.api.Send(message)
		if err != nil {
			log.Printf("update inline keyboard error: %s", err)
		}
	}

	a.lastMessageStorage.StoreByUpdate(&update, "")
}

func (a *EasyApi) DownloadFileString(fileId string) (string, error) {
	fileDirectLink, err := a.api.GetFileDirectURL(fileId)
	if err != nil {
		return "", fmt.Errorf("Getting file direct link error: %s", err)
	}

	cli := gentleman.New()
	cli.URL(fileDirectLink)

	req := cli.Request()

	res, err := req.Send()
	if err != nil {
		return "", fmt.Errorf("Downloading sended file error: %s", err)
	}
	if !res.Ok {
		return "", fmt.Errorf("Downloading sended file error invalid server response: %d\n", res.StatusCode)
	}

	return res.String(), nil
}

func (a EasyApi) keyboardMessage(
	update tgbotapi.Update,
	markup tgbotapi.InlineKeyboardMarkup,
	text string,
	tryInline bool,
) []tgbotapi.Chattable {
	if cq := update.CallbackQuery; cq != nil {
		if tryInline {
			return a.keyBoardMessageByCallbackQuery(cq, markup, text)
		} else {
			return a.keyBoardNewMessage(cq.Message.Chat, text, markup)
		}
	} else if message := update.Message; message != nil {
		return a.keyBoardNewMessage(message.Chat, text, markup)
	}

	return nil
}
func (a EasyApi) keyBoardMessageByCallbackQuery(
	callbackQuery *tgbotapi.CallbackQuery,
	markup tgbotapi.InlineKeyboardMarkup,
	text string,
) (messages []tgbotapi.Chattable) {
	messageText := tgbotapi.NewEditMessageText(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		text,
	)
	messageText.ReplyMarkup = &markup

	messages = append(messages, messageText)

	return messages
}

func (a EasyApi) keyBoardNewMessage(
	chat *tgbotapi.Chat,
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) []tgbotapi.Chattable {
	newMessage := tgbotapi.NewMessage(chat.ID, text)
	newMessage.ReplyMarkup = markup
	return []tgbotapi.Chattable{newMessage}
}
