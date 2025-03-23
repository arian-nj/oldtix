package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/arian-nj/master-card/back/pkg/response"
	"github.com/arian-nj/master-card/back/pkg/validator"
)

func (app *CommonGlobals) ReportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		// method  = r.Method
		// url     = r.URL.String()
		trace = string(debug.Stack())
	)

	// requestAttrs := slog.Group("request", "method", method, "url", url)
	// app.Logger.Error(message, requestAttrs, "trace", trace)
	// log.Println(message, method, url)
	app.Logger.Error(message)
	log.Println(trace)
}
func (app *CommonGlobals) ReportError(err error) {
	var (
		message = err.Error()
		trace   = string(debug.Stack())
	)

	// requestAttrs := slog.Group("request", "method", method, "url", url)
	// app.Logger.Error(message, requestAttrs, "trace", trace)
	app.Logger.Error(message)
	log.Println(trace)
}

func (app *CommonGlobals) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.ReportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *CommonGlobals) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.ReportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *CommonGlobals) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *CommonGlobals) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *CommonGlobals) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (app *CommonGlobals) FailedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.ServerError(w, r, err)
	}
}

func (app *CommonGlobals) invalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	app.errorMessage(w, r, http.StatusUnauthorized, "Invalid authentication token", headers)
}

func (app *CommonGlobals) AuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	app.errorMessage(w, r, http.StatusUnauthorized, "You must be authenticated to access this resource", nil)
}
