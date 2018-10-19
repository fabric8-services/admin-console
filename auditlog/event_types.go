package auditlog

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// UserSearch the UUID of the event for the user search action
var UserSearch uuid.UUID

func init() {
	var err error
	UserSearch, err = uuid.FromString("7aea0277-d6fa-4df9-8224-a27fa4096ec7")
	if err != nil {
		panic(fmt.Sprintf("UserSearch event type ID os not an UUID: %v", err))
	}
}
