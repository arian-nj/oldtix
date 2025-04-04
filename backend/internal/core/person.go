package core_api

import (
	"context"

	"github.com/arian-nj/master-card/back/pkg/password"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *ApiApplication) CreateBrandNewPerson(username, displayname, palin_password string) error {
	hashedPassword, err := password.Hash(palin_password)
	if err != nil {
		return err
	}
	personRow, err := app.Queries.InsertPerson(context.Background(), sqldb.InsertPersonParams{
		Username:       username,
		DisplayName:    pgtype.Text{String: displayname, Valid: true},
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return err
	}
	_, err = app.Queries.InsertUserStatistic(context.Background(), personRow.ID)
	return err

}
