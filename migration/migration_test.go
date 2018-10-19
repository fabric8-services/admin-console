package migration

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/fabric8-services/admin-console/configuration"
	"github.com/fabric8-services/fabric8-common/resource"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentMigrations(t *testing.T) {
	resource.Require(t, resource.Database)

	config, err := configuration.New()
	if err != nil {
		panic(fmt.Errorf("Failed to setup the configuration: %s", err.Error()))
	}

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			db, err := sql.Open("postgres", config.GetPostgresConfigString())
			if err != nil {
				t.Fatalf("Cannot connect to DB: %s\n", err)
			}
			err = Migrate(db, config.GetPostgresDatabase())
			assert.Nil(t, err)
		}()

	}
	wg.Wait()
}
