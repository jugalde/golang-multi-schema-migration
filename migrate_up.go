package db

import (
	"context"
	"strings"

	"github.com/golang-multi-schema-migration/migrations"

	"github.com/jackc/pgx/v4"
)

func (m *Migrator) migrateUpToLatest(ctx context.Context, schema string, isPublic bool, tx pgx.Tx) error {
	return m.latestMigration(ctx, schema, isPublic, true, "", tx)
}

func (m *Migrator) migrateUpToDateString(ctx context.Context, schema string, dateString string, isPublic bool, tx pgx.Tx) error {
	return m.latestMigration(ctx, schema, isPublic, true, dateString, tx)
}

func (m *Migrator) ifMigrationTableDoesNotExistMakeIt(ctx context.Context, schema string, tx pgx.Tx) error {
	rows, err := tx.Query(ctx, `
   SELECT * FROM information_schema.tables 
   WHERE  table_schema = $1
   AND    table_name   = 'migrations';`, schema)
	defer rows.Close()
	// no err if migrations migsExists
	if err == nil {
		for rows.Next() {
			vals, err := rows.Values()
			if err != nil {
				return err
			}
			if len(vals) != 0 {
				return nil
			}
		}
	} else if !strings.Contains(err.Error(), "does not migsExist") {
		return err
	}

	err = migrations.CreateMigrationTable(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}
