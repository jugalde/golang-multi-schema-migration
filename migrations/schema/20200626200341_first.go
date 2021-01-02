package schema

import (
	"context"

	"github.com/golang-multi-schema-migration/models"

	"github.com/jackc/pgx/v4"
)

func init() {
	up := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
      CREATE TABLE users (
				id uuid DEFAULT public.uuid_generate_v4 (),
				email text UNIQUE NOT NULL,
				PRIMARY KEY(id)
			);`)
		return err
	}

	down := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			DROP TABLE users;
		`)
		return err
	}

	migrations = append(migrations, models.Migration{
		DateString: "20200626200341",
		Name:       "first",
		Up:         up,
		Down:       down,
	})
}
