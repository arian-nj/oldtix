package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
func (app *CommonGlobals) CorsMiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		// In a production environment, you might want to restrict the origin:
		// w.Header().Set("Access-Control-Allow-Origin", "https://your-specific-domain.com")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Optional

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			// You can add more logic here if needed, e.g., checking specific headers
			// for the preflight request.
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// func (app *CommonGlobals) ValidateToken(w http.ResponseWriter, r *http.Request, token string) *http.Request {
// 	claims, err := jwt.HMACCheck([]byte(token), []byte(app.Config.Jwt.SecretKey))
// 	if err != nil {
// 		app.invalidAuthenticationToken(w, r)
// 		return nil
// 	}

// 	if !claims.Valid(time.Now()) {
// 		app.invalidAuthenticationToken(w, r)
// 		return nil
// 	}

// 	if claims.Issuer != app.Config.BaseURL {
// 		app.invalidAuthenticationToken(w, r)
// 		return nil
// 	}

// 	if !claims.AcceptAudience(app.Config.BaseURL) {
// 		app.invalidAuthenticationToken(w, r)
// 		return nil
// 	}

// 	userID, err := strconv.Atoi(claims.Subject)
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return nil
// 	}

// 	user, err := app.Queries.GetPerson(r.Context(), userID)
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return nil
// 	}

// 	if user.ID != 0 {
// 		return ContextSetAuthenticatedUser(r, &user)
// 	}
// 	return nil
// }

func (app *CommonGlobals) ValidateToken(w http.ResponseWriter, r *http.Request, tokenString string) *http.Request {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return app.Config.Jwt.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		app.invalidAuthenticationToken(w, r)
		return nil
	}

	expireAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}
	if expireAt.Time.Unix() < time.Now().Unix() {

		app.invalidAuthenticationToken(w, r)
		return nil
	}

	notBefore, err := token.Claims.GetNotBefore()
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	if notBefore.Time.Unix() > time.Now().Unix() {
		app.invalidAuthenticationToken(w, r)
		return nil
	}

	sub, err := token.Claims.GetSubject()

	if err != nil {
		app.Logger.Info("here")
		app.ServerError(w, r, err)
		return nil
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		app.Logger.Info("here")
		app.ServerError(w, r, err)
		return nil
	}

	user, err := app.Queries.GetPerson(r.Context(), userID)
	if err != nil {
		app.Logger.Info("here")
		app.ServerError(w, r, err)
		return nil
	}

	if user.ID == 0 {
		app.Logger.Info("here")
		app.invalidAuthenticationToken(w, r)
	}
	return ContextSetAuthenticatedUser(r, &user)

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
