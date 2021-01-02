package db

import (
	"context"
)

// MigrateSchemaUpToLatest does as the name implies
func (m *Migrator) MigrateSchemaUpToLatest(ctx context.Context, schema string) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateUpToLatest(ctx, schema, false, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

// MigrateSchemaUpToDateString does as the name implies
func (m *Migrator) MigrateSchemaUpToDateString(ctx context.Context, schema string, dateString string) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateUpToDateString(ctx, schema, dateString, false, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}
