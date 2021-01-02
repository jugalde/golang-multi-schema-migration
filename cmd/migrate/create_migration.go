package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

const timeFormat = "20060102150405"

var template = `package %s

import (
	"context"

	"github.com/golang-multi-schema-migration/models"

	"github.com/jackc/pgx/v4"
)

func init() {
	up := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "")
		return err
	}

	down := func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "")
		return err
	}

	migrations = append(migrations, models.Migration{
		DateString: "%s",
		Name:       "%s",
		Up:         up,
		Down:       down,
	})
}
`

// Create creates a migration file
func Create(directory, name string, isPublic bool) error {
	version := time.Now().UTC().Format(timeFormat)
	fullname := fmt.Sprintf("%s_%s", version, name)

	level := "public"
	if !isPublic {
		level = "schema"
	}
	directory += "/" + level

	filename := path.Join(directory, fmt.Sprintf("%s.go", fullname))
	fmt.Printf("Creating %s...\n", filename)

	return ioutil.WriteFile(filename, []byte(fmt.Sprintf(template, level, version, name)), 0644)
}
