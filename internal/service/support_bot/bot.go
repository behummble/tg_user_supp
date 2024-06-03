package supportbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	errorMessage = "Sorry, something went wrong, try to repeat the request later"
	startMessage = "Hi, %s"
)

const (
	commandNotFound = "Can`t find command %s"
)

type MessageSaver interface {
	Save(ctx context.Context, botName, msg string) error
}

type BotService struct {
	log *slog.Logger
	Saver MessageSaver
}

type Message struct {
	BotID int64
	ChatID int64
	UserID int64
	UserName string
	Payload string
}

func New(log *slog.Logger, saver MessageSaver) *BotService {
	return &BotService{
		log,
		saver,
	}
}

func(sbot *BotService) ProcessUpdate(upd tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	resp, err := sbot.handleEvent(upd)
	if err != nil {
		return errorResponse(upd.Message.Chat.ID), err
	}

	return resp, nil
}

func(sbot *BotService) handleEvent(upd tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	// switch event and return 
	if upd.Message.IsCommand() {
		return handleCommand(upd)
	}
	
	if upd.Message.Text != "" {
		msg, err := prepareMessage(upd)
		if err == nil {
			err = sbot.Saver.Save(
				context.Background(),
				upd.Message.ViaBot.UserName,
				msg,
			)
		}
		
		return tgbotapi.MessageConfig{}, err
	}

	return tgbotapi.MessageConfig{}, nil
}

func handleCommand(upd tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	switch upd.Message.Command() {
	case "start":
		return handleStart(upd), nil
	case "action1":
		return handleAction1(upd), nil
	case "action2":
		return handleAction2(upd), nil
	default:
		return tgbotapi.MessageConfig{}, fmt.Errorf(commandNotFound, upd.Message.Command())
	}
}

func prepareMessage(upd tgbotapi.Update) (string, error) {
	msg := Message {
		upd.Message.ViaBot.ID,
		upd.Message.Chat.ID,
		upd.Message.From.ID,
		upd.Message.From.UserName,
		upd.Message.Text,
	}

	payload, err := json.Marshal(msg)
	return string(payload), err
}

func handleStart(upd tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(
		upd.Message.Chat.ID, 
		fmt.Sprintf(startMessage, upd.Message.From.UserName),
	)
}

func handleAction1(upd tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(
		upd.Message.Chat.ID, 
		"Thanks for your help",
	)
}

func handleAction2(upd tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(
		upd.Message.Chat.ID, 
		"Your phone belongs to me now",
	)
}

func errorResponse(chatID int64) (tgbotapi.MessageConfig) {
	return tgbotapi.NewMessage(chatID, errorMessage)
}
