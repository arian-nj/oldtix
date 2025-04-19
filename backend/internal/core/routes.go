package core_api

import (
	"os"

	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/go-chi/chi/v5"
)

type ApiApplication struct {
	*server.CommonGlobals
	ReleaseMode string
}

func NewApiApplication(globalStructs *server.CommonGlobals) *ApiApplication {
	return &ApiApplication{
		CommonGlobals: globalStructs,
		ReleaseMode:   os.Getenv("RELEASE_MODE"),
	}
}

func (app *ApiApplication) CoreRoutes() *chi.Mux {
	// http router
	mux := chi.NewRouter()

	mux.NotFound(app.NotFound)
	mux.MethodNotAllowed(app.MethodNotAllowed)

	mux.Use(app.RecoverPanic)

	mux.Group(func(r chi.Router) {
		r.Use(app.Authenticate)

		// All routes
		r.Get("/status", app.status)
		r.Get("/version", app.getLatestVersion)
		r.Post("/register", app.register)
		r.Post("/token", app.createAuthenticationToken)

		// Authenticated REST routes
		r.Group(func(r2 chi.Router) {
			r2.Use(app.RequireAuthenticatedUser)
			r2.Put("/update", app.updateUserData)
			r2.Get("/me", app.getMeData)
			r2.Get("/person/{user_id}", app.getUserData)
			r2.Get("/person/{user_id}/stat", app.getUserStatisticsData)
		})
	})

	return mux
}
