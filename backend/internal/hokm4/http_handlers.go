package hokm4

import (
	"net/http"

	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/pkg/response"
)

func (app *ApplicationHokm4) isInActiveGame(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	var output struct {
		IsActive     bool  `json:"is_active"`
		ActiveGameId int64 `json:"active_game_id,omitempty"`
	}
	game, exist := app.Lobby.UserGames[user.ID]
	if exist {
		output.IsActive = true
		output.ActiveGameId = game.ID
	}

	if err := response.JSON(w, http.StatusOK, output); err != nil {
		app.ServerError(w, r, err)
		return
	}
}
