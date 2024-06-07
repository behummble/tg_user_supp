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

var (
	messagesQueue = "messages:{%s}"
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
	MessageID int64
}

func New(log *slog.Logger, saver MessageSaver) *BotService {
	return &BotService{
		log,
		saver,
	}
}

func(sbot *BotService) ProcessUpdate(upd tgbotapi.Update, botName string, botID int64) (tgbotapi.MessageConfig, error) {
	resp, err := sbot.handleEvent(upd, botID)
	if err != nil {
		return errorResponse(upd.Message.Chat.ID), err
	}

	return resp, nil
}

func(sbot *BotService) handleEvent(upd tgbotapi.Update, botID int64) (tgbotapi.MessageConfig, error) {
	if upd.Message.IsCommand() {
		return handleCommand(upd)
	}
	
	if upd.Message.Text != "" {
		msg, err := prepareMessage(upd, botID)
		if err == nil {
			err = sbot.Saver.Save(
				context.Background(),
				fmt.Sprintf(messagesQueue, botID),
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

func prepareMessage(upd tgbotapi.Update, botID int64) (string, error) {
	msg := Message {
		botID,
		upd.Message.Chat.ID,
		upd.Message.From.ID,
		fmt.Sprintf("%s %s",upd.Message.From.FirstName, upd.Message.From.LastName),
		upd.Message.Text,
		int64(upd.Message.MessageID),
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
