package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.Resource("tenants_update", func() {
	a.BasePath("/tenants/update")

	a.Action("show", func() {
		a.Security("jwt")
		a.Routing(
			a.GET(""),
		)
		a.Headers(func() {
			a.Header("Authorization", d.String, "the authorization header")
		})
		a.Description("Get information about last/ongoing update.")
		a.Response(d.OK) // here we don't specify a media type, because we're just proxying to `tenant`
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("start", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Headers(func() {
			a.Header("Authorization", d.String, "the authorization header")
		})
		a.Params(func() {
			a.Param("cluster_url", d.String, "the URL of the OSO cluster the update should be limited to")
			a.Param("env_type", d.String, "environment type the update should be executed for", func() {
				a.Enum("user", "che", "jenkins", "stage", "run")
			})
		})

		a.Description("Start new cluster-wide update.")
		a.Response(d.Accepted)
		a.Response(d.Conflict) // here we don't specify a media type, because we're just proxying to `tenant`
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("stop", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE(""),
		)
		a.Headers(func() {
			a.Header("Authorization", d.String, "the authorization header")
		})
		a.Description("Stops an ongoing cluster-wide update.")
		a.Response(d.Accepted)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
})
