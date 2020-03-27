package main

import (
	"github.com/goadesign/goa"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

// JWTController implements the jwt resource.
type JWTController struct {
	*goa.Controller
}

// NewJWTController creates a jwt controller.
func NewJWTController(service *goa.Service) *JWTController {
	return &JWTController{Controller: service.NewController("JWTController")}
}

// Refresh runs the refresh action.
func (c *JWTController) Refresh(ctx *app.RefreshJWTContext) error {
	// JWTController_Refresh: start_implement

	userID, err := userIDFromContext(ctx)
	if err != nil {
		log.WithError(err).Error("unable to retrieve user id from context")
		return ErrBadRequest("Incorrect bearrer specfied")
	}

	return c.createJWTAndAddToHeader(ctx.ResponseData, userID)

	// JWTController_Refresh: end_implement
}

// Signin runs the signin action.
func (c *JWTController) Signin(ctx *app.SigninJWTContext) error {
	// JWTController_Signin: start_implement

	email, _, _ := ctx.RequestData.Request.BasicAuth()

	user, err := DB.UserDB.GetOnEmail(ctx, email)
	if err != nil {
		return ErrUnauthorized("User and or password incorrect")
	}

	return c.createJWTAndAddToHeader(ctx.ResponseData, user.ID)

	// JWTController_Signin: end_implement
}
