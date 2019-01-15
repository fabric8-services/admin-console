package migration_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/fabric8-services/admin-console/migration"

	migrationsupport "github.com/fabric8-services/fabric8-common/migration"

	"github.com/fabric8-services/fabric8-common/gormsupport"
	"github.com/fabric8-services/fabric8-common/resource"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	dbName      = "test"
	defaultHost = "localhost"
	defaultPort = "5435"
)

type MigrationTestSuite struct {
	suite.Suite
}

const (
	databaseName = "test"
)

var (
	sqlDB *sql.DB
	host  string
	port  string
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}

func (s *MigrationTestSuite) SetupTest() {
	resource.Require(s.T(), resource.Database)

	host = os.Getenv("ADMIN_POSTGRES_HOST")
	if host == "" {
		host = defaultHost
	}
	port = os.Getenv("ADMIN_POSTGRES_PORT")
	if port == "" {
		port = defaultPort
	}

	dbConfig := fmt.Sprintf("host=%s port=%s user=postgres password=mysecretpassword sslmode=disable connect_timeout=5", host, port)

	db, err := sql.Open("postgres", dbConfig)
	require.NoError(s.T(), err, "cannot connect to database: %s", dbName)
	defer db.Close()

	_, err = db.Exec("DROP DATABASE " + dbName)
	if err != nil && !gormsupport.IsInvalidCatalogName(err) {
		require.NoError(s.T(), err, "failed to drop database '%s'", dbName)
	}

	_, err = db.Exec("CREATE DATABASE " + dbName)
	require.NoError(s.T(), err, "failed to create database '%s'", dbName)
}

func (s *MigrationTestSuite) TestMigrate() {
	dbConfig := fmt.Sprintf("host=%s port=%s user=postgres password=mysecretpassword dbname=%s sslmode=disable connect_timeout=5",
		host, port, dbName)
	var err error
	sqlDB, err = sql.Open("postgres", dbConfig)
	require.NoError(s.T(), err, "cannot connect to DB '%s'", dbName)
	defer sqlDB.Close()

	gormDB, err := gorm.Open("postgres", dbConfig)
	require.NoError(s.T(), err, "cannot connect to DB '%s'", dbName)
	defer gormDB.Close()

	dialect := gormDB.Dialect()
	dialect.SetDB(sqlDB)

	s.T().Run("checkMigration001", checkMigration001)

}

func checkMigration001(t *testing.T) {
	err := migrationsupport.Migrate(sqlDB, databaseName, migration.Steps()[:2])
	require.NoError(t, err)

	t.Run("insert without id", func(t *testing.T) {
		// when
		_, err := sqlDB.Exec("INSERT INTO audit_log (identity_id, event_type_id, event_params) VALUES (uuid_generate_v4(),'7aea0277-d6fa-4df9-8224-a27fa4096ec7', '{}')")
		// ok to insert audit log record without ID
		require.NoError(t, err)
	})

	t.Run("insert without identity_id", func(t *testing.T) {
		// when
		_, err := sqlDB.Exec("INSERT INTO audit_log (event_type_id, event_params) VALUES ('7aea0277-d6fa-4df9-8224-a27fa4096ec7', '{}')")
		// ok to insert audit log record without ID
		require.Error(t, err)
	})

	t.Run("insert without event_type_id", func(t *testing.T) {
		// when
		_, err := sqlDB.Exec("INSERT INTO audit_log (identity_id, event_params) VALUES (uuid_generate_v4(),'{}')")
		// ok to insert audit log record without ID
		require.Error(t, err)
	})

	t.Run("insert with invalid event_type_id", func(t *testing.T) {
		// when
		_, err := sqlDB.Exec("INSERT INTO audit_log (identity_id, event_type_id, event_params) VALUES (uuid_generate_v4(),'7aea0277-cafe-cafe-cafe-a27fa4096ec7', '{}')")
		// ok to insert audit log record without ID
		require.Error(t, err)
	})

	t.Run("insert without event_params", func(t *testing.T) {
		// when
		_, err := sqlDB.Exec("INSERT INTO audit_log (identity_id, event_type_id) VALUES (uuid_generate_v4(),'7aea0277-d6fa-4df9-8224-a27fa4096ec7')")
		// ok to insert audit log record without ID
		require.Error(t, err)
	})

}
