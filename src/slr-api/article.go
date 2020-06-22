package main

import (
	"encoding/csv"
	"github.com/goadesign/goa"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/logfields"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
	"io"
	"os"
	"strconv"
)

// ArticleController implements the article resource.
type ArticleController struct {
	*goa.Controller
}

// NewArticleController creates a article controller.
func NewArticleController(service *goa.Service) *ArticleController {
	return &ArticleController{Controller: service.NewController("ArticleController")}
}

// Create runs the create action.
func (c *ArticleController) Create(ctx *app.CreateArticleContext) error {
	// ArticleController_Create: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	article, err := models.FromCreateArticlePayload(ctx.Payload)
	if err != nil {
		log.WithError(err).WithField(logfields.ProjectID, projectID).Error("unable to create article from payload")
		return ErrBadRequest("incorrect payload")
	}
	err = article.SetAbstractDoc()
	if err != nil {
		log.Errorf("unable to set abstract document: %s", err)
	}
	err = article.SetFullTextDoc()
	if err != nil {
		log.Errorf("unable to set full text document: %s", err)
	}
	article.Preprocessed = true
	article.ProjectID = projectID
	err = DB.ArticleDB.Add(ctx, article)
	if err != nil {
		log.WithError(err).WithField(logfields.ProjectID, projectID).Error("unable to add article to database")
		return ErrInternal("Unable to create article")
	}
	enhancementChan <- ArticleEnhancement{
		ID:    article.ID,
		Title: article.Title,
	}
	return ctx.OK(article.ArticleToArticle())

	// ArticleController_Create: end_implement
}

// Delete runs the delete action.
func (c *ArticleController) Delete(ctx *app.DeleteArticleContext) error {
	// ArticleController_Delete: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	err = DB.ArticleDB.PolicyScope(projectID).Delete(ctx, ctx.ArticleID)
	if err != nil {
		log.WithError(err).Error("error deleting article")
		return ErrInternal("Unable to delete article")
	}
	return ctx.OK(&app.Article{})

	// ArticleController_Delete: end_implement
}

// Download runs the download action.
func (c *ArticleController) Download(ctx *app.DownloadArticleContext) error {
	// ArticleController_Download: start_implement

	path := "csv/" + uuid.Must(uuid.NewV4()).String()

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.WithError(err).Errorf("unable to create csv file path")
		return ErrInternal("Unable to return file")
	}
	// Create empty csv file to write to
	csvFile, err := os.Create(path + "/articles.csv")
	if err != nil {
		log.WithError(err).Errorf("unable to create csv file")
		return ErrInternal("Unable to return file")
	}
	csvwriter := csv.NewWriter(csvFile)
	// List articles
	articles, err := DB.ArticleDB.ListForProject(ctx, ctx.ProjectID)
	if err != nil {
		return ErrInternal("Unable to return file")
	}
	// Create header
	err = csvwriter.Write([]string{
		"id",
		"authors",
		"abstract",
		"full_text",
		"status",
		"cited_amount",
		"doi",
		"journal",
		"platform",
		"title",
		"year",
		"search_result_number",
		"url",
		"publisher",
		"query",
		"comment",
		"bibtex",
	})
	if err != nil {
		log.WithError(err).Errorf("unable to write csv file")
		return ErrInternal("Unable to return file")
	}

	// Write article rows
	for _, article := range articles {
		err = csvwriter.Write([]string{
			article.ID.String(),
			article.Authors,
			article.Abstract,
			article.FullText,
			article.Status.String(),
			strconv.Itoa(article.CitedAmount),
			article.Doi,
			article.Journal,
			article.Platform.String(),
			article.Title,
			strconv.Itoa(article.Year),
			strconv.Itoa(article.SearchResultNumber),
			article.URL,
			article.Publisher,
			article.Query,
			article.Comment,
			article.Bibtex,
		})
		if err != nil {
			log.WithError(err).Errorf("unable to write csv file")
			return ErrInternal("Unable to return file")
		}
	}
	csvwriter.Flush()

	createdCSV, err := os.Open(path + "/articles.csv")
	if err != nil {
		log.WithError(err).Errorf("unable to open created csv file")
		return ErrInternal("Unable to return file")
	}
	defer func() {
		err = createdCSV.Close() //Close after function return
		if err != nil {
			log.WithError(err).Errorf("unable to close file")
		}
		err = os.RemoveAll(path)
		if err != nil {
			log.WithError(err).Errorf("unable to remove csv file")
		}
	}()

	// Copy csv file content to response writer
	_, err = io.Copy(ctx.ResponseWriter, createdCSV)
	if err != nil {
		log.WithError(err).Errorf("unable to copy csv file contents")
		return ErrInternal("Unable to return file")
	}

	return nil

	// ArticleController_Download: end_implement
}

// List runs the list action.
func (c *ArticleController) List(ctx *app.ListArticleContext) error {
	// ArticleController_List: start_implement

	page := 0
	requestedPageHeader := ctx.RequestData.Header.Get("X-List-Page")
	if requestedPageHeader != "" {
		parsed, err := strconv.Atoi(requestedPageHeader)
		if err == nil {
			page = parsed
		} else {
			log.WithField("page_header", requestedPageHeader).Warningf("incorrect page header")
			return ErrBadRequest("Incorrect page header provided")
		}
	}
	articles, count, err := DB.ArticleDB.ListPaginated(ctx, ctx.ProjectID, page)
	if err != nil {
		log.WithError(err).Error("error deleting article")
		return ErrInternal("Unable to list articles")

	}
	ctx.ResponseData.Header().Set("X-List-Count", strconv.Itoa(count))
	res := app.ArticleCollection{}

	for _, article := range articles {
		res = append(res, article.ArticleToArticle())
	}
	return ctx.OK(res)

	// ArticleController_List: end_implement
}

// Snowball runs the snowball action.
func (c *ArticleController) Snowball(ctx *app.SnowballArticleContext) error {
	// ArticleController_Snowball: start_implement

	art, err := DB.ArticleDB.RetrieveNotSnowBalled(ctx, ctx.ProjectID)
	if err != nil {
		log.WithError(err).Error("unable to get latest not snowballed")
		return ErrInternal("Unable to retrieve not snowballed article")
	}
	if art == nil {
		return ErrBadRequest("No articles need to be snowballed")
	}

	return ctx.OK(art.ArticleToArticle())

	// ArticleController_Snowball: end_implement
}

// Update runs the update action.
func (c *ArticleController) Update(ctx *app.UpdateArticleContext) error {
	// ArticleController_Update: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	orgArticle, err := DB.ArticleDB.Get(ctx, ctx.ArticleID)
	if err != nil {
		return ErrBadRequest("Article not found")
	}

	article, err := models.FromUpdateArticlePayload(ctx.Payload)
	if err != nil {
		log.WithError(err).WithField(logfields.ProjectID, projectID).Error("unable to create article from payload")
		return ErrBadRequest("Incorrect payload")
	}
	err = article.SetAbstractDoc()
	if err != nil {
		log.Errorf("unable to set abstract document: %s", err)
	}
	err = article.SetFullTextDoc()
	if err != nil {
		log.Errorf("unable to set full text document: %s", err)
	}
	article.ID = ctx.ArticleID
	if ctx.Payload.CitedBy != nil {
		err = article.SetCitedByFromString(*ctx.Payload.CitedBy)
		if err != nil {
			log.WithError(err).Error("unable to set cited by")
			return ErrInternal("Unable to update article")
		}
		article.BackwardSnowball = true
		citedBy, err := article.GetCitedBy()
		if err != nil {
			log.WithError(err).Error("unable to retrieve cited by")
			return ErrInternal("Unable to update article")
		}
		for _, cited := range citedBy {
			art := models.Article{
				ProjectID:   projectID,
				Title:       cited.Title,
				Doi:         cited.Doi,
				URL:         cited.URL,
				CitedAmount: cited.CitedAmount,
				Status:      models.ArticleStatusUnprocessed,
				Keywords:    []byte("[]"),
				Metadata:    []byte("{}"),
				DocAbstract: []byte("{}"),
				DocFullText: []byte("{}"),
				CitedBy:     []byte("[]"),
			}
			err = DB.ArticleDB.Add(ctx, &art)
			if err != nil {
				log.WithError(err).Error("unable to add cited by article")
				return ErrInternal("Unable to update article")
			}
			enhancementChan <- ArticleEnhancement{
				ID:    article.ID,
				Title: article.Title,
			}

		}
	}

	err = DB.ArticleDB.PolicyScope(projectID).Update(ctx, article)
	if err != nil {
		log.WithError(err).Error("unable to update article")
		return ErrInternal("Unable to update article")
	}

	if article.Doi != "" || ctx.Payload.CitedBy != nil {
		if article.Doi != orgArticle.Doi {
			err := DB.DB.Model(&article).Updates(map[string]interface{}{"checked_by_crossref": false}).Error
			if err != nil {
				log.Errorf("Error setting crossref to false")
			}
		}
	}

	return ctx.OK(article.ArticleToArticle())

	// ArticleController_Update: end_implement
}
