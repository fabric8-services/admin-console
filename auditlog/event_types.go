package auditlog

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

const (
	// UserSearchEvent the name of the "search user" event
	UserSearchEvent = "user_search"
	// ShowTenantUpdateEvent the name of the "show tenant update" event
	ShowTenantUpdateEvent = "show_tenant_update"
	// StartTenantUpdateEvent the name of the "start tenant update" event
	StartTenantUpdateEvent = "start_tenant_update"
	// StopTenantUpdateEvent the name of the "stop tenant update" event
	StopTenantUpdateEvent = "stop_tenant_update"
	// UserDeactivationNotificationEvent the name of the "user deactivation notification" event
	UserDeactivationNotificationEvent = "user_deactivation_notification"
	// UserDeactivationEvent the name of the "user deactivation" event
	UserDeactivationEvent = "user_deactivation"
	// ListAuditLogsEvent the name of the "list audit logs" event
	ListAuditLogsEvent = "list_audit_logs"
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

// ListAuditLogs the UUID of the event for the "list audit logs" action
var ListAuditLogs uuid.UUID

// EventTypesByName the event types indexed by their name. At least, those that can be created from an endpoint
var EventTypesByName map[string]uuid.UUID

// EventTypesByID the event types indexed by their UUID.
var EventTypesByID map[uuid.UUID]string

func init() {
	var err error
	EventTypesByName = map[string]uuid.UUID{}
	EventTypesByID = map[uuid.UUID]string{}

	UserSearch, err = uuid.FromString("7aea0277-d6fa-4df9-8224-a27fa4096ec7")
	if err != nil {
		panic(fmt.Sprintf("UserSearch event type ID is not an UUID: %v", err))
	}
	EventTypesByID[UserSearch] = UserSearchEvent

	ShowTenantUpdate, err = uuid.FromString("a2633717-12f0-4edd-bcd9-bbdf900a8ec5")
	if err != nil {
		panic(fmt.Sprintf("ShowTenantUpdate event type ID is not an UUID: %v", err))
	}
	EventTypesByID[ShowTenantUpdate] = ShowTenantUpdateEvent

	StartTenantUpdate, err = uuid.FromString("2d51de09-2ab7-4e15-9e0d-030a71756c8d")
	if err != nil {
		panic(fmt.Sprintf("StartTenantUpdate event type ID is not an UUID: %v", err))
	}
	EventTypesByID[StartTenantUpdate] = StartTenantUpdateEvent

	StopTenantUpdate, err = uuid.FromString("3dd22424-27b6-494a-a550-9611bfe41cac")
	if err != nil {
		panic(fmt.Sprintf("StopTenantUpdate event type ID is not an UUID: %v", err))
	}
	EventTypesByID[StopTenantUpdate] = StopTenantUpdateEvent

	UserDeactivationNotification, err = uuid.FromString("9f924fc3-403b-4167-b20c-a543adc4ff3c")
	if err != nil {
		panic(fmt.Sprintf("UserDeactivationNotification event type ID is not an UUID: %v", err))
	}
	EventTypesByName[UserDeactivationNotificationEvent] = UserDeactivationNotification
	EventTypesByID[UserDeactivationNotification] = UserDeactivationNotificationEvent

	UserDeactivation, err = uuid.FromString("777ede15-4b18-4720-bada-1519d1915f2e")
	if err != nil {
		panic(fmt.Sprintf("UserDeactivation event type ID is not an UUID: %v", err))
	}
	EventTypesByName[UserDeactivationEvent] = UserDeactivation
	EventTypesByID[UserDeactivation] = UserDeactivationEvent

	ListAuditLogs, err = uuid.FromString("3a7cc30b-1b7f-4764-9a35-d1bbb5cfe38a")
	if err != nil {
		panic(fmt.Sprintf("UserDeactivationNotification event type ID is not an UUID: %v", err))
	}
	EventTypesByName[ListAuditLogsEvent] = ListAuditLogs

}
