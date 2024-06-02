package appsupportbot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	csupp "github.com/behummble/csupp_bot/internal/service/support_bot"
)

type Support struct {
	log *slog.Logger
	bot *tgbotapi.BotAPI
	botService *csupp.BotService
}

func New(log *slog.Logger, token string, botService *csupp.BotService) (*Support, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &Support{log, bot, botService}, nil
} 

func(support *Support) StartListenUpdates(timeout int) {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = timeout
	updates := support.bot.GetUpdatesChan(update)
	for {
		upd := <-updates
		resp, err := support.botService.ProcessUpdate(upd)
		if err != nil {
			support.log.Error("ProcessUpdate", err)
		}
		support.bot.Send(resp)
	}
}