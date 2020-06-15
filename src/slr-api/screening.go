package main

import (
	"context"
	"fmt"
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

	err := VerifyModel(ctx.ProjectID)
	if err != nil {
		return ErrBadRequest(err)
	}

	screeningChan <- ctx.ProjectID
	articles, err := DB.ArticleDB.ListOnStatus(context.Background(), ctx.ProjectID, models.ArticleStatusUnprocessed)
	if err != nil {
		log.Errorf("Err listing: %s", err)
	}
	result := fmt.Sprintf("Screening will take approximatly %d seconds", len(articles)/3)

	res := &app.Autoscreenabstract{
		Message: result,
	}
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
	if article.ProjectID != projectID {
		return ctx.BadRequest(fmt.Errorf("Incorrect article ID"))
	}
	if article.Title == "" || article.Abstract == "" {
		return ErrBadRequest("Both title and abstract needs to be present before screening of articles")
	}
	res, err := GetScreeningMediaForProject(projectID, article.Title, article.Abstract)
	if err != nil {
		log.WithError(err).Error("unable to retrieve screening media")
		return ctx.InternalServerError()
	}
	res.ID = article.ID
	return ctx.OK(res)

	// ScreeningController_Show: end_implement
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
	if article.ProjectID != projectID {
		return ctx.NotFound()
	}
	if article.Status != models.ArticleStatusUnprocessed {
		return ErrBadRequest(fmt.Errorf("This article is not unprocessed and therefore can't be used to train the model"))
	}
	err = TrainModel(ctx.ProjectID, article.Title, article.Abstract, ctx.Payload.Include)
	if err != nil {
		return ctx.InternalServerError()
	}
	article.Status = models.ArticleStatusExcluded
	if ctx.Payload.Include {
		article.Status = models.ArticleStatusIncludedOnAbstract
	}
	err = DB.ArticleDB.Update(ctx, article)
	if err != nil {
		return ctx.InternalServerError()
	}

	res := &app.Articlescreening{}
	return ctx.OK(res)

	// ScreeningController_Update: end_implement
}
