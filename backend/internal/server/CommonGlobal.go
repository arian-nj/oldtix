package server

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/arian-nj/master-card/back/internal/dbconf"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	// err := run(logger)
	// if err != nil {
	// 	trace := string(debug.Stack())
	// 	logger.Error(err.Error(), "trace", trace)
	// 	os.Exit(1)
	// }
}

type Config struct {
	HTTPPort int

	BaseURL string
	Jwt     struct {
		SecretKey []byte
	}
	DatabaseUrl string
}

type CommonGlobals struct {
	Config  *Config
	Logger  *zap.Logger
	Queries *sqldb.Queries
	wg      sync.WaitGroup
}

func NewCommonGlobals(http_port string) (*CommonGlobals, *pgxpool.Pool, error) {
	// read configs
	Glob := &CommonGlobals{}

	cfg := &Config{}
	Glob.Config = cfg

	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderConfig

	logger, err := config.Build()

	// logger, err :=
	if err != nil {
		return nil, nil, err
	}
	Glob.Logger = logger

	// err = godotenv.Load()
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("error loading .env file")
	// }

	cfg, err = readConfigs(http_port)
	if err != nil {
		return nil, nil, err
	}
	Glob.Config = cfg

	queries, poll, err := dbconf.SetupDB(cfg.DatabaseUrl)
	if err != nil {
		return nil, nil, err
	}
	Glob.Queries = queries

	return Glob, poll, nil
}

func readConfigs(http_port string) (*Config, error) {
	var cfg *Config = new(Config)
	port := os.Getenv(http_port)
	if port == "" {
		return nil, fmt.Errorf("can't read port from .env %s", port)
	}

	port_int, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	cfg.HTTPPort = port_int

	base_url := os.Getenv("BASE_URL")
	if base_url == "" {
		return nil, fmt.Errorf("can't base url from .env : %s", base_url)
	}
	cfg.BaseURL = base_url

	secret_key := os.Getenv("SecretKey")
	if secret_key == "" {
		return nil, fmt.Errorf("can't read secret key from .env : %s", secret_key)
	}
	cfg.Jwt.SecretKey = []byte(secret_key)

	database_url := os.Getenv("DATABASE_URL")
	log.Println("database url is ", database_url)
	if database_url == "" {
		return nil, fmt.Errorf("can't read database url from .env : %s", database_url)
	}
	cfg.DatabaseUrl = database_url

	return cfg, nil
}
