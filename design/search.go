package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.Resource("search", func() {
	a.BasePath("/search")

	a.Action("search_users", func() {
		a.Security("jwt")
		a.Routing(
			a.GET("users"),
		)
		a.Description("Search by fullname")
		a.Headers(func() {
			a.Header("Authorization", d.String, "the authorization header")
		})
		a.Params(func() {
			a.Param("q", d.String)
			a.Param("page[offset]", d.String, "Paging start position") // #428
			a.Param("page[limit]", d.Integer, "Paging size")
			a.Required("q")
		})
		a.Response(d.OK) // here we don't need to specify a media type, because we're just proxying to `auth`
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})
