package db

import (
	"context"
)

// MigratePublicDownFully does as the name implies
func (m *Migrator) MigratePublicDownFully(ctx context.Context) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateDownFully(ctx, "public", true, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err

}

// MigratePublicDownToDateString does as the name implies
func (m *Migrator) MigratePublicDownToDateString(ctx context.Context, dateString string) error {
	tx, deferFunc, err := m.migratePreamble(ctx)
	if err != nil {
		return err
	}

	defer deferFunc(err)
	err = m.migrateDownToDateString(ctx, "public", dateString, true, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err

}
