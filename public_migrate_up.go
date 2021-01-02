package db

import (
	"context"
)

// MigratePublicUpToLatest does as the name implies
func (m *Migrator) MigratePublicUpToLatest(ctx context.Context) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateUpToLatest(ctx, "public", true, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

// MigratePublicUpToDateString does as the name implies
func (m *Migrator) MigratePublicUpToDateString(ctx context.Context, dateString string) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateUpToDateString(ctx, "public", dateString, true, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}
