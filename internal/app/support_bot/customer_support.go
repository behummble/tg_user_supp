package appsupportbot

import (
	"log/slog"
	"time"
	"gopkg.in/telebot.v3"
	csupp "github.com/behummble/csupp_bot/internal/service/support_bot"
)

type Support struct {
	log *slog.Logger
	bot *telebot.Bot
	botService *csupp.BotService
}

func New(log *slog.Logger, token string, timeout int, botService *csupp.BotService) (*Support, error) {
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token: token,
			Poller: &telebot.LongPoller{Timeout: time.Second * time.Duration(timeout)},
		},
	)
	if err != nil {
		return nil, err
	}

	return &Support{log, bot, botService}, nil
} 

func(support *Support) ListenUpdates(timeout int, botName string) {
	support.bot.Handle(telebot.OnText, support.TextHandler)
	support.bot.Start()
}

func(support *Support) TextHandler(upd telebot.Context) error {
	name := support.bot.Me.Username
	id := support.bot.Me.ID
	resp, err := support.botService.ProcessUpdate(upd, name, id)
	if err != nil {
		return err
	}

	return upd.Send(resp)
}