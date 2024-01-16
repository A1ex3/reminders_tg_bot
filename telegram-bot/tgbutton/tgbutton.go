package tgbutton

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgButton struct {
	button      *tgbotapi.InlineKeyboardMarkup
	buttonsList [][]tgbotapi.InlineKeyboardButton
}

type ITgButton interface {
	Create(text, callBackData string) *tgbotapi.InlineKeyboardButton
	Add(button *tgbotapi.InlineKeyboardButton)
	Build() tgbotapi.InlineKeyboardMarkup
}

func Init() TgButton {
	button := &tgbotapi.InlineKeyboardMarkup{}
	buttonsList := make([][]tgbotapi.InlineKeyboardButton, 0)

	return TgButton{button: button, buttonsList: buttonsList}
}

func (*TgButton) Create(text, callBackData string) *tgbotapi.InlineKeyboardButton {
	_callBackData := new(string)
	*_callBackData = callBackData
	return &tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: _callBackData,
	}
}

func (tgb *TgButton) Add(button *tgbotapi.InlineKeyboardButton) {
	arrayButton := make([]tgbotapi.InlineKeyboardButton, 0)
	arrayButton = append(arrayButton, *button)
	tgb.buttonsList = append(tgb.buttonsList, arrayButton)
}

func (tgb *TgButton) Build() tgbotapi.InlineKeyboardMarkup {
	tgb.button.InlineKeyboard = tgb.buttonsList
	return tgbotapi.NewInlineKeyboardMarkup(tgb.button.InlineKeyboard...)
}
