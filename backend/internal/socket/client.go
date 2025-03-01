package socket

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = pongWait * 9 / 10
)

type Client struct {
	Conn      *websocket.Conn
	Egres     chan Event
	NewEvents chan Event
}

func (c *Client) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn:      conn,
		Egres:     make(chan Event),
		NewEvents: make(chan Event),
	}
}

func (c *Client) ReadMessage(l *slog.Logger, ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, payload, err := c.Conn.ReadMessage()
			if err != nil {
				l.Debug(err.Error())
				return nil
			}

			event, err := EventUnmarshal(payload)
			if err != nil {
				l.Debug("err in unmarshalling event: " + err.Error())
				continue
			}
			c.NewEvents <- *event
		}

	}
}

func (c *Client) WriteMessage(l *slog.Logger, ctx context.Context) error {

	// defer func() {
	// 	c.manager.removeClient(c)
	// }()

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			l.Debug(err.Error())
		}
		// l.Debug("pong")
		return err
	})

	for {
		select {
		case event := <-c.Egres:
			payload, err := event.GetJsonByte()
			if err != nil {
				l.Debug("err in marshalling event: " + err.Error())
			}
			err = c.Conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				l.Debug(err.Error())
			}

		case <-pingTicker.C:
			err := c.Conn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				l.Debug("error in writing ping msg: " + err.Error())
			}
			// l.Debug("ping")
		case <-ctx.Done():
			return nil
		}

	}
}
