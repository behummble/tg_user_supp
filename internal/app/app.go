package app

import (
	"log/slog"

	appbot "github.com/behummble/csupp_bot/internal/app/support_bot"
	"github.com/behummble/csupp_bot/internal/repo/db/redis"
	"github.com/behummble/csupp_bot/internal/service/support_bot"
)

type App struct {
	Bot *appbot.Support
}

func New(log *slog.Logger, token, dbHost, dbPort, dbPassword string) App {
	db, err := redis.New(log, dbHost, dbPort, dbPassword)
	if err != nil {
		panic(err)
	}
	botService := supportbot.New(log, db)
	appBot, err := appbot.New(log, token, botService)
	if err != nil {
		panic(err)
	}
	return App{appBot}
}