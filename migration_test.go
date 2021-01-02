package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func createOrg(ctx context.Context, t *testing.T, db *Migrator) (string, string, string) {
	orgID := uuid.New().String()
	orgName := uuid.New().String()
	schemaID := strings.ReplaceAll(orgID, "-", "_")

	_, err := db.Pool.Exec(ctx, `INSERT INTO public.organizations(id, org_name) VALUES ($1, $2)`, orgID, orgName)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Pool.Exec(ctx, fmt.Sprintf(`CREATE SCHEMA "%s"`, schemaID))
	if err != nil {
		t.Fatal(err)
	}
	return orgID, orgName, schemaID
}
func dropSchema(ctx context.Context, t *testing.T, db *Migrator, schemaID string) {
	_, err := db.Pool.Exec(ctx, fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schemaID))
	if err != nil {
		t.Fatal(err)
	}
}
func resetPublic(ctx context.Context, t *testing.T, db *Migrator) {
	_, err := db.Pool.Exec(ctx, `DROP SCHEMA public CASCADE;`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Pool.Exec(ctx, `CREATE SCHEMA public;`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMigrateUpAndDownToLatestPublic(t *testing.T) {
	db, err := CreateMigrator("127.0.0.1", "8001", "db", "user", "pass", &log.Logger{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Pool.Close()

	ctx := context.Background()
	defer resetPublic(ctx, t, db)

	if err := db.MigratePublicUpToLatest(ctx); err != nil {
		t.Fatal(err)
	}

	orgID, _, schema := createOrg(ctx, t, db)
	defer dropSchema(ctx, t, db, schema)

	id := uuid.New().String()
	_, err = db.Pool.Exec(ctx, "INSERT INTO users(id,views,name,email,org_id) VALUES ($1,23,'name','johndough@gmail.com',$2)", id, orgID)
	if err != nil {
		t.Fatal(err)
	}

	// views is only present in the latest migration
	rows1, err := db.Pool.Query(ctx, "SELECT views FROM users;")
	defer rows1.Close()
	if err != nil {
		t.Fatal(err)
	}

	for rows1.Next() {
		rowErr := rows1.Err()
		if rowErr != nil {
			t.Fatal(rowErr)
		}
	}

	if err := db.MigratePublicDownFully(ctx); err != nil {
		t.Fatal(err)
	}

	// Error should be thrown when trying to get field because table doesn't exist
	rows2, err := db.Pool.Query(ctx, "SELECT id FROM users;")
	defer rows2.Close()
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatal(errors.New("users table should not longer be present after migrate down"))
	}

}

func TestMigrateUpAndDownToLatestSchema(t *testing.T) {
	db, err := CreateMigrator("127.0.0.1", "8001", "db", "user", "pass", &log.Logger{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Pool.Close()

	ctx := context.Background()
	// "public.organizations" needs to be created to start creating schemas, so we need to migrate public to at least the first migration
	if err := db.MigratePublicUpToLatest(ctx); err != nil {
		t.Fatal(err)
	}
	defer resetPublic(ctx, t, db)

	_, _, schema := createOrg(ctx, t, db)
	defer dropSchema(ctx, t, db, schema)

	if err := db.MigrateSchemaUpToLatest(ctx, schema); err != nil {
		t.Fatal(err)
	}

	// views is only present in the latest migration
	rows1, err := db.Pool.Query(ctx, "SELECT views FROM users;")
	defer rows1.Close()
	if err != nil {
		t.Fatal(err)
	}

	if err = db.MigrateSchemaDownFully(ctx, schema); err != nil {
		t.Fatal(err)
	}

	// Error should be thrown when trying to get field because table doesn't exist
	rows2, err := db.Pool.Query(ctx, "SELECT id FROM users;")
	defer rows2.Close()
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatal(errors.New("users table should not longer be present after migrate down"))
	}
}

func TestMigrateUpAndDownToDateStringPublic(t *testing.T) {
	db, err := CreateMigrator("127.0.0.1", "8001", "db", "user", "pass", &log.Logger{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Pool.Close()

	ctx := context.Background()

	defer resetPublic(ctx, t, db)

	secondDateString := "20200626004429"

	if err := db.MigratePublicUpToDateString(ctx, secondDateString); err != nil {
		t.Fatal(err)
	}

	// name is only present in the second migration
	rows1, err := db.Pool.Query(ctx, "SELECT name FROM users;")
	defer rows1.Close()
	if err != nil {
		t.Fatal(err)
	}

	firstDateString := "20200626003329"

	if err := db.MigratePublicDownToDateString(ctx, firstDateString); err != nil {
		t.Fatal(err)
	}

	// Name shouldn't exist but id should
	rows2, err := db.Pool.Query(ctx, "SELECT name FROM users;")
	defer rows2.Close()
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatal(errors.New("users table should not longer be present after migrate down"))
	}

	rows3, err := db.Pool.Query(ctx, "SELECT id FROM users;")
	defer rows3.Close()
	if err != nil {
		t.Fatal(err)
	}

}

func TestMigrateUpAndDownToDateStringSchema(t *testing.T) {
	db, err := CreateMigrator("127.0.0.1", "8001", "db", "user", "pass", &log.Logger{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Pool.Close()

	ctx := context.Background()
	// "public.organizations" needs to be created to start creating schemas, so we need to migrate public to at least the first migration
	if err := db.MigratePublicUpToLatest(ctx); err != nil {
		t.Fatal(err)
	}
	defer resetPublic(ctx, t, db)

	_, _, schema := createOrg(ctx, t, db)
	defer dropSchema(ctx, t, db, schema)

	secondDateString := "20200626200342"

	if err := db.MigrateSchemaUpToDateString(ctx, schema, secondDateString); err != nil {
		t.Fatal(err)
	}

	// name is only present in the second migration
	rows1, err := db.Pool.Query(ctx, "SELECT name FROM users;")
	defer rows1.Close()
	if err != nil {
		t.Fatal(err)
	}

	firstDateString := "20200626200341"

	if err := db.MigrateSchemaDownToDateString(ctx, schema, firstDateString); err != nil {
		t.Fatal(err)
	}

	rows2, err := db.Pool.Query(ctx, "SELECT name FROM users;")
	defer rows2.Close()
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatal(errors.New("users table should not longer be present after migrate down"))
	}

	// Migrate down fully to clean database
	err = db.MigrateSchemaDownFully(ctx, schema)
	if err != nil {
		t.Fatal(err)
	}
}
