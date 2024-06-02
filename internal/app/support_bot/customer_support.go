package appsupportbot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	csupp "github.com/behummble/csupp_bot/internal/service/support_bot"
)

type Support struct {
	log *slog.Logger
	bot *tgbotapi.BotAPI
}

func New(log *slog.Logger, token string) *Support {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	return &Support{log, bot}
} 

func(support *Support) StartListenUpdates() {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 10
	updates := support.bot.GetUpdatesChan(update)
	for {
		upd := <-updates
		resp, err := csupp.ProcessUpdate(upd)
		if err != nil {
			support.log.Error("ProcessUpdate", err)
		}
		support.bot.Send(resp)
	}
}