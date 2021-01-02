package public

import "github.com/golang-multi-schema-migration/models"

var migrations = []models.Migration{}

// GetMigrations exposes this packages migrations
func GetMigrations() []models.Migration {
	return migrations
}
