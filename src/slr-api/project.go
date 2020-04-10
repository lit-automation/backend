package main

import (
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/database"
	"github.com/wimspaargaren/slr-automation/src/packages/logfields"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

// ProjectController implements the project resource.
type ProjectController struct {
	*goa.Controller
}

// NewProjectController creates a project controller.
func NewProjectController(service *goa.Service) *ProjectController {
	return &ProjectController{Controller: service.NewController("ProjectController")}
}

// Create runs the create action.
func (c *ProjectController) Create(ctx *app.CreateProjectContext) error {
	// ProjectController_Create: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized("Should be logged in to retrieve projects")
	}

	project := models.Project{
		UserID: userID,
		Status: models.ProjectStatusInitPhase,
	}
	if ctx.Payload.ScrapeState != nil {
		err := project.SetScrapeStateFromString(*ctx.Payload.ScrapeState)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{logfields.UserID: userID}).Error("unable to decode scrape state")
			return ErrInternal("unable to decode scraping state")
		}
	} else {
		err = project.SetDefaultScrapeState()
		if err != nil {
			log.WithError(err).Error("unable to set scraping state")
			return ErrInternal("Unable to create project")
		}
	}

	if ctx.Payload.Name != nil {
		project.Name = *ctx.Payload.Name
	}
	if ctx.Payload.SearchString != nil {
		project.SearchString = *ctx.Payload.SearchString
	}

	err = DB.ProjectDB.Add(ctx, &project)
	if err != nil {
		log.WithError(err).Error("unable to create project")
		return ErrInternal("Unable to create project")
	}

	resp, err := c.createProjectResponse(&project)
	if err != nil {
		return err
	}
	return ctx.OK(resp)

	// ProjectController_Create: end_implement
}

// CreateFromCSV runs the createFromCSV action.
func (c *ProjectController) CreateFromCSV(ctx *app.CreateFromCSVProjectContext) error {
	// ProjectController_CreateFromCSV: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized("Should be logged in to retrieve projects")
	}

	project := models.Project{
		UserID:      userID,
		Name:        ctx.Payload.Name,
		Status:      models.ProjectStatusArticlesGathered,
		ScrapeState: []byte("{}"),
	}
	err = database.Transact(DB.DB, func(tx *gorm.DB) error {
		err = models.NewProjectDB(tx).Add(ctx, &project)
		if err != nil {
			log.WithError(err).Error("unable to create project")
			return ErrInternal("Unable to create project")
		}

		err = c.ProcessCSV(tx, project.ID, ctx.Payload.CsvContent)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	resp, err := c.createProjectResponse(&project)
	if err != nil {
		return err
	}
	return ctx.OK(resp)

	// ProjectController_CreateFromCSV: end_implement
}

// Graph runs the graph action.
func (c *ProjectController) Graph(ctx *app.GraphProjectContext) error {
	// ProjectController_Graph: start_implement

	articles, err := DB.ArticleDB.ListForProject(ctx, ctx.ProjectID)
	if err != nil {
		log.WithError(err).WithField(logfields.ProjectID, ctx.ProjectID).Error("unable to list articles")
		return ErrInternal("Unable to retrieve graph")
	}
	idMap, err := c.BuildGraphIDList(articles)
	if err != nil {
		return ErrInternal("Unable to retrieve graph")
	}

	res := app.GraphmediaCollection{}

	for _, article := range articles {
		if article.Doi == "" {
			continue
		}
		graphMedia := app.Graphmedia{
			ArticleID:   article.ID,
			URL:         article.URL,
			CitedAmount: article.CitedAmount,
			Doi:         article.Doi,
			ID:          idMap[article.Doi],
			Title:       article.Title,
		}
		citedByList, err := article.GetCitedBy()
		if err != nil {
			log.WithError(err).WithField(logfields.ArticleID, article.ID).Error("unable to get cited by list for article")
			return ErrInternal("Unable to retrieve graph")
		}
		for _, citedBy := range citedByList {
			if citedBy.Doi == "" {
				continue
			}
			graphMedia.Children = append(graphMedia.Children, &app.Articlesmallmedia{
				ID:          idMap[citedBy.Doi],
				CitedAmount: citedBy.CitedAmount,
				Doi:         citedBy.Doi,
				Title:       citedBy.Title,
				URL:         citedBy.URL,
			})
		}
		res = append(res, &graphMedia)

	}
	return ctx.OK(res)

	// ProjectController_Graph: end_implement
}

// Latest runs the latest action.
func (c *ProjectController) Latest(ctx *app.LatestProjectContext) error {
	// ProjectController_Latest: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized("Should be logged in to retrieve projects")
	}
	project, err := DB.ProjectDB.PolicyScope(userID).GetLatest(ctx)
	if err == gorm.ErrRecordNotFound {
		return ErrBadRequest("Project does not exist")
	} else if err != nil {
		log.WithError(err).WithField(logfields.UserID, userID).Error("unable to get project for user")
		return ErrInternal("Unable to retrieve projects")
	}
	resp, err := c.createProjectResponse(project)
	if err != nil {
		return err
	}
	return ctx.OK(resp)

	// ProjectController_Latest: end_implement
}

// List runs the list action.
func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	// ProjectController_List: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized("Should be logged in to retrieve projects")
	}

	projects, err := DB.ProjectDB.PolicyScope(userID).ListOrderCreatedAt(ctx)
	if err != nil {
		log.WithError(err).WithField(logfields.UserID, userID).Error("unable to list projects for user")
		return ErrInternal("Unable to list projects")
	}
	res := app.ProjectCollection{}
	for _, project := range projects {
		resp, err := c.createProjectResponse(project)
		if err != nil {
			return err
		}
		res = append(res, resp)
	}
	return ctx.OK(res)

	// ProjectController_List: end_implement
}

// Show runs the show action.
func (c *ProjectController) Show(ctx *app.ShowProjectContext) error {
	// ProjectController_Show: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized("Should be logged in to retrieve projects")
	}
	project, err := DB.ProjectDB.PolicyScope(userID).Get(ctx, ctx.ProjectID)
	if err == gorm.ErrRecordNotFound {
		return ErrBadRequest("Project does not exist")
	} else if err != nil {
		log.WithError(err).WithField(logfields.UserID, userID).Error("unable to get project for user")
		return ErrInternal("Unable to retrieve projects")
	}

	resp, err := c.createProjectResponse(project)
	if err != nil {
		return err
	}
	return ctx.OK(resp)

	// ProjectController_Show: end_implement
}

// Update runs the update action.
func (c *ProjectController) Update(ctx *app.UpdateProjectContext) error {
	// ProjectController_Update: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized("Should be logged in to retrieve projects")
	}
	project, err := DB.ProjectDB.PolicyScope(userID).Get(ctx, ctx.ProjectID)
	if err == gorm.ErrRecordNotFound {
		return ErrBadRequest("Project does not exist")
	} else if err != nil {
		log.WithError(err).WithField(logfields.UserID, userID).Error("unable to get project for user")
		return ErrInternal("Unable to update project")
	}
	if ctx.Payload.Name != nil {
		project.Name = *ctx.Payload.Name
	}

	if ctx.Payload.SearchString != nil {
		if project.Status != models.ProjectStatusInitPhase {
			return ErrBadRequest("Can't update the search query after you've started scraping")
		}
		project.SearchString = *ctx.Payload.SearchString
	}
	if ctx.Payload.Status != nil {
		newStatus := models.ProjectStatus(*ctx.Payload.Status)
		if newStatus > project.Status {
			project.Status = newStatus
		}
	}

	if ctx.Payload.ScrapeState != nil {
		err := project.SetScrapeStateFromString(*ctx.Payload.ScrapeState)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{logfields.UserID: userID, logfields.ProjectID: ctx.ProjectID}).Error("unable to decode scrape state")
			return ErrInternal("Unable to decode scraping state")
		}
	}
	err = DB.ProjectDB.PolicyScope(userID).Update(ctx, project)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{logfields.UserID: userID, logfields.ProjectID: ctx.ProjectID}).Error("unable to get project for user")
		return ErrInternal("Unable to update project")
	}

	resp, err := c.createProjectResponse(project)
	if err != nil {
		return err
	}
	return ctx.OK(resp)

	// ProjectController_Update: end_implement
}
