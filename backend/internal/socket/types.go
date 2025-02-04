package socket

type EventType string

const (
	TypeChat        EventType = "chat"
	TypeStatus      EventType = "status"
	TypeMakeMatch   EventType = "make_match"
	TypeGameData    EventType = "game_data"
	TypeGetData     EventType = "get_data"
	TypeNewCard     EventType = "new_card"
	TypeHokmChoosed EventType = "hokm_choosed"
)

type EventMessage string

// status messages
var (
	StatusConnected EventMessage = "connected"
)
