package design

import (
	"lab.weave.nl/forks/gorma"
	. "lab.weave.nl/forks/gorma/dsl"
)

var _ = StorageGroup("slrautomation", func() {
	Store("postgres", gorma.Postgres, func() {

		Enum("Platform", func() {
			Value("GoogleScholar", 1)
			Value("ACM", 2)
			Value("Springer", 3)
			Value("IEEE", 4)
			Value("WebOfScience", 5)
			Value("ScienceDirect", 6)
		})

		Enum("Status", func() {
			Value("Unprocessed", 1)
			Value("NotUseful", 2)
			Value("Useful", 3)
			Value("Unknown", 4)
			Value("Duplicate", 5)
		})

		Enum("ScopeType", func() {
			Value("access", 1)
		})

		Enum("ProjectStatus", func() {
			Value("InitPhase", 1)
			Value("ConductingSearch", 2)
			Value("ArticlesGathered", 3)
			Value("ScreeningPhase", 4)
			Value("DataExtractionPhase", 5)
			Value("Done", 6)
		})

		Model("Article", func() {
			Description("Article represent an article found in a database")
			RendersTo(ArticleMedia)
			Field("id", gorma.UUID, func() {
				PrimaryKey()
			})
			Field("year", gorma.Integer)
			Field("cited_amount", gorma.Integer)
			Field("search_result_number", gorma.Integer)
			Field("abstract", gorma.String)
			Field("language", gorma.String)
			Field("type", gorma.String)
			Field("query", gorma.String)
			Field("query_platform", gorma.String)
			Field("authors", gorma.String)
			Field("journal", gorma.String)
			Field("publisher", gorma.String)
			Field("title", gorma.String)
			Field("checked_by_crossref", gorma.Boolean)
			Field("backward_snowball", gorma.Boolean)
			Field("got_pdf", gorma.Boolean)
			Field("doi", gorma.String)
			Field("url", gorma.String)
			Field("bibtex", gorma.String)
			Field("platform", gorma.Enum, "Platform")
			Field("metadata", gorma.JSON)
			Field("cited_by", gorma.JSON)
			Field("keywords", gorma.JSON)
			Field("status", gorma.Enum, "Status")
			Field("comment", gorma.String)
		})

		Model("User", func() {
			Description("User model represents a user in the platform")
			RendersTo(UserMedia)

			Field("id", gorma.UUID, func() {
				PrimaryKey()
			})
			Field("first_name", gorma.String)
			Field("middle_name", gorma.String)
			Field("family_name", gorma.String)
			Field("email", gorma.String, func() {
				SQLTag("unique")
			})
			Field("password", gorma.String)
			Field("reset_token", gorma.String)
			Field("scope", gorma.Enum, "ScopeType")
			HasMany("Projects", "Project")

		})

		Model("Project", func() {
			Description("Project model represents an SLR project")
			RendersTo(ProjectMedia)

			Field("id", gorma.UUID, func() {
				PrimaryKey()
			})
			Field("name", gorma.String)
			Field("status", gorma.Enum, "ProjectStatus")
			Field("search_string", gorma.String)
			Field("scrape_state", gorma.JSON)

			HasMany("Articles", "Article")
		})

	})
})
