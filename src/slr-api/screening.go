package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/goadesign/goa"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

// ScreeningController implements the screening resource.
type ScreeningController struct {
	*goa.Controller
}

// NewScreeningController creates a screening controller.
func NewScreeningController(service *goa.Service) *ScreeningController {
	return &ScreeningController{Controller: service.NewController("ScreeningController")}
}

// Auto runs the auto action.
func (c *ScreeningController) Auto(ctx *app.AutoScreeningContext) error {
	// ScreeningController_Auto: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	err = VerifyModel(projectID)
	if err != nil {
		return ErrBadRequest(err)
	}

	articles, err := DB.ArticleDB.ListOnStatus(context.Background(), projectID, models.ArticleStatusUnprocessed)
	if err != nil {
		log.Errorf("unable to list articles for screening: %s", err)
		return ctx.InternalServerError()
	}
	for _, article := range articles {
		if !article.Preprocessed {
			continue
		}
		res, err := GetScreeningMediaForProject(projectID, article, true)
		if err != nil {
			log.Errorf("Err predicting: %s", err)
		}
		article.Status = models.ArticleStatusIncludedOnAbstract
		if res.Tfidf.Class == "Exclude" {
			article.Status = models.ArticleStatusExcluded
		}
		err = DB.ArticleDB.Update(context.Background(), article)
		if err != nil {
			log.Errorf("Unable to update article status: %s", err)
		}
	}

	res := &app.Autoscreenabstract{}
	return ctx.OK(res)

	// ScreeningController_Auto: end_implement
}

// Show runs the show action.
func (c *ScreeningController) Show(ctx *app.ShowScreeningContext) error {
	// ScreeningController_Show: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	article, err := DB.ArticleDB.Get(ctx, ctx.ArticleID)
	if err != nil {
		return ErrBadRequest(fmt.Errorf("Article not found"))
	}
	if !article.Preprocessed {
		return ErrBadRequest(fmt.Errorf("Article is not yet pre processed for screening, if you just imported your articles give the backend some time to process"))
	}
	if article.ProjectID != projectID {
		return ctx.BadRequest(fmt.Errorf("Incorrect article ID"))
	}
	res, err := GetScreeningMediaForProject(projectID, article, true)
	if err != nil {
		log.WithError(err).Error("unable to retrieve screening media")
		return ctx.InternalServerError()
	}
	res.ID = article.ID
	return ctx.OK(res)

	// ScreeningController_Show: end_implement
}

// Shownext runs the shownext action.
func (c *ScreeningController) Shownext(ctx *app.ShownextScreeningContext) error {
	// ScreeningController_Shownext: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	articles, err := DB.ArticleDB.ListOnStatus(context.Background(), projectID, models.ArticleStatusUnprocessed)
	if err != nil {
		log.Errorf("unable to list articles for screening: %s", err)
		return ctx.InternalServerError()
	}
	if len(articles) == 0 {
		return ctx.BadRequest(fmt.Errorf("No articles found for screening"))
	}
	screeningResult := []*app.Articlescreening{}
	for _, article := range articles {
		if !article.Preprocessed {
			continue
		}
		res, err := GetScreeningMediaForProject(projectID, article, true)
		if err != nil {
			log.Errorf("Err predicting: %s", err)
		}

		res.ID = article.ID
		screeningResult = append(screeningResult, res)
	}

	sort.Slice(screeningResult, func(i, j int) bool {
		return screeningResult[i].Tfidf.Confidence < screeningResult[j].Tfidf.Confidence
	})
	counterInc := 0
	counterEx := 0
	for _, temp := range screeningResult {
		if temp.Tfidf.Class == "Include" {
			counterInc++
		} else {
			counterEx++
		}
	}
	totalArticles, err := DB.ArticleDB.CountOnStatusList(ctx, projectID, []models.ArticleStatus{models.ArticleStatusExcluded, models.ArticleStatusIncluded, models.ArticleStatusIncludedOnAbstract, models.ArticleStatusUnprocessed, models.ArticleStatusUnknown})
	if err != nil {
		log.Errorf("unable to retrieve total amount of articles: %s", err)
	}
	screenedArticles, err := DB.ArticleDB.CountOnStatusList(ctx, projectID, []models.ArticleStatus{models.ArticleStatusExcluded, models.ArticleStatusIncluded, models.ArticleStatusIncludedOnAbstract})
	if err != nil {
		log.Errorf("unable to retrieve amount of articles screened: %s", err)
	}
	modelDetails := &app.Modeldetails{
		AutoExclude:      float64(counterEx),
		AutoInclude:      float64(counterInc),
		ScreenedArticles: float64(screenedArticles),
		TotalArticles:    float64(totalArticles),
	}
	fmt.Println("Current situation:", counterInc, counterEx)
	for _, result := range screeningResult {
		if result.Tfidf.Abstract != "" {
			result.ModelDetails = modelDetails
			return ctx.OK(result)
		}
	}
	screeningResult[0].ModelDetails = modelDetails
	return ctx.OK(screeningResult[0])

	// ScreeningController_Shownext: end_implement
}

// Update runs the update action.
func (c *ScreeningController) Update(ctx *app.UpdateScreeningContext) error {
	// ScreeningController_Update: start_implement

	projectID, err := ProjectIDFromContext(ctx, ctx.ProjectID)
	if err != nil {
		return err
	}
	article, err := DB.ArticleDB.Get(ctx, ctx.ArticleID)
	if err != nil {
		return ErrBadRequest(fmt.Errorf("Article not found"))
	}
	if !article.Preprocessed {
		return ErrBadRequest(fmt.Errorf("Article is not yet pre processed for screening, if you just imported your articles give the backend some time to process"))
	}
	if article.ProjectID != projectID {
		return ctx.NotFound()
	}
	if article.Status != models.ArticleStatusUnprocessed {
		return ErrBadRequest(fmt.Errorf("This article is not unprocessed and therefore can't be used to train the model"))
	}
	err = TrainModel(ctx.ProjectID, article, true, ctx.Payload.Include)
	if err != nil {
		return ctx.InternalServerError()
	}
	article.Status = models.ArticleStatusExcluded
	if ctx.Payload.Include {
		article.Status = models.ArticleStatusIncludedOnAbstract
	}
	fmt.Println("Updating article: ", article.Status)
	err = DB.ArticleDB.Update(ctx, article)
	if err != nil {
		return ctx.InternalServerError()
	}

	res := &app.Articlescreening{}
	return ctx.OK(res)

	// ScreeningController_Update: end_implement
}
