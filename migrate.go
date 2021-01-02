package db

import (
	"context"
	"fmt"
	"sort"

	"github.com/golang-multi-schema-migration/migrations"
	"github.com/golang-multi-schema-migration/models"
	"github.com/jackc/pgx/v4"
)

func (m *Migrator) migratePreamble(ctx context.Context) (pgx.Tx, func(err error), error) {
	tx, err := m.Pool.Begin(ctx)

	if err != nil {
		return nil, nil, err
	}

	deferFunc := func(err error) {
		if err != nil {
			rErr := tx.Rollback(ctx)
			if rErr != nil {
				m.Log.Panic("Failed to rollback migration")
			}
		}
	}

	return tx, deferFunc, err
}

// Order matters for migration rows
func (m *Migrator) getMigrationRows(ctx context.Context, schema string, isPublic bool, up bool, dateString string, tx pgx.Tx) ([]models.Migration, error) {

	migs := migrations.GetAllSchemaMigrations()
	if isPublic {
		migs = migrations.GetAllPublicMigrations()
	}
	if up {
		sort.Slice(migs, func(i, j int) bool {
			return migs[i].DateString < migs[j].DateString
		})
	} else {
		sort.Slice(migs, func(i, j int) bool {
			return migs[i].DateString > migs[j].DateString
		})
	}

	// dateString is only empty if we are migrating to the latest date
	mostRecentMigrationDateString := ""
	rows, err := tx.Query(ctx, `SELECT date_string, name FROM migrations ORDER BY date_string DESC LIMIT 1;`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	name := ""
	for rows.Next() {
		err := rows.Scan(&mostRecentMigrationDateString, &name)
		if err != nil {
			return nil, err
		}
	}

	filteredRows := []models.Migration{}
	for _, mig := range migs {

		if up {
			// if migrating up, ignore migrations that we've already done
			if mostRecentMigrationDateString != "" && mig.DateString <= mostRecentMigrationDateString {
				continue
			}

			if dateString == "" {
				filteredRows = append(filteredRows, mig)
			} else if mig.DateString <= dateString {
				filteredRows = append(filteredRows, mig)
			}
			continue
		}
		// if migrating up, ignore migrations that we haven't done
		if mostRecentMigrationDateString != "" && mig.DateString > mostRecentMigrationDateString {
			continue
		}
		if dateString == "" {
			filteredRows = append(filteredRows, mig)
		} else if mig.DateString > dateString {
			filteredRows = append(filteredRows, mig)
		}
	}

	return filteredRows, nil
}

func (m *Migrator) latestMigration(ctx context.Context, schema string, isPublic bool, up bool, dateString string, tx pgx.Tx) error {
	if !SchemaIsSafe(schema) {
		return fmt.Errorf("schema is not safe: %s", schema)
	}

	pre := fmt.Sprintf(`SET search_path TO "%s";`, schema)
	_, err := tx.Exec(ctx, pre)
	if err != nil {
		return err
	}

	if up {
		if err := m.ifMigrationTableDoesNotExistMakeIt(ctx, schema, tx); err != nil {
			return err
		}
	}
	filteredRows, err := m.getMigrationRows(ctx, schema, isPublic, up, dateString, tx)
	if err != nil {
		return err
	}
	for _, mig := range filteredRows {
		if up {
			if err := mig.Up(ctx, tx); err != nil {
				return err
			}

			_, err = tx.Exec(ctx, `INSERT INTO migrations (date_string, name) VALUES ($1,$2);`, mig.DateString, mig.Name)
			if err != nil {
				return err
			}
			continue
		}
		if err := mig.Down(ctx, tx); err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `DELETE FROM migrations WHERE name=$1;`, mig.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
