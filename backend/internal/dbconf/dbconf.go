package dbconf

import (
	"context"
	"os"

	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDB() (*sqldb.Queries, *pgxpool.Pool, error) {

	DbConnStr := os.Getenv("DATABASE_URL")
	if DbConnStr == "" {
		panic(DbConnStr)
	}
	conn, err := pgxpool.New(context.Background(), DbConnStr)
	if err != nil {
		return nil, nil, err
	}

	queries := sqldb.New(conn)

	return queries, conn, err
}
