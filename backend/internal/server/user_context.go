package server

import (
	"context"
	"net/http"

	"github.com/arian-nj/master-card/back/sqldb"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func ContextSetAuthenticatedUser(r *http.Request, user *sqldb.Person) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetAuthenticatedUser(r *http.Request) *sqldb.Person {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*sqldb.Person)
	if !ok {
		return nil
	}

	return user
}
