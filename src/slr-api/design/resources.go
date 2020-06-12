package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var articlePayload = func() {
	Attribute("year", Integer)
	Attribute("cited_amount", Integer)
	Attribute("search_result_number", Integer)
	Attribute("abstract", String)
	Attribute("language", String)
	Attribute("query", String)
	Attribute("query_platform", String)
	Attribute("authors", String)
	Attribute("journal", String)
	Attribute("publisher", String)
	Attribute("title", String)
	Attribute("got_pdf", Boolean)
	Attribute("doi", String)
	Attribute("url", String)
	Attribute("type", String)
	Attribute("platform", Integer)
	Attribute("status", Integer)
	Attribute("comment", String)
	Attribute("cited_by", String)
}

var projectPayload = func() {
	Attribute("name", String)
	Attribute("search_string", String)
	Attribute("scrape_state", String)
	Attribute("status", Integer)
}

var projectCSVPayload = func() {
	Attribute("name", String)
	Attribute("csv_content", String)
	Required("name", "csv_content")
}

var _ = Resource("health", func() {
	DefaultMedia(HealthMedia)
	BasePath("/health")
	Action("health", func() {
		Routing(
			GET(""),
		)
		Description("Checks API health")
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("article", func() {
	Security(JWT, func() {
		Scope("api:access")
	})
	Parent("project")
	DefaultMedia(ArticleMedia)
	BasePath("/article")
	Action("create", func() {
		Routing(
			POST(""),
		)
		Payload(articlePayload)
		Response(OK)
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
	Action("snowball", func() {
		Routing(
			GET("/snowball"),
		)
		Params(func() {
			Param("projectID", UUID, "Project ID")
		})
		Description("Show article which has not been snowballed")
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
	Action("update", func() {
		Routing(
			PUT("/:articleID"),
		)
		Params(func() {
			Param("articleID", UUID, "Article ID")
		})
		Payload(articlePayload)
		Response(OK)
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
	Action("list", func() {
		Routing(
			GET("list"),
		)
		Params(func() {
			Param("status", Integer)
			Param("title", String)
			Param("doi", String)
			Param("abstract", String)
			Param("type", String)
			Param("year", Integer)
			Param("amount_cited", Integer)
		})
		Description("List articles")
		Response(OK, func() {
			Media(CollectionOf(ArticleMedia))
		})
		Response(BadRequest, ErrorMedia)
	})
	Action("delete", func() {
		Routing(
			DELETE("/:articleID"),
		)
		Description("List articles")
		Params(func() {
			Param("articleID", UUID, "Article ID")
		})
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
	Action("download", func() {
		Routing(
			GET("download"),
		)
		Response(OK)
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("project", func() {
	Security(JWT, func() {
		Scope("api:access")
	})
	DefaultMedia(ProjectMedia)
	BasePath("/project")
	Action("show", func() {
		Routing(
			GET("/:projectID"),
		)
		Params(func() {
			Param("projectID", UUID, "Project ID")
		})
		Description("Show project")
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
	Action("graph", func() {
		Routing(
			GET("/:projectID/graph"),
		)
		Params(func() {
			Param("projectID", UUID, "Project ID")
		})
		Description("Give graph presentation of project")
		Response(OK, func() {
			Media(CollectionOf(GraphMedia))
		})
		Response(BadRequest, ErrorMedia)
	})
	Action("latest", func() {
		Routing(
			GET("/latest"),
		)
		Description("Show the latest project")
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Create new project")
		Payload(projectPayload)
		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
	Action("createFromCSV", func() {
		Routing(
			POST("/csv"),
		)
		Description("Create new project")
		Payload(projectCSVPayload)
		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
	Action("list", func() {
		Routing(
			GET("/list"),
		)
		Description("List projects")
		Response(OK, func() {
			Media(CollectionOf(ProjectMedia))
		})
		Response(BadRequest, ErrorMedia)
	})
	Action("update", func() {
		Routing(
			PUT("/:projectID"),
		)
		Params(func() {
			Param("projectID", UUID, "Project ID")
		})
		Payload(projectPayload)
		Response(OK)
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
	Action("removeDuplicates", func() {
		Routing(
			POST("/:projectID/removeduplicates"),
		)
		Params(func() {
			Param("projectID", UUID, "Project ID")
		})
		Response(OK, func() {
			Media(DuplicateMedia)
		})
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("user", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})
	DefaultMedia(UserMedia)
	BasePath("/user")
	Action("create", func() {
		NoSecurity()
		Routing(
			POST(""),
		)
		Description("Create new user")
		Payload(func() {
			Member("first_name", String)
			Member("middle_name", String)
			Member("family_name", String)
			Member("email")
			Member("password")

			Required("first_name", "family_name", "email", "password")
		})
		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
		Response(Conflict, ErrorMedia)
	})
})

var _ = Resource("screening", func() {
	Parent("project")
	BasePath("/screen")
	DefaultMedia(ArticleScreeningMedia)
	Action("show", func() {
		Routing(
			GET("/:articleID"),
		)
		Params(func() {
			Param("articleID", UUID, "Article ID")
		})
		Description("Show article screening")
		Response(OK)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
	Action("update", func() {
		Routing(
			PUT("/:articleID"),
		)
		Params(func() {
			Param("articleID", UUID, "Article ID")
		})
		Payload(func() {
			Member("include", Boolean)
			Required("include")
		})
		Response(OK)
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
	Action("auto", func() {
		Routing(
			POST("auto"),
		)
		Response(OK, func() {
			Media(AutoScreenAbstract)
		})
		Response(NotFound)
		Response(InternalServerError)
		Response(BadRequest, ErrorMedia)
	})
})
