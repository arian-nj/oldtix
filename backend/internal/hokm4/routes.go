package hokm4

import (
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/go-chi/chi/v5"
)

type ApplicationHokm4 struct {
	*server.CommonGlobals
	Lobby *Lobby
}

func NewHokm4Application(globalStructs *server.CommonGlobals) *ApplicationHokm4 {
	return &ApplicationHokm4{
		CommonGlobals: globalStructs,
		Lobby:         NewLobby(),
	}
}

func (app *ApplicationHokm4) Hokm4Router() *chi.Mux {
	// http router
	mux := chi.NewRouter()

	mux.NotFound(app.NotFound)
	mux.MethodNotAllowed(app.MethodNotAllowed)

	mux.Use(app.RecoverPanic)
	mux.Use(app.CorsMiddlewareFunc)
	// Authenticated websocket routes
	mux.Group(func(r chi.Router) {
		r.Use(app.AuthenticateQuery)
		r.Use(app.RequireAuthenticatedUser)

		r.Get("/ws", app.WsUpgradeHandler)
	})

	mux.Group(func(r chi.Router) {
		r.Use(app.Authenticate)
		r.Use(app.RequireAuthenticatedUser)

		r.Get("/active_game", app.isInActiveGame)
	})

	return mux
}
