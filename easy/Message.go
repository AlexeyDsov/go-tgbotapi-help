package easy

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Message struct {
	update    *tgbotapi.Update
	text      string
	command   string
	deepSetup  func(message tgbotapi.Message) error
	ParseMode  string
	oldReplace bool
}

func NewMessage(update *tgbotapi.Update, text string) *Message {
	return &Message{update: update, text: text}
}

func (m *Message) WithCommand(command string) *Message {
	m.command = command
	return m
}

func (m *Message) WithDeepSetup(f func(message tgbotapi.Message) error) *Message {
	m.deepSetup = f
	return m
}

func (m *Message) WithParseMk2() *Message {
	m.ParseMode = tgbotapi.ModeMarkdownV2
	return m
}

func (m *Message) WithOldReplace() *Message {
	m.oldReplace = true
	return m
}
