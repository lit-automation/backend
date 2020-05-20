package main

import (
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

// Show runs the show action.
func (c *ScreeningController) Show(ctx *app.ShowScreeningContext) error {
	// ScreeningController_Show: start_implement

	// projectID, err := c.ProjectIDFromContext(ctx, ctx.ProjectID)
	// if err != nil {
	// 	return err
	// }
	// TODO add security for article & project
	article, err := DB.ArticleDB.Get(ctx, ctx.ArticleID)
	if err != nil {
		return ctx.BadRequest(fmt.Errorf("Article not found"))
	}
	if article.Title == "" || article.Abstract == "" {
		return ErrBadRequest("Both title and abstract needs to be present before screening of articles")
	}
	res, err := GetScreeningMediaForProject(ctx.ProjectID, article.Title, article.Abstract)
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

	article, err := DB.ArticleDB.Get(ctx, ctx.ArticleID)
	if err != nil {
		return ctx.BadRequest(fmt.Errorf("Article not found"))
	}
	err = TrainModel(ctx.ProjectID, article.Title, article.Abstract, ctx.Payload.Include)
	if err != nil {
		return ctx.InternalServerError()
	}
	article.Status = models.ArticleStatusNotUseful
	if ctx.Payload.Include {
		article.Status = models.ArticleStatusUseful
	}
	err = DB.ArticleDB.Update(ctx, article)
	if err != nil {
		return ctx.InternalServerError()
	}

	res := &app.Articlescreening{}
	return ctx.OK(res)

	// ScreeningController_Update: end_implement
}
