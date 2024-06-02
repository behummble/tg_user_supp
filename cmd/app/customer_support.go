package main

import (
	"log/slog"
	"os"

	"github.com/behummble/csupp_bot/internal/config"
	"github.com/behummble/csupp_bot/internal/app/support_bot"
)

func main() {
	log := initLog()
	config := config.MustLoad()
	app := appsupportbot.New(log, config.Bot.Token)
	app.StartListenUpdates()
}

func initLog() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(
		os.Stdout, 
		&slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}