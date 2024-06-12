package supportline

import (
	"context"
	"fmt"

	"golang.org/x/net/websocket"
)

type SupportLine struct {
	conn connect
}

type connect struct {
	host string
	port int
	path string
}

func New(host, path string, port int)  SupportLine {
	sLine := SupportLine{
		conn: connect{host, port, path},
	}
	err := sLine.ping()
	if err != nil {
		panic(err)
	}

	return sLine
}

func (sLine SupportLine) Send(ctx context.Context, payload string) error {
	ws, err := sLine.connect()
	if err != nil {
		return err
	}
	defer ws.Close()
	err = websocket.Message.Send(ws, payload)
	return err
}

func (sLine SupportLine) connect() (*websocket.Conn, error) {
	origin := fmt.Sprintf("http://%s:%d/", sLine.conn.host, sLine.conn.port)
	url := fmt.Sprintf("ws://%s:%d/%s", sLine.conn.host, sLine.conn.port, sLine.conn.path)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (sLine SupportLine) ping() error {
	ws, err := sLine.connect()
	if err != nil {
		return err
	}
	defer ws.Close()

	err = websocket.Message.Send(ws, "ping")
	return err
}