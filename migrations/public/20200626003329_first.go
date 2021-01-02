package public

import (
	"context"

	"github.com/golang-multi-schema-migration/models"

	"github.com/jackc/pgx/v4"
)

func init() {
	up := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			CREATE TABLE organizations (
				id uuid NOT NULL,
				org_name text NOT NULL,
				PRIMARY KEY(id)
			);

			CREATE TABLE users (
				id uuid NOT NULL,
				email text NOT NULL,
				org_id uuid NOT NULL,
				PRIMARY KEY(id),
				FOREIGN KEY(org_id) REFERENCES organizations (id)
		);

		ALTER TABLE organizations
		ADD CONSTRAINT o_name UNIQUE (org_name);

		ALTER TABLE users
		ADD CONSTRAINT u_email UNIQUE (email, org_id);
			
		`)
		return err
	}

	down := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
		ALTER TABLE users
		DROP CONSTRAINT u_email;

		ALTER TABLE organizations
		DROP CONSTRAINT o_name;
	
	  	DROP TABLE users;
	  	DROP TABLE organizations;
		`)
		return err
	}

	migrations = append(migrations, models.Migration{
		DateString: "20200626003329",
		Name:       "first",
		Up:         up,
		Down:       down,
	})
}
