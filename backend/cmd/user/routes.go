package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) profileRoutes() *chi.Mux {

	// http router
	mux := chi.NewRouter()

	mux.NotFound(app.NotFound)
	mux.MethodNotAllowed(app.MethodNotAllowed)

	mux.Use(app.RecoverPanic)

	mux.Group(func(r chi.Router) {
		r.Use(app.Authenticate)

		// All routes
		r.Get("/status", app.status)
		r.Post("/user/register", app.register)
		r.Post("/user/token", app.createAuthenticationToken)

		// Authenticated REST routes
		r.Group(func(r2 chi.Router) {
			r2.Use(app.RequireAuthenticatedUser)
			r2.Put("/user/update", app.updateUserData)
			r2.Get("/user/me", app.getUserData)
		})
	})

	return mux
}
