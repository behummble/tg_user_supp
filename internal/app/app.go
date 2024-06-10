package app

import (
	"log/slog"

	appbot "github.com/behummble/csupp_bot/internal/app/support_bot"
	"github.com/behummble/csupp_bot/internal/repo/db/redis"
	"github.com/behummble/csupp_bot/internal/service/support_bot"
	"github.com/behummble/csupp_bot/internal/config"
)

type App struct {
	Bot *appbot.Support
}

func New(log *slog.Logger, config *config.Config) App {
	db, err := redis.New(
		log, 
		config.Redis.Host, 
		config.Redis.Port, 
		config.Redis.Password)
	if err != nil {
		panic(err)
	}
	botService := supportbot.New(
		log, 
		db,
		config.Server.Host,
		config.Server.Path,
		config.Server.Port)

	appBot, err := appbot.New(
		log, 
		config.Bot.Token, 
		config.Bot.UpdateTimeout, 
		botService)

	if err != nil {
		panic(err)
	}

	return App{appBot}
}