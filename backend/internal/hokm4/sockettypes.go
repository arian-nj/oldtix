package hokm4

import "github.com/arian-nj/master-card/back/internal/socket"

const (
	TypeChat        socket.EventType = "chat"
	TypeGameData    socket.EventType = "game_data"
	TypeGetData     socket.EventType = "get_data"
	TypeNewCard     socket.EventType = "new_card"
	TypeHokmChoosed socket.EventType = "hokm_choosed"

	TypeGetMyCards socket.EventType = "get_my_cards"

	TypeMakeMatch   socket.EventType = "make_match"
	TypeDisconnect  socket.EventType = "disconnect"
	TypeMatchFound  socket.EventType = "found_match"
	TypeRejoinMatch socket.EventType = "rejoin_match"

	TypeNewTrick socket.EventType = "new_trick"
	TypeEndTrick socket.EventType = "end_trick"

	// Game Turn Stuff
	TypeTurnStart   socket.EventType = "turn_start"
	TypeYourTurn    socket.EventType = "your_turn"
	TypeValidPlay   socket.EventType = "valid_play"
	TypeInvalidPlay socket.EventType = "invalid_play"
	TypePlayTimeout socket.EventType = "play_timeout"
	TypeTurnPlayed  socket.EventType = "turn_played"
	TypeTurnEnd     socket.EventType = "turn_end"

	TypeTheEnd socket.EventType = "the_end"
)
