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

	// All routes
	mux.Get("/status", app.status)
	mux.Get("/version", app.getLatestVersion)
	// mux.Post("/register", app.register)
	// mux.Post("/token", app.createAuthenticationToken)
	mux.Get("/auth/guest/create", app.createGuest)
	mux.Post("/auth/guest/token", app.createGuestToken)

	// Authenticated REST routes
	mux.Group(func(authRouter chi.Router) {

		authRouter.Use(app.Authenticate)
		authRouter.Use(app.RequireAuthenticatedUser)

		authRouter.Put("/update", app.updateUserData)
		authRouter.Get("/me", app.getMeData)

		authRouter.Get("/person/{user_id}", app.getUserData)
		authRouter.Get("/person/{user_id}/stat", app.getUserStatisticsData)
	})

	return mux
}
