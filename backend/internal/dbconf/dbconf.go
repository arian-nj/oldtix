package dbconf

import (
	"context"

	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDB(db_url string) (*sqldb.Queries, *pgxpool.Pool, error) {

	conn, err := pgxpool.New(context.Background(), db_url)
	if err != nil {
		return nil, nil, err
	}

	queries := sqldb.New(conn)

	return queries, conn, err
}
