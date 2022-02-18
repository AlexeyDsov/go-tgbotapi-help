package easy

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type EasyApiKeyboard struct {
	api *tgbotapi.BotAPI
	lastMessageStorage lastMessageStorage
}

func (a *EasyApiKeyboard) SendKeyboard(update tgbotapi.Update, keyboard EasyKeyboard, text string, tryInline bool) {
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

func (a *EasyApiKeyboard) keyboardMessage(
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
func (a *EasyApiKeyboard) keyBoardMessageByCallbackQuery(
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

func (a *EasyApiKeyboard) keyBoardNewMessage(
	chat *tgbotapi.Chat,
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) []tgbotapi.Chattable {
	newMessage := tgbotapi.NewMessage(chat.ID, text)
	newMessage.ReplyMarkup = markup
	return []tgbotapi.Chattable{newMessage}
}
