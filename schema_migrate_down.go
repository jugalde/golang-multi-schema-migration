package db

import (
	"context"
)

// MigrateSchemaDownFully does as the name implies
func (m *Migrator) MigrateSchemaDownFully(ctx context.Context, schema string) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateDownFully(ctx, schema, false, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

// MigrateSchemaDownToDateString does as the name implies
func (m *Migrator) MigrateSchemaDownToDateString(ctx context.Context, schema string, dateString string) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateDownToDateString(ctx, schema, dateString, false, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}
