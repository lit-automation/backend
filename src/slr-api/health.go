package main

import (
	"github.com/goadesign/goa"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

// HealthController implements the health resource.
type HealthController struct {
	*goa.Controller
}

// NewHealthController creates a health controller.
func NewHealthController(service *goa.Service) *HealthController {
	return &HealthController{Controller: service.NewController("HealthController")}
}

// Health runs the health action.
func (c *HealthController) Health(ctx *app.HealthHealthContext) error {
	// HealthController_Health: start_implement
	res := &app.Health{
		Health: true,
	}
	return ctx.OK(res)
	// HealthController_Health: end_implement
}
