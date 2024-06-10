package main

import (
	"log/slog"
	"os"

	"github.com/behummble/csupp_bot/internal/config"
	"github.com/behummble/csupp_bot/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	log := initLog()
	setEnv()
	config := config.MustLoad()
	app := app.New(
		log, 
		config)
	app.Bot.ListenUpdates(config.Bot.UpdateTimeout, config.Bot.Name)	
}

func initLog() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(
		os.Stdout, 
		&slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}

func setEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
}