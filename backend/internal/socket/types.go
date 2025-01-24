package socket

type EventType string

const (
	EventTypeChat   EventType = "chat"
	EventTypeStatus EventType = "status"
)

type EventMessage string

// status messages

const (
	StatusMatchFound EventMessage = "match found"
	StatusConnected  EventMessage = "connected"
)
