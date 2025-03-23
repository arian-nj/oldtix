package hokm4

import (
	"github.com/go-chi/chi/v5"
)

func (app *ApplicationHokm4) Hokm4Router() *chi.Mux {
	// http router
	mux := chi.NewRouter()

	mux.NotFound(app.NotFound)
	mux.MethodNotAllowed(app.MethodNotAllowed)

	mux.Use(app.RecoverPanic)

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
