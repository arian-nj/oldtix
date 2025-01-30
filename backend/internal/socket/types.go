package socket

type EventType string

const (
	TypeChat       EventType = "chat"
	TypeStatus     EventType = "status"
	TypeMakeMatch  EventType = "make_match"
	TypeMatchFound EventType = "match_found"
)

type EventMessage string

// status messages
const (
	StatusConnected EventMessage = "connected"
)
