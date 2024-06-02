package supportbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ErrorMessage = "Sorry, something went wrong, try to repeat the request later"
)

func ProcessUpdate(upd tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	resp := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	return resp, nil
}
