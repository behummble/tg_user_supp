package redis

import (
	"log/slog"
	"os"
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	log *slog.Logger
	conn *redis.Client
}

func New(log *slog.Logger) (*Client, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}

	return &Client{log, conn}, nil
}

func connect() (*redis.Client, error) {
	options := &redis.Options{
		Addr: os.Getenv("REDIS_ADDRES"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	conn := redis.NewClient(options)

	_, err := conn.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}
	
	return conn, nil
}

func (client *Client) Push(botName, message string) {
	_, err := client.conn.LPush(context.Background(), botName, message).Result()
	if err != nil {
		client.log.Error("PushMessage", err)
	}
}

func (client *Client) Pop(botName string) []string {
	messages, err := client.conn.BLPop(context.Background(), time.Second * 3, botName).Result()
	if err != nil {
		client.log.Error("PopMessage", err)
		return []string{}
	}
	return messages
} 