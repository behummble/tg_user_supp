package supportbot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ErrorMessage = "Sorry, something went wrong, try to repeat the request later"
)

type MessageSaver interface {
	Save(ctx context.Context, botName, msg string) error
}

type BotService struct {
	log *slog.Logger
	Saver MessageSaver
}

func New(log *slog.Logger, saver MessageSaver) *BotService {
	return &BotService{
		log,
		saver,
	}
}

func(sbot *BotService) ProcessUpdate(upd tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	resp := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	return resp, nil
}
