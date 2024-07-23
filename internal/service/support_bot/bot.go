package supportbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/behummble/csupp_bot/pkg/crypto"
	"gopkg.in/telebot.v3"
)

const (
	errorMessage = "Sorry, something went wrong, try to repeat the request later"
	startMessage = "Hi, %s!"
)

const (
	commandNotFound = "Can`t find command %s"
)

var (
	messagesQueue = "messages:{%s}"
	topicQueue = "topic:{%d}"
)

type DB interface {
	Save(ctx context.Context, botName, msg string) error
	Topic(ctx context.Context, topicKey string) (string, error)
}

type UserSupport interface {
	SendToUser(ctx context.Context, payload string) error
	SendToSupport(ctx context.Context, payload string) error
}

type BotService struct {
	log *slog.Logger
	db DB
	userSupport UserSupport
	groupChatID int64
}

type Message struct {
	BotToken string
	ChatID int64
	UserID int64
	UserName string
	Payload string
	MessageID int64
	GroupChatID int64
}

type TopicData struct {
	BotToken string
	ChatID int64
	UserID int64
	TopicID int
}

type SupportMessage struct {
	BotToken string
	ChatID int64
	TopicID int
	Payload string
}

func New(log *slog.Logger, saver DB, userSupport UserSupport, chatID int64) *BotService {
	
	return &BotService{
		log,
		saver,
		userSupport,
		chatID,
	}
}

func(sbot *BotService) ProcessUpdate(upd telebot.Context, botName string, botID int64) (string, error) {
	resp, err := sbot.handleEvent(upd, botID, botName)
	if err != nil {
		sbot.log.Error("ProcessError", err)
		return errorResponse(), nil
	}

	return resp, nil
}

func(sbot *BotService) handleEvent(upd telebot.Context, botID int64, botName string) (string, error) {
	if isCommand(upd.Entities()) {
		return handleCommand(upd)
	}
	
	topicID := executeTopicID(upd)
	if topicID != 0 { 
		err := handleSupportMessage(
			sbot.userSupport, 
			upd.Bot().Token,
			upd.Chat().ID,
			topicID, 
			upd.Text())
		if err != nil {
			sbot.log.Error("HandleSupportMessage", err)
		}
		return "", nil
	}

	msg, err := prepareMessage(upd, botID, sbot.groupChatID)

	if err == nil {
		err = sbot.userSupport.SendToSupport(context.Background(), msg)	
	}
	return "", err
}

func executeTopicID(upd telebot.Context) int {
	topic := upd.Topic()
	if topic != nil {
		return topic.ThreadID
	}
	msg := upd.Message()
	if msg != nil {
		return msg.ThreadID
	}
	return 0
}

func handleCommand(upd telebot.Context) (string, error) {
	switch upd.Text() {
	case "/start":
		return handleStart(upd), nil
	case "/action1":
		return handleAction1(), nil
	case "/action2":
		return handleAction2(), nil
	default:
		return "", fmt.Errorf(commandNotFound, upd.Text())
	}
}

func handleSupportMessage(userSupport UserSupport ,token string, chatID int64, topicID int, payload string) error {
	
	msg, err := prepareSupportMessage(token, chatID, topicID, payload)
	if err != nil {
		return err
	} 
	
	if err = userSupport.SendToUser(context.Background(), msg); err != nil {
		return err
	}
	return nil
}

func prepareMessage(upd telebot.Context, botID int64, groupChatID int64) (string, error) {
	token, err := crypto.EncryptData([]byte(upd.Bot().Token))

	if err != nil {
		return "", err
	}
	
	msg := Message {
		token,
		upd.Chat().ID,
		upd.Message().Sender.ID,
		fmt.Sprintf("%s %s",upd.Message().Sender.FirstName, upd.Message().Sender.LastName),
		upd.Text(),
		int64(upd.Message().ID),
		groupChatID,
	}

	payload, err := json.Marshal(msg)
	
	return string(payload), err
}

func handleStart(upd telebot.Context) string {
	if upd.Sender() != nil {
		return fmt.Sprintf(startMessage, upd.Sender().Username)
	}
	
	return fmt.Sprintf(startMessage, "User")
}

func handleAction1() string {
	 return "Thanks for your help"
}

func handleAction2() string {
	 return "Your phone belongs to me now"
}

func errorResponse() string {
	return errorMessage
}

func isCommand(entities []telebot.MessageEntity) bool {
	if len(entities) == 0 {
		return false
	}

	return entities[0].Type == telebot.EntityCommand
}

func parseTopic(data string) (TopicData, error) {
	var topic TopicData
	err := json.Unmarshal([]byte(data), &topic)
	return topic, err
}

func prepareSupportMessage(token string, chatID int64, topicID int, payload string) (string, error) {
	encToken, err := crypto.EncryptData([]byte(token))
	if err != nil {
		return "", err
	}
	data := SupportMessage{
		encToken,
		chatID,
		topicID,
		payload,
	}

	res, err := json.Marshal(data)
	return string(res), err
}