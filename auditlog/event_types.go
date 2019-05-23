package auditlog

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

const (
	// UserDeactivationNotificationEvent the type name for a user deactivation notification
	UserDeactivationNotificationEvent = "user_deactivation_notification"
	// UserDeactivationEvent the type name for a user deactivation
	UserDeactivationEvent = "user_deactivation"
)

// UserSearch the UUID of the event for the "user search" action
var UserSearch uuid.UUID

// ShowTenantUpdate the UUID of the event for the "show tenant update" action
var ShowTenantUpdate uuid.UUID

// StartTenantUpdate the UUID of the event for the "start tenant update" action
var StartTenantUpdate uuid.UUID

// StopTenantUpdate the UUID of the event for the "stop tenant update" action
var StopTenantUpdate uuid.UUID

// UserDeactivationNotification the UUID of the event when a user is notified before account deactivation
var UserDeactivationNotification uuid.UUID

// UserDeactivation the UUID of the event when a user account is deactivated
var UserDeactivation uuid.UUID

// EventTypes the event types indexed buy their name. At least, those that can be created from an endpoint
var EventTypes map[string]uuid.UUID

func init() {
	var err error
	EventTypes = map[string]uuid.UUID{}

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

	UserDeactivationNotification, err = uuid.FromString("9f924fc3-403b-4167-b20c-a543adc4ff3c")
	if err != nil {
		panic(fmt.Sprintf("UserDeactivationNotification event type ID is not an UUID: %v", err))
	}
	EventTypes[UserDeactivationNotificationEvent] = UserDeactivationNotification
	UserDeactivation, err = uuid.FromString("777ede15-4b18-4720-bada-1519d1915f2e")
	if err != nil {
		panic(fmt.Sprintf("UserDeactivation event type ID is not an UUID: %v", err))
	}
	EventTypes[UserDeactivationEvent] = UserDeactivation

}
