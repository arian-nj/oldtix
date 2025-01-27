package socket

import (
	"log"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 1 * time.Second
	pingInterval = pongWait * 9 / 10
)

type Client struct {
	conn           *websocket.Conn
	Egres          chan Event
	Handlers       *HandlerMap
	DefaultHandler WsEventHandler
}

type HandlerMap map[EventType]WsEventHandler
type WsEventHandler func(event *Event, client *Client) error

func (h HandlerMap) RegisterEventHandler(eventType EventType, handler WsEventHandler) {
	h[eventType] = handler
}
func NewHandlerMap() *HandlerMap {
	return &HandlerMap{}
}
func NewClient(conn *websocket.Conn, Handlers *HandlerMap) *Client {
	return &Client{
		conn:     conn,
		Handlers: Handlers,
		Egres:    make(chan Event),
	}
}

func (c *Client) ReadMessage(l *slog.Logger) {
	// defer func() {
	// 	c.manager.removeClient(c)
	// }()
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			l.Debug(err.Error())
			return
			// continue
		}

		event, err := EventUnmarshal(payload)
		if err != nil {
			log.Println("err in unmarshalling event: ", err)
			continue
		}
		if err := c.routeEvent(event, c); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (c *Client) WriteMessage() {
	// defer func() {
	// 	c.manager.removeClient(c)
	// }()

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Println("pong")

		return err
	})

	for {
		select {
		case event := <-c.Egres:
			payload, err := event.GetJsonByte()
			if err != nil {
				log.Println("err in marshalling event: ", err)
			}
			err = c.conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				log.Println(err)
			}

		case <-pingTicker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				log.Println("error in writing ping msg: ", err)
			}
			log.Println("ping")
		}
	}
}
