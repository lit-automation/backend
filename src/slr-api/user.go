package main

import (
	"fmt"
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/database"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
	"golang.org/x/crypto/bcrypt"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service) *UserController {
	return &UserController{Controller: service.NewController("UserController")}
}

// Create runs the create action.
func (c *UserController) Create(ctx *app.CreateUserContext) error {
	// UserController_Create: start_implement

	passWordHash, err := bcrypt.GenerateFromPassword([]byte(ctx.Payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("unable to hash password")
		return ctx.BadRequest(fmt.Errorf("Unable to create user"))
	}
	user := models.User{
		FirstName:  ctx.Payload.FirstName,
		FamilyName: ctx.Payload.FamilyName,
		Email:      ctx.Payload.Email,
		Password:   string(passWordHash),
	}
	if ctx.Payload.MiddleName != nil {
		user.MiddleName = *ctx.Payload.MiddleName
	}
	err = database.Transact(DB.DB, func(tx *gorm.DB) error {
		err = models.NewUserDB(tx).Add(ctx, &user)
		if err != nil {
			log.WithError(err).Error("unable to add user to database")
			return ErrBadRequest("Unable to create user")
		}
		project := models.Project{
			UserID:       user.ID,
			Name:         "MyInitialProject",
			SearchString: "The search I want to perform",
		}
		err = project.SetDefaultScrapeState()
		if err != nil {
			log.WithError(err).Error("unable to create initial user project")
			return ErrBadRequest("Unable to create user")
		}
		err = models.NewProjectDB(tx).Add(ctx, &project)
		if err != nil {
			log.WithError(err).Error("unable to create initial user project")
			return ErrBadRequest("Unable to create user")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ctx.OK(user.UserToUser())

	// UserController_Create: end_implement
}
