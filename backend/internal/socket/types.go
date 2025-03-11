package socket

type EventType string

const (
	TypeChat        EventType = "chat"
	TypeGameData    EventType = "game_data"
	TypeGetData     EventType = "get_data"
	TypeNewCard     EventType = "new_card"
	TypeHokmChoosed EventType = "hokm_choosed"

	TypeMakeMatch  EventType = "make_match"
	TypeMatchFound EventType = "found_match"

	TypeNewTrick EventType = "new_trick"
	TypeEndTrick EventType = "end_trick"

	// Game Turn Stuff
	TypeTurnStart   EventType = "turn_start"
	TypeYourTurn    EventType = "your_turn"
	TypeValidPlay   EventType = "valid_play"
	TypeInvalidPlay EventType = "invalid_play"
	TypePlayTimeout EventType = "play_timeout"
	TypeTurnPlayed  EventType = "turn_played"
	TypeTurnEnd     EventType = "turn_end"
)

type EventMessage string
