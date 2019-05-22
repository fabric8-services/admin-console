package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.Resource("audit_log", func() {

	a.BasePath("/auditlogs")

	a.Action("create", func() {
		a.Security("jwt")
		a.Routing(
			a.POST("users/:username"),
		)
		a.Description("Add an auditlog for a user")
		a.Params(func() {
			a.Param("username", d.String)
			a.Required("username")
		})
		a.Payload(createAuditLog)
		a.Response(d.NoContent)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})

var createAuditLog = a.MediaType("application/vnd.createauditlog+json", func() {
	a.UseTrait("jsonapi-media-type")
	a.TypeName("CreateAuditLog")
	a.Description("Create an auditlog")
	a.Attributes(func() {
		a.Attribute("data", createAuditLogData)
		a.Required("data")
	})
	a.View("default", func() {
		a.Attribute("data")
		a.Required("data")
	})
})

// createAuditLogData represents the data of an audit log to create for a user
var createAuditLogData = a.Type("CreateAuditLogData", func() {
	a.Attribute("type", d.String, "type of the audit log", func() {
		a.Enum("audit_logs")
	})
	a.Attribute("attributes", createAuditLogDataAttributes, "Attributes of the audit log. ")
	a.Attribute("links", genericLinks)
	a.Required("type", "attributes")
})

var createAuditLogDataAttributes = a.Type("CreateAuditLogDataAttributes", func() {
	a.Attribute("event_type", d.String, "the type of event")
	a.Attribute("event_params", a.HashOf(d.String, d.Any), "a generic map holding the params of the event to log")
	a.Required("event_type", "event_params")
})
