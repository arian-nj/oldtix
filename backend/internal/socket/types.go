package socket

type EventType string

const (
	TypeChat      EventType = "chat"
	TypeStatus    EventType = "status"
	TypeMakeMatch EventType = "make_match"
)

type EventMessage string

// status messages
const (
	StatusMatchFound EventMessage = "match found"
	StatusConnected  EventMessage = "connected"
)
