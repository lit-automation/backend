package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

//HealthMedia media type for healthroute
var HealthMedia = MediaType("application/vnd.health+json", func() {
	Description("Media type used to indicate if services are healthy")
	Attributes(func() {
		Attribute("health", Boolean, "True if API is healthy")
		Required("health")
	})
	View("default", func() {
		Attribute("health")
	})
})

// UserMedia is the mediatype for a user in the system
var UserMedia = MediaType("application/vnd.user+json", func() {
	Description("User mediatype")
	Attributes(func() {
		Attribute("id", UUID, "ID of the user in the DB")
		Attribute("first_name", String, "First name of the user")
		Attribute("family_name", String, "Family name of the user")
		Attribute("middle_name", String, "Middle name of the user")
		Attribute("email", String, "Email of the user")

		Required("id", "first_name", "family_name", "email")
	})
	View("default", func() {
		Attribute("id")
		Attribute("email")
		Attribute("family_name")
		Attribute("middle_name")
		Attribute("first_name")
	})
})

//ArticleMedia article media
var ArticleMedia = MediaType("application/vnd.article+json", func() {
	Description("Article mediatype")
	Attributes(func() {
		Attribute("id", UUID)
		Attribute("project_id", UUID)
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
		Attribute("platform", Integer)
		Attribute("metadata", Any)
		Attribute("keywords", Any)
		Attribute("status", Integer)
		Attribute("comment", String)
		Attribute("type", String)
		Required("id")
	})
	View("default", func() {
		Attribute("id")
		Attribute("project_id")
		Attribute("year")
		Attribute("cited_amount")
		Attribute("search_result_number")
		Attribute("abstract")
		Attribute("language")
		Attribute("query")
		Attribute("query_platform")
		Attribute("authors")
		Attribute("journal")
		Attribute("publisher")
		Attribute("title")
		Attribute("got_pdf")
		Attribute("doi")
		Attribute("url")
		Attribute("platform")
		Attribute("status")
		Attribute("comment")
		Attribute("type")
	})
})

//ArticleMetadataMedia media type for ArticleMetadataMedia
var ArticleMetadataMedia = MediaType("application/vnd.articlemetadata+json", func() {
	Attributes(func() {
		Attribute("rq_ids", ArrayOf(Integer), "List of rq ids")
		Required("rq_ids")
	})
	View("default", func() {
		Attribute("rq_ids")
	})
})

//ProjectMedia project media
var ProjectMedia = MediaType("application/vnd.project+json", func() {
	Description("Project mediatype")
	Attributes(func() {
		Attribute("id", UUID)
		Attribute("name", String)
		Attribute("status", Integer)
		Attribute("search_string", String)
		Attribute("scrape_state", Any)
		Attribute("amount_of_articles", Integer)
		Required("id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("status")
		Attribute("search_string")
		Attribute("scrape_state")
		Attribute("amount_of_articles", Integer)
	})
})

var GraphMedia = MediaType("application/vnd.graphmedia+json", func() {
	Attributes(func() {
		Attribute("id", Integer, "graph id")
		Attribute("article_id", UUID)
		Attribute("title", String)
		Attribute("cited_amount", Integer)
		Attribute("doi", String)
		Attribute("url", String)
		Attribute("children", ArrayOf(ArticleSmall))
		Required("id", "article_id", "title", "cited_amount", "doi", "url", "children")
	})
	View("default", func() {
		Attribute("id")
		Attribute("article_id")
		Attribute("title")
		Attribute("cited_amount")
		Attribute("doi")
		Attribute("url")
		Attribute("children")
	})
})

var ArticleSmall = MediaType("application/vnd.articlesmallmedia+json", func() {
	Attributes(func() {
		Attribute("id", Integer, "graph id")
		Attribute("title", String)
		Attribute("cited_amount", Integer)
		Attribute("doi", String)
		Attribute("url", String)
		Required("id", "title", "cited_amount", "doi", "url")
	})
	View("default", func() {
		Attribute("id")
		Attribute("title")
		Attribute("cited_amount")
		Attribute("doi")
		Attribute("url")
	})
})

//ArticleScreeningMedia article media
var ArticleScreeningMedia = MediaType("application/vnd.articlescreening+json", func() {
	Description("Article Screening mediatype")
	Attributes(func() {
		Attribute("id", UUID)
		Attribute("title", TextPredictMedia)
		Attribute("abstract", TextPredictMedia)
		Attribute("sentences", ArrayOf(TextPredictMedia))
		Required("id", "title", "abstract", "sentences")
	})
	View("default", func() {
		Attribute("id")
		Attribute("title")
		Attribute("abstract")
		Attribute("sentences")
	})
})

var TextPredictMedia = MediaType("application/vnd.textpredictmedia+json", func() {
	Attribute("text", String)
	Attribute("class", String)
	Attribute("confidence", Number)
	Required("text", "class", "confidence")
	View("default", func() {
		Attribute("text")
		Attribute("class")
		Attribute("confidence")
	})
})
