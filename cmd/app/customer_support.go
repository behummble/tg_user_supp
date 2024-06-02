package main

import (
	"log/slog"
	"os"

	"github.com/behummble/csupp_bot/internal/config"
	"github.com/behummble/csupp_bot/internal/app"
)

func main() {
	log := initLog()
	config := config.MustLoad()
	app := app.New(
		log, 
		config.Bot.Token, 
		config.Redis.Host, 
		config.Redis.Port,
		config.Redis.Password,
	)
	go app.Bot.StartListenUpdates(config.Bot.UpdateTimeout)
	
}

func initLog() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(
		os.Stdout, 
		&slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}