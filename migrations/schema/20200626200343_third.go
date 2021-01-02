package schema

import (
	"context"

	"github.com/golang-multi-schema-migration/models"

	"github.com/jackc/pgx/v4"
)

func init() {
	up := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
		ALTER TABLE users ADD COLUMN views integer;
			`)
		return err
	}

	down := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
		ALTER TABLE users DROP COLUMN views;
		`)
		return err
	}

	migrations = append(migrations, models.Migration{
		DateString: "20200626200343",
		Name:       "third",
		Up:         up,
		Down:       down,
	})
}
