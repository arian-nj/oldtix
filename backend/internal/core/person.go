package core_api

import (
	"context"

	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *ApiApplication) CreateBrandNewPerson(displayname string) (*sqldb.Person, error) {
	personRow, err := app.Queries.InsertPerson(context.Background(), sqldb.InsertPersonParams{
		DisplayName: pgtype.Text{String: displayname, Valid: true},
		Coin:        50,
	})
	if err != nil {
		return &personRow, err
	}
	_, err = app.Queries.InsertUserStatistic(context.Background(), personRow.ID)
	return &personRow, err
}

func (app *ApiApplication) CreateNewGuest(uid_string string) (*sqldb.GuestPerson, error) {
	newPersonRow, err := app.CreateBrandNewPerson("guest")
	if err != nil {
		return nil, err
	}

	guestPersonRow, err := app.Queries.InsertGuestPerson(context.Background(), sqldb.InsertGuestPersonParams{
		UidString: uid_string,
		UserID:    newPersonRow.ID,
	})
	return &guestPersonRow, err

}
