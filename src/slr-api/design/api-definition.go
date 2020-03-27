package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("SLR Automation", func() {
	Title("SLR Automation")
	Description("API handler for SLR Automation")
	Contact(func() {
		Name("Wim")
		Email("wim_spaargaren@live.nl")
		URL("")
	})
	Host("localhost:80")
	Scheme("http")
	BasePath("/v1")
	Origin("*", func() {
		Headers("Authorization, X-Pin",
			"X-Platform",
			"content-type",
			"X-Vendor-id",
			"x-list-limit",
			"x-list-page",
			"x-list-filter",
			"X-List-Count",
			"Access-Control-Allow-Origin",
			"accept")
		Expose("X-List-Count")
		Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")
		MaxAge(600)
		Credentials()
	})
	Consumes("application/json")
	Consumes("application/x-www-form-urlencoded", func() {
		Package("github.com/goadesign/goa/encoding/form")
	})
})
