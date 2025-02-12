package socket

type EventType string

const (
	TypeChat        EventType = "chat"
	TypeMakeMatch   EventType = "make_match"
	TypeGameData    EventType = "game_data"
	TypeGetData     EventType = "get_data"
	TypeNewCard     EventType = "new_card"
	TypeHokmChoosed EventType = "hokm_choosed"

	// Game Turn Stuff
	TypeTurnStart   EventType = "turn_start"
	TypeYourTurn    EventType = "your_turn"
	TypePlayTurn    EventType = "play_turn"
	TypeValidPlay   EventType = "valid_play"
	TypeInvalidPlay EventType = "invalid_play"
	TypePlayTimeout EventType = "play_timeout"
	TypeTurnPlayed  EventType = "turn_played"
)

type EventMessage string
