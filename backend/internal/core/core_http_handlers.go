package core_api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/pkg/request"
	"github.com/arian-nj/master-card/back/pkg/response"
	"github.com/arian-nj/master-card/back/pkg/validator"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

type JWTTokenOutput struct {
	Token string `json:"token"`
}

var ErrUserIsNotValid error = fmt.Errorf("user id is not valid")

func (app *ApiApplication) status(writer http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok u",
	}

	err := response.JSON(writer, http.StatusOK, data)
	if err != nil {
		app.ServerError(writer, r, err)
	}
}

func (app *ApiApplication) getLatestVersion(w http.ResponseWriter, r *http.Request) {
	// pwRow, err := app.Queries.GetVersion(r.Context(), app.ReleaseMode)
	// if err != nil {
	// 	app.ServerError(w, r, err)
	// 	return

	// }
	project_version := os.Getenv("PROJECT_VERSION")
	if project_version == "" {
		app.ServerError(w, r, fmt.Errorf("PROJECT_VERSION is empty"))
		return
	}

	data := map[string]string{
		"version": project_version,
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) getUserStatisticsData(w http.ResponseWriter, r *http.Request) {
	user_param := chi.URLParam(r, "user_id")
	user_int, err := strconv.Atoi(user_param)
	if err != nil {
		app.BadRequest(w, r, ErrUserIsNotValid)
		return
	}

	statisticsRow, err := app.Queries.GetPersonStatisticsById(r.Context(), user_int)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.NotFound(w, r)
		} else {
			app.Logger.Info("wanted id is string" + user_param)
			app.ServerError(w, r, err)
		}
		return
	}

	var output struct {
		UserID     int `json:"user_id"`
		Win        int `json:"win"`
		Lose       int `json:"lose"`
		TricksWon  int `json:"tricks_won"`
		TricksLost int `json:"tricks_lost"`
		TurnsWon   int `json:"turns_won"`
		TurnsLost  int `json:"turns_lost"`
	}

	output.UserID = statisticsRow.UserID
	output.Win = statisticsRow.Win
	output.Lose = statisticsRow.Lose
	output.TricksWon = statisticsRow.TotalTricksWon
	output.TricksLost = statisticsRow.TotalTricksLost
	output.TurnsWon = statisticsRow.TotalTurnsWon
	output.TurnsLost = statisticsRow.TotalTurnsLost

	err = response.JSON(w, http.StatusOK, output)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) getUserData(w http.ResponseWriter, r *http.Request) {
	user_param := chi.URLParam(r, "user_id")
	userid_int, err := strconv.Atoi(user_param)
	if err != nil {
		app.BadRequest(w, r, ErrUserIsNotValid)
		return
	}

	user, err := app.Queries.GetPerson(r.Context(), userid_int)
	if err != nil {
		app.BadRequest(w, r, ErrUserIsNotValid)
		return
	}

	var output struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
		Coin        int    `json:"coin"`
	}

	output.Username = "remove me"
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
	output.Username = "remove me to getMeData"
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

func createToken(userId int) *jwt.Token {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userId),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpiryDuration)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

func (app *ApiApplication) createGuest(w http.ResponseWriter, r *http.Request) {
	newRandString := randutils.GenerateRandomString(128)
	_, err := app.CreateNewGuest(newRandString)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	var output struct {
		UIDString string `json:"uid_string"`
	}
	output.UIDString = newRandString

	err = response.JSON(w, http.StatusCreated, output)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}
func (app *ApiApplication) createGuestToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UIDString string              `json:"uid_string"`
		validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}
	input.validator.CheckField(input.UIDString != "", "uid_string", "uid can't be empty")
	if input.validator.HasErrors() {
		app.FailedValidation(w, r, input.validator)
		return
	}

	guestPersonRow, err := app.Queries.GetGuestPersonByUid(r.Context(), input.UIDString)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			input.validator.AddFieldError("uid_string", "no user found")
			app.FailedValidation(w, r, input.validator)
			return
		} else {
			app.ServerError(w, r, err)
			return
		}
	}

	token := createToken(guestPersonRow.UserID)
	tokenString, err := token.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, JWTTokenOutput{
		Token: tokenString,
	})
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

// Auth

// func (app *ApiApplication) register(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username  string              `json:"username"`
// 		Password  string              `json:"password"`
// 		Validator validator.Validator `json:"-"`
// 	}
// 	err := request.DecodeJSON(w, r, &input)
// 	if err != nil {
// 		app.BadRequest(w, r, err)
// 		return
// 	}

// 	input.Validator.CheckField(len(input.Username) >= 5, "username", "to short min lenth is 5")
// 	input.Validator.CheckField(len(input.Username) < 64, "username", "to long")

// 	input.Validator.CheckField(input.Password != "", "password", "Password is required")
// 	input.Validator.CheckField(len(input.Password) >= 8, "password", "Password is too short")
// 	input.Validator.CheckField(len(input.Password) <= 72, "password", "Password is too long")
// 	// input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "password", "Password is too common")

// 	if input.Validator.HasErrors() {
// 		app.FailedValidation(w, r, input.Validator)
// 		return
// 	}
// 	userRow, err := app.Queries.GetPersonByUsername(r.Context(), input.Username)
// 	if err != nil {
// 		if !errors.Is(err, pgx.ErrNoRows) {
// 			app.ServerError(w, r, err)
// 			return
// 		}
// 	}
// 	if userRow.ID != 0 {
// 		input.Validator.AddFieldError("username", "username already exists")
// 	}

// 	if input.Validator.HasErrors() {
// 		app.FailedValidation(w, r, input.Validator)
// 		return
// 	}

// 	err = app.CreateBrandNewPerson(input.Username, input.Username, input.Password)
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return
// 	}

// 	// end
// 	w.WriteHeader(http.StatusCreated)

// }

// func (app *ApiApplication) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username  string              `json:"username"`
// 		Password  string              `json:"password"`
// 		Validator validator.Validator `json:"-"`
// 	}

// 	err := request.DecodeJSON(w, r, &input)
// 	if err != nil {
// 		app.BadRequest(w, r, err)
// 		return
// 	}

// 	person, err := app.Queries.GetPersonByUsername(r.Context(), input.Username)
// 	username_found := true
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			username_found = false
// 		} else {
// 			app.ServerError(w, r, err)
// 			return
// 		}
// 	}

// 	input.Validator.CheckField(input.Username != "", "username", "username is required")
// 	input.Validator.CheckField(username_found, "username", "username address could not be found")

// 	if username_found {
// 		passwordMatches, err := password.Matches(input.Password, person.HashedPassword)
// 		if err != nil {
// 			app.ServerError(w, r, err)
// 			return
// 		}

// 		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
// 		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
// 	}

// 	if input.Validator.HasErrors() {
// 		app.FailedValidation(w, r, input.Validator)
// 		return
// 	}

// 	token := createToken(person.ID)
// 	tokenString, err := token.SignedString([]byte(app.Config.Jwt.SECRET_KEY_FILE))
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return
// 	}

// 	data := map[string]string{
// 		"token": tokenString,
// 	}

// 	err = response.JSON(w, http.StatusOK, data)
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 	}
// }
