package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Migrator manages migrations
type Migrator struct {
	Pool *pgxpool.Pool
	Log  *log.Logger
}

// CreateMigrator creates a database
func CreateMigrator(hostName, port, database, user, password string, lg *log.Logger) (*Migrator, error) {
	var err error
	config, err := pgxpool.ParseConfig(fmt.Sprintf("host=%s port=%s", hostName, port))
	if err != nil {
		return nil, err
	}
	config.ConnConfig.Config.Database = database
	config.ConnConfig.Config.User = user
	config.ConnConfig.Config.Password = password
	config.MaxConns = 16
	config.MaxConnIdleTime = time.Minute
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Migrator{Pool: pool, Log: lg}, nil
}
