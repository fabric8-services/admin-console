package migration

import (
	"database/sql"

	"github.com/fabric8-services/fabric8-common/migration"
)

// Migrate performs the database migration to update the schema to the latest level
func Migrate(db *sql.DB, catalog string) error {
	return migration.Migrate(db, catalog, migrateData{})
}

type migrateData struct {
}

func (d migrateData) Asset(name string) ([]byte, error) {
	return Asset(name)
}

func (d migrateData) AssetNameWithArgs() [][]string {
	names := [][]string{
		{"000-bootstrap.sql"},
		{"001-audit-logs"},
	}
	return names
}
