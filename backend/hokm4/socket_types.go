package hokm4

import "github.com/arian-nj/master-card/back/internal/socket"

const (
	// TypeChat               socket.EventType = "chat"
	// TypeGameData           socket.EventType = "game_data"
	TypeNewCard            socket.EventType = "new_card"             // client <- server
	TypeNewCardOne         socket.EventType = "new_card_one"         // client <- server
	TypePlayerSelectedHokm socket.EventType = "player_selected_hokm" // server <- client
	TypeNewHokm            socket.EventType = "new_hokm"             // client <- server

	TypeMakeMatch   socket.EventType = "make_match"   // server <- client
	TypeDisconnect  socket.EventType = "disconnect"   // server <- client
	TypeMatchFound  socket.EventType = "found_match"  // client <- server
	TypeRejoinMatch socket.EventType = "rejoin_match" // client <- server

	TypeNewTrick socket.EventType = "new_trick" // client <- server
	TypeEndTrick socket.EventType = "end_trick" // client <- server

	// Game Turn Stuff
	TypeTurnStart   socket.EventType = "turn_start"   // client <- server
	TypeYourTurn    socket.EventType = "your_turn"    // client <- server
	TypeValidPlay   socket.EventType = "valid_play"   // client <- server
	TypeInvalidPlay socket.EventType = "invalid_play" // client <- server
	TypePlayTimeout socket.EventType = "play_timeout" // client <- server
	TypeTurnPlayed  socket.EventType = "turn_played"  // server <- client
	TypeTurnEnd     socket.EventType = "turn_end"     // client <- server

	TypeTheEnd socket.EventType = "the_end" // client <- server
)
