package easy

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Key struct {
	text string
	data string
}

type KeySplit struct {}

type EasyKeyboard struct {
	list []interface{}
	width int
}

func (k *EasyKeyboard) SetWidth(width int) *EasyKeyboard {
	if width > 0 && width <= 10 {
		k.width = width
	} else {
		k.width = 0
	}
	return k
}

func (k *EasyKeyboard) Add(text string, data string) {
	k.list = append(k.list, Key{text: text, data: data})
}

func (k *EasyKeyboard) AddSplit() {
	k.list = append(k.list, KeySplit{})
}

func (k EasyKeyboard) InlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var buttons []tgbotapi.InlineKeyboardButton
	for _, key := range k.list {
		switch key.(type) {
		case KeySplit:
			if len(buttons) > 0 {
				rows = append(rows, buttons)
				buttons = nil
			}
		case Key:
			button := tgbotapi.NewInlineKeyboardButtonData(key.(Key).text, key.(Key).data)
			buttons = append(buttons, button)
			if len(buttons) >= k.realWidth() {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons...))
				buttons = nil
			}
		}
	}

	if len(buttons) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons...))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (k EasyKeyboard) realWidth() int {
	if k.width == 0 {
		return 3
	}
	return k.width
}