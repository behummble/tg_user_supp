package supportbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"gopkg.in/telebot.v3"
	"golang.org/x/net/websocket"
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

type BotService struct {
	log *slog.Logger
	db DB
	ws *websocket.Conn
}

type Message struct {
	BotID int64
	ChatID int64
	UserID int64
	UserName string
	Payload string
	MessageID int64
}

type TopicData struct {
	BotID int64
	ChatID int64
	UserID int64
	TopicID int
}

type SupportMessage struct {
	ChatID int64
	TopicID int
	Text string
}

func New(log *slog.Logger, saver DB, host, path string, port int) *BotService {
	origin := fmt.Sprintf("http://%s:%d/", host, port)
	url := fmt.Sprintf("ws://%s:%d/%s", host, port, path)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	return &BotService{
		log,
		saver,
		ws,
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
		data, err := sbot.db.Topic(
			context.Background(), 
			fmt.Sprintf(topicQueue, topicID))
		if err != nil {
			sbot.log.Error("GetTopicData", err)
			return "", nil
		}
		err = handleSupportMessage(sbot.ws, data, upd.Text())
		if err != nil {
			sbot.log.Error("HandleSupportMessage", err)
		}
		return "", nil
	}

	if upd.Text() != "" {
		msg, err := prepareMessage(upd, botID)
		if err == nil {
			err = sbot.db.Save(
				context.Background(),
				fmt.Sprintf(messagesQueue, botName),
				msg,
			)
		}
		
		return "", err
	}

	return "", nil
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

func handleSupportMessage(ws *websocket.Conn, topicStr, payload string) error {
	topicData, err := parseTopic(topicStr)
	if err != nil {
		return err
	}
	msg, err := prepareSupportMessage(topicData, payload)
	if err != nil {
		return err
	} 
	
	if err = websocket.Message.Send(ws, msg); err != nil {
		return err
	}
	return nil
}

func prepareMessage(upd telebot.Context, botID int64) (string, error) {
	msg := Message {
		botID,
		upd.Chat().ID,
		upd.Message().Sender.ID,
		fmt.Sprintf("%s %s",upd.Message().Sender.FirstName, upd.Message().Sender.LastName),
		upd.Text(),
		int64(upd.Message().ID),
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

func prepareSupportMessage(topic TopicData, payload string) (string, error) {
	data := SupportMessage{
		topic.ChatID,
		topic.TopicID,
		payload,
	}

	res, err := json.Marshal(data)
	return string(res), err
}