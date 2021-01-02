package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func (m *Migrator) migrateDownFully(ctx context.Context, schema string, isPublic bool, tx pgx.Tx) error {
	return m.latestMigration(ctx, schema, isPublic, false, "", tx)
}

func (m *Migrator) migrateDownToDateString(ctx context.Context, schema string, dateString string, isPublic bool, tx pgx.Tx) error {
	return m.latestMigration(ctx, schema, isPublic, false, dateString, tx)
}
