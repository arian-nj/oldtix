package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *ApplicationH2) wsHokm2Router() *chi.Mux {
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

	return mux
}
