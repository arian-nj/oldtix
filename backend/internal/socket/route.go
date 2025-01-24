package socket

import (
	"errors"
)

func (c *Client) routeEvent(event *Event, client *Client) error {
	handler := *c.Handlers
	if _, ok := handler[event.Type]; ok {
		err := handler[event.Type](event, client)
		return err
	}
	return errors.New("handler for " + string(event.Type) + " type does not exist")
}
