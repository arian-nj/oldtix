package socket

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = pongWait * 9 / 10
)

type ConncectionState int

const (
	OPEN ConncectionState = iota
	OPENING
	CLOSED
)

type Client struct {
	Conn      *websocket.Conn
	State     ConncectionState
	Egres     chan Event
	NewEvents chan Event

	CancelCtx context.Context
	Cancel    context.CancelFunc
}

func NewClient(conn *websocket.Conn) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		Conn:      conn,
		State:     OPENING,
		Egres:     make(chan Event),
		NewEvents: make(chan Event),
		CancelCtx: ctx,
		Cancel:    cancel,
	}

	err := conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return nil, err
	}

	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return err
	})

	conn.SetCloseHandler(func(code int, text string) error {
		client.State = CLOSED
		return nil
	})

	return client, nil
}

func (client *Client) Close() error {
	if client.Conn != nil {
		return client.Conn.Close()
	}
	client.State = CLOSED
	client.Cancel()
	return nil
}

func (client *Client) ReadMessage(l *zap.Logger, ctx context.Context) error {
	defer func() {
		client.State = CLOSED
		client.Cancel()
		client.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, payload, err := client.Conn.ReadMessage()
			if err != nil {
				return err
			}

			event, err := EventUnmarshal(payload)
			if err != nil {
				l.Debug("err in unmarshalling event: " + err.Error())
				continue
			}
			client.NewEvents <- *event
		}

	}

}

func (client *Client) WriteMessage(l *zap.Logger, ctx context.Context) error {
	defer func() {
		client.State = CLOSED
		client.Cancel()
		client.Close()
	}()

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	client.State = OPEN

	for {
		select {
		case event := <-client.Egres:
			payload, err := event.GetJsonByte()
			if err != nil {
				l.Debug("err in marshalling event: " + err.Error())
			}
			err = client.Conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				l.Debug(err.Error())
			}

		case <-pingTicker.C:
			err := client.Conn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				l.Debug("error in writing ping msg: " + err.Error())
			}
			// l.Debug("ping")
		case <-ctx.Done():
			return nil
		}

	}
}
