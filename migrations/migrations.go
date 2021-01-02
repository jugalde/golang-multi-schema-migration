package migrations

import (
	"context"

	"github.com/golang-multi-schema-migration/migrations/public"
	"github.com/golang-multi-schema-migration/migrations/schema"
	"github.com/golang-multi-schema-migration/models"
	"github.com/jackc/pgx/v4"
)

// GetAllPublicMigrations returns all public migrations
func GetAllPublicMigrations() []models.Migration {
	return public.GetMigrations()
}

// GetAllSchemaMigrations returns all schema migrations
func GetAllSchemaMigrations() []models.Migration {
	return schema.GetMigrations()
}

// CreateMigrationTable needs to be done before any migration is run
func CreateMigrationTable(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
		CREATE TABLE migrations(
			id uuid DEFAULT public.uuid_generate_v4 (),
			date_string text NOT NULL,
			name text NOT NULL UNIQUE,
			PRIMARY KEY (id)
		);
`)
	return err
}
