package application

import (
	"runtime/debug"
	"time"

	"github.com/fabric8-services/fabric8-common/log"

	"github.com/pkg/errors"
)

var databaseTransactionTimeout = 5 * time.Minute

// SetDatabaseTransactionTimeout sets the global timeout variable to the given
// duration.
func SetDatabaseTransactionTimeout(t time.Duration) {
	databaseTransactionTimeout = t
}

// Transactional executes the given function in a transaction. If todo returns an error, the transaction is rolled back
func Transactional(db DB, todo func(f Application) error) error {
	var tx Transaction
	var err error
	if tx, err = db.BeginTransaction(); err != nil {
		log.Error(nil, map[string]interface{}{
			"err": err,
		}, "database BeginTransaction failed!")

		return errors.WithStack(err)
	}

	return func() error {
		errorChan := make(chan error, 1)
		txTimeout := time.After(databaseTransactionTimeout)

		go func(tx Transaction) {
			defer func() {
				if err := recover(); err != nil {
					errorChan <- errors.Errorf("recovered %v. stack: %s", err, debug.Stack())
				}
			}()
			errorChan <- todo(tx)
		}(tx)

		select {
		case err := <-errorChan:
			if err != nil {
				log.Debug(nil, map[string]interface{}{"error": err}, "Rolling back the transaction...")
				tx.Rollback()
				log.Error(nil, map[string]interface{}{
					"err": err,
				}, "database transaction failed!")
				return errors.WithStack(err)
			}

			tx.Commit()
			log.Debug(nil, nil, "Commit the transaction!")
			return nil
		case <-txTimeout:
			log.Debug(nil, nil, "Rolling back the transaction...")
			tx.Rollback()
			log.Error(nil, nil, "database transaction timeout!")
			return errors.New("database transaction timeout")
		}
	}()
}
