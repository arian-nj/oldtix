package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *CommonGlobals) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *CommonGlobals) ValidateToken(w http.ResponseWriter, r *http.Request, token string) *http.Request {
	claims, err := jwt.HMACCheck([]byte(token), []byte(app.Config.Jwt.SecretKey))
	if err != nil {
		app.invalidAuthenticationToken(w, r)
		return nil
	}

	if !claims.Valid(time.Now()) {
		app.invalidAuthenticationToken(w, r)
		return nil
	}

	if claims.Issuer != app.Config.BaseURL {
		app.invalidAuthenticationToken(w, r)
		return nil
	}

	if !claims.AcceptAudience(app.Config.BaseURL) {
		app.invalidAuthenticationToken(w, r)
		return nil
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	user, err := app.Queries.GetUser(context.Background(), int32(userID))
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	if user.ID != 0 {
		return ContextSetAuthenticatedUser(r, &user)
	}
	return nil
}

func (app *CommonGlobals) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")

			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token := headerParts[1]
				new_request := app.ValidateToken(w, r, token)
				if new_request == nil {
					return
				}
				r = new_request
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *CommonGlobals) AuthenticateQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("auth_token")
		if token != "" {
			new_request := app.ValidateToken(w, r, token)
			if new_request == nil {
				return
			}
			r = new_request
		}
		next.ServeHTTP(w, r)
	})
}

func (app *CommonGlobals) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := ContextGetAuthenticatedUser(r)

		if authenticatedUser == nil {
			app.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
