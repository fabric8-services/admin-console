package migration

import (
	"database/sql"

	"github.com/fabric8-services/fabric8-common/migration"
)

// Migrate performs the database migration to update the schema to the latest level
func Migrate(db *sql.DB, catalog string) error {
	return migration.Migrate(db, catalog, Steps())
}

// Scripts the structure that provides the SQL scripts to migrate the database schema
type Scripts [][]string

// Steps returns the array of scripts to run to migrate the database
func Steps() Scripts {
	return [][]string{
		{"000-bootstrap.sql"},
		{"001-audit-log.sql"},
		{"002-tenant-updates-event-types.sql"},
		{"003-audit-log-with-username.sql"},
		{"004-deactivation-event-types.sql"},
		{"005-list-audit-logs-event-type.sql"},
	}
}

// Asset returns the content of a file given its name
func (d Scripts) Asset(name string) ([]byte, error) {
	return Asset(name)
}

func (d Scripts) AssetNameWithArgs() [][]string {
	return d
}
