package models

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// Migration are used to track migrations
type Migration struct {
	ID         string
	DateString string
	Name       string
	Up         func(ctx context.Context, tx pgx.Tx) error
	Down       func(ctx context.Context, tx pgx.Tx) error
}
