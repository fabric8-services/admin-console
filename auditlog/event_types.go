package auditlog

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// UserSearch the UUID of the event for the "user search" action
var UserSearch uuid.UUID

// ShowTenantUpdate the UUID of the event for the "show tenant update" action
var ShowTenantUpdate uuid.UUID

// StartTenantUpdate the UUID of the event for the "start tenant update" action
var StartTenantUpdate uuid.UUID

// StopTenantUpdate the UUID of the event for the "stop tenant update" action
var StopTenantUpdate uuid.UUID

func init() {
	var err error
	UserSearch, err = uuid.FromString("7aea0277-d6fa-4df9-8224-a27fa4096ec7")
	if err != nil {
		panic(fmt.Sprintf("UserSearch event type ID is not an UUID: %v", err))
	}
	ShowTenantUpdate, err = uuid.FromString("a2633717-12f0-4edd-bcd9-bbdf900a8ec5")
	if err != nil {
		panic(fmt.Sprintf("ShowTenantUpdate event type ID is not an UUID: %v", err))
	}
	StartTenantUpdate, err = uuid.FromString("2d51de09-2ab7-4e15-9e0d-030a71756c8d")
	if err != nil {
		panic(fmt.Sprintf("StartTenantUpdate event type ID is not an UUID: %v", err))
	}
	StopTenantUpdate, err = uuid.FromString("3dd22424-27b6-494a-a550-9611bfe41cac")
	if err != nil {
		panic(fmt.Sprintf("StopTenantUpdate event type ID is not an UUID: %v", err))
	}
}
