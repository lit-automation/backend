package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/wimspaargaren/slr-automation/src/packages/environment"
	"github.com/wimspaargaren/slr-automation/src/packages/goalogadapter"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

var (
	DB              *DBManager
	ErrUnauthorized = goa.NewErrorClass("unauthorized", 401)
	ErrBadRequest   = goa.NewErrorClass("badrequest", 400)
	ErrInternal     = goa.NewErrorClass("internal", 500)
)

func main() {
	// Initialize the service environment
	environment.Initialize()

	// Initialize database
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("unable to connect to the database: %s", err)
	}
	DB = NewDBManagaer(db)
	// Run cron to gather additional article info
	go enhanceArticles()

	screeningChan = make(chan uuid.UUID, 500)
	go autoScreenAbstract()

	// Create service
	service := goa.New("SLR Automation")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Set goa default logger to logrus
	logger := log.StandardLogger()
	adapter := goalogadapter.New(logger)
	service.WithLogger(adapter)

	jwtMiddleware, err := NewJWTMiddleware()
	if err != nil {
		log.Fatal(err)
	}
	app.UseJWTMiddleware(service, jwtMiddleware)
	app.UseSigninBasicAuthMiddleware(service, NewBasicAuthMiddleware())
	// Mount "article" controller
	c := NewArticleController(service)
	app.MountArticleController(service, c)
	// Mount "jwt" controller
	c2 := NewJWTController(service)
	app.MountJWTController(service, c2)
	// Mount "user" controller
	c3 := NewUserController(service)
	app.MountUserController(service, c3)
	// Mount "project" controller
	c4 := NewProjectController(service)
	app.MountProjectController(service, c4)
	// Mount "health" controller
	c5 := NewHealthController(service)
	app.MountHealthController(service, c5)
	// Mount "screening" controller
	c6 := NewScreeningController(service)
	app.MountScreeningController(service, c6)

	runMigrations(DB)

	// Start service
	if err := service.ListenAndServe(":9001"); err != nil {
		service.LogError("startup", "err", err)
	}
}
