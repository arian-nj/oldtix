package server

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/arian-nj/master-card/back/internal/dbconf"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Run() {
	// err := run(logger)
	// if err != nil {
	// 	trace := string(debug.Stack())
	// 	logger.Error(err.Error(), "trace", trace)
	// 	os.Exit(1)
	// }
}

type config struct {
	HttpPort string
	BaseURL  string
	Jwt      struct {
		SecretKey string
	}
	DatabaseUrl string
}

type CommonGlobals struct {
	Config  *config
	Logger  *slog.Logger
	Queries *sqldb.Queries
	wg      sync.WaitGroup
}

func NewCommonGlobals() (*CommonGlobals, *pgxpool.Pool, error) {
	// read configs
	Glob := &CommonGlobals{}

	cfg := &config{}
	Glob.Config = cfg

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	Glob.Logger = logger

	err := godotenv.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading .env file")
	}

	err = readConfigs(cfg)
	if err != nil {
		return nil, nil, err
	}

	queries, poll, err := dbconf.SetupDB()
	if err != nil {
		return nil, nil, err
	}
	Glob.Queries = queries

	return Glob, poll, nil
}

func readConfigs(cfg *config) error {

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		return fmt.Errorf("can't read port from .env %s", port)
	}
	cfg.HttpPort = port

	base_url := os.Getenv("BASE_URL")
	if base_url == "" {
		return fmt.Errorf("can't base url from .env : %s", base_url)
	}
	cfg.BaseURL = base_url

	secret_key := os.Getenv("SecretKey")
	if secret_key == "" {
		return fmt.Errorf("can't read secret key from .env : %s", secret_key)
	}
	cfg.Jwt.SecretKey = secret_key

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		return fmt.Errorf("can't read database url from .env : %s", database_url)
	}
	cfg.DatabaseUrl = database_url

	return nil
}
