package supportline

import (
	"context"
	"fmt"

	"golang.org/x/net/websocket"
)

const (
	userMessagePath = "user/message"
	supportMessagePath = "support/message"
)

type SupportLine struct {
	conn connect
}

type connect struct {
	host string
	port int
}

func New(host string, port int)  SupportLine {
	sLine := SupportLine{
		conn: connect{host, port},
	}
	err := sLine.ping()
	if err != nil {
		panic(err)
	}

	return sLine
}

func (sLine SupportLine) SendToUser(ctx context.Context, payload string) error {
	ws, err := sLine.connect(supportMessagePath)
	if err != nil {
		return err
	}
	defer ws.Close()
	err = websocket.Message.Send(ws, payload)
	return err
}

func (sLine SupportLine) SendToSupport(ctx context.Context, payload string) error {
	ws, err := sLine.connect(userMessagePath)
	if err != nil {
		return err
	}
	defer ws.Close()
	err = websocket.Message.Send(ws, payload)
	return err
}

func (sLine SupportLine) connect(path string) (*websocket.Conn, error) {
	origin := fmt.Sprintf("http://%s:%d/", sLine.conn.host, sLine.conn.port)
	url := fmt.Sprintf("ws://%s:%d/%s", sLine.conn.host, sLine.conn.port, path)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (sLine SupportLine) ping() error {
	ws, err := sLine.connect("ping")
	if err != nil {
		return err
	}
	defer ws.Close()

	err = websocket.Message.Send(ws, "ping")
	return err
}