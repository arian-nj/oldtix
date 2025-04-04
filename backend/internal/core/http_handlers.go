package core_api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/pkg/password"
	"github.com/arian-nj/master-card/back/pkg/request"
	"github.com/arian-nj/master-card/back/pkg/response"
	"github.com/arian-nj/master-card/back/pkg/validator"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pascaldekloe/jwt"
)

func (app *ApiApplication) status(writer http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	err := response.JSON(writer, http.StatusOK, data)
	if err != nil {
		app.ServerError(writer, r, err)
	}
}

func (app *ApiApplication) getLatestVersion(w http.ResponseWriter, r *http.Request) {
	pwRow, err := app.Queries.GetVersion(r.Context(), app.ReleaseMode)
	if err != nil {
		app.ServerError(w, r, err)
		return

	}

	data := map[string]string{
		"version": pwRow.VersionNumber,
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) getUserData(w http.ResponseWriter, r *http.Request) {
	var UserIsNotValidErr error = fmt.Errorf("user id is not valid")
	user_param := chi.URLParam(r, "user_id")
	userid_int, err := strconv.Atoi(user_param)
	if err != nil {
		app.BadRequest(w, r, UserIsNotValidErr)
		return
	}

	user, err := app.Queries.GetPerson(r.Context(), userid_int)
	if err != nil {
		app.BadRequest(w, r, UserIsNotValidErr)
		return
	}

	var output struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
		Coin        int    `json:"coin"`
	}

	output.Username = user.Username
	output.DisplayName = user.DisplayName.String
	output.Bio = user.Bio.String
	output.Coin = user.Coin

	err = response.JSON(w, http.StatusOK, output)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) getMeData(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	var output struct {
		ID          int    `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
		Coin        int    `json:"coin"`
	}
	output.ID = user.ID
	output.Username = user.Username
	output.DisplayName = user.DisplayName.String
	output.Bio = user.Bio.String
	output.Coin = user.Coin

	err := response.JSON(w, http.StatusOK, output)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) updateUserData(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	var input struct {
		DisplayName *string             `json:"display_name"`
		Validator   validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	if input.DisplayName != nil {
		user.DisplayName.String = *input.DisplayName
		user.DisplayName.Valid = true

		input.Validator.CheckField(len(user.DisplayName.String) >= 5, "username", "to short min lenth is 5")
		input.Validator.CheckField(len(user.DisplayName.String) < 64, "username", "to long")
	}

	if input.Validator.HasErrors() {
		app.FailedValidation(w, r, input.Validator)
		return
	}

	err = app.Queries.UpdatePersonDisplayName(r.Context(), sqldb.UpdatePersonDisplayNameParams{
		DisplayName: user.DisplayName,
		ID:          user.ID,
	})
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Auth

func (app *ApiApplication) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string              `json:"username"`
		Password  string              `json:"password"`
		Validator validator.Validator `json:"-"`
	}
	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	input.Validator.CheckField(len(input.Username) >= 5, "username", "to short min lenth is 5")
	input.Validator.CheckField(len(input.Username) < 64, "username", "to long")

	input.Validator.CheckField(input.Password != "", "password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "password", "Password is too long")
	// input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "password", "Password is too common")

	if input.Validator.HasErrors() {
		app.FailedValidation(w, r, input.Validator)
		return
	}
	userRow, err := app.Queries.GetPersonByUsername(r.Context(), input.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			app.ServerError(w, r, err)
			return
		}
	}
	if userRow.ID != 0 {
		input.Validator.AddFieldError("username", "username already exists")
	}

	if input.Validator.HasErrors() {
		app.FailedValidation(w, r, input.Validator)
		return
	}

	err = app.CreateBrandNewPerson(input.Username, input.Username, input.Password)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	// end
	w.WriteHeader(http.StatusCreated)

}

func (app *ApiApplication) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string              `json:"username"`
		Password  string              `json:"password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	user, err := app.Queries.GetPersonByUsername(r.Context(), input.Username)
	username_found := true
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			username_found = false
		} else {
			app.ServerError(w, r, err)
			return
		}
	}

	input.Validator.CheckField(input.Username != "", "username", "username is required")
	input.Validator.CheckField(username_found, "username", "username address could not be found")

	if username_found {
		passwordMatches, err := password.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.FailedValidation(w, r, input.Validator)
		return
	}

	var claims jwt.Claims
	claims.Subject = strconv.Itoa(int(user.ID))

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.Config.BaseURL
	claims.Audiences = []string{app.Config.BaseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.Config.Jwt.SecretKey))
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.ServerError(w, r, err)
	}
}
