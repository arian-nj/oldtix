package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/arian-nj/master-card/back/internal/dbconf"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/internal/version"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	httpPort string
	baseURL  string
	jwt      struct {
		secretKey string
	}
	databaseUrl string
}

type Application struct {
	config      config
	logger      *slog.Logger
	Queries     *sqldb.Queries
	wg          sync.WaitGroup
	eventRouter *socket.HandlerMap
}

func run(logger *slog.Logger) error {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	err = readConfigs(&cfg)
	if err != nil {
		return err
	}

	showVersion := flag.Bool("version", false, "display version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	app := &Application{
		config: cfg,
		logger: logger,
	}

	queries, conn, err := dbconf.SetupDB()
	if err != nil {
		return err
	}
	defer conn.Close()
	app.Queries = queries

	return app.serveHTTP()
}

func readConfigs(cfg *config) error {

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		return fmt.Errorf("can't read port from .env %s", port)
	}
	cfg.httpPort = port

	base_url := os.Getenv("BASE_URL")
	if base_url == "" {
		return fmt.Errorf("can't base url port from base url %s", base_url)
	}
	cfg.baseURL = base_url

	secret_key := os.Getenv("SecretKey")
	if secret_key == "" {
		return fmt.Errorf("can't read secret key from .env", secret_key)
	}
	cfg.jwt.secretKey = secret_key

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		return fmt.Errorf("can't read database url from .env", database_url)
	}
	cfg.databaseUrl = database_url

	// databaseUrl

	return nil
}
