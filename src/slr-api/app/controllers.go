// Code generated by goagen v1.4.1, DO NOT EDIT.
//
// API "SLR Automation": Application Controllers
//
// Command:
// $ goagen
// --design=github.com/wimspaargaren/slr-automation/src/slr-api/design
// --out=$(GOPATH)/src/github.com/wimspaargaren/slr-automation/src/slr-api
// --version=v1.4.3

package app

import (
	"context"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/cors"
	"github.com/goadesign/goa/encoding/form"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(form.NewDecoder, "application/x-www-form-urlencoded")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// ProjectController is the controller interface for the Project actions.
type ProjectController interface {
	goa.Muxer
	Create(*CreateProjectContext) error
	CreateFromCSV(*CreateFromCSVProjectContext) error
	Graph(*GraphProjectContext) error
	Latest(*LatestProjectContext) error
	List(*ListProjectContext) error
	RemoveDuplicates(*RemoveDuplicatesProjectContext) error
	Show(*ShowProjectContext) error
	Update(*UpdateProjectContext) error
}

// MountProjectController "mounts" a Project resource controller on the given service.
func MountProjectController(service *goa.Service, ctrl ProjectController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/v1/project", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/csv", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/graph", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/latest", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/list", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/removeduplicates", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID", ctrl.MuxHandler("preflight", handleProjectOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateProjectPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("POST", "/v1/project", ctrl.MuxHandler("create", h, unmarshalCreateProjectPayload))
	service.LogInfo("mount", "ctrl", "Project", "action", "Create", "route", "POST /v1/project", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateFromCSVProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateFromCSVProjectPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.CreateFromCSV(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("POST", "/v1/project/csv", ctrl.MuxHandler("createFromCSV", h, unmarshalCreateFromCSVProjectPayload))
	service.LogInfo("mount", "ctrl", "Project", "action", "CreateFromCSV", "route", "POST /v1/project/csv", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewGraphProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Graph(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("GET", "/v1/project/:projectID/graph", ctrl.MuxHandler("graph", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "Graph", "route", "GET /v1/project/:projectID/graph", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewLatestProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Latest(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("GET", "/v1/project/latest", ctrl.MuxHandler("latest", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "Latest", "route", "GET /v1/project/latest", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("GET", "/v1/project/list", ctrl.MuxHandler("list", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "List", "route", "GET /v1/project/list", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewRemoveDuplicatesProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.RemoveDuplicates(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("POST", "/v1/project/:projectID/removeduplicates", ctrl.MuxHandler("removeDuplicates", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "RemoveDuplicates", "route", "POST /v1/project/:projectID/removeduplicates", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("GET", "/v1/project/:projectID", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "Show", "route", "GET /v1/project/:projectID", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateProjectPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleProjectOrigin(h)
	service.Mux.Handle("PUT", "/v1/project/:projectID", ctrl.MuxHandler("update", h, unmarshalUpdateProjectPayload))
	service.LogInfo("mount", "ctrl", "Project", "action", "Update", "route", "PUT /v1/project/:projectID", "security", "jwt")
}

// handleProjectOrigin applies the CORS response headers corresponding to the origin.
func handleProjectOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "X-List-Count")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Pin, X-Platform, content-type, X-Vendor-id, x-list-limit, x-list-page, x-list-filter, X-List-Count, Access-Control-Allow-Origin, accept")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// unmarshalCreateProjectPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateProjectPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createProjectPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalCreateFromCSVProjectPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateFromCSVProjectPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createFromCSVProjectPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalUpdateProjectPayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateProjectPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateProjectPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// ArticleController is the controller interface for the Article actions.
type ArticleController interface {
	goa.Muxer
	Create(*CreateArticleContext) error
	Delete(*DeleteArticleContext) error
	Download(*DownloadArticleContext) error
	List(*ListArticleContext) error
	Snowball(*SnowballArticleContext) error
	Update(*UpdateArticleContext) error
}

// MountArticleController "mounts" a Article resource controller on the given service.
func MountArticleController(service *goa.Service, ctrl ArticleController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/article", ctrl.MuxHandler("preflight", handleArticleOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/article/:articleID", ctrl.MuxHandler("preflight", handleArticleOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/article/download", ctrl.MuxHandler("preflight", handleArticleOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/article/list", ctrl.MuxHandler("preflight", handleArticleOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/article/snowball", ctrl.MuxHandler("preflight", handleArticleOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateArticleContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateArticlePayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleArticleOrigin(h)
	service.Mux.Handle("POST", "/v1/project/:projectID/article", ctrl.MuxHandler("create", h, unmarshalCreateArticlePayload))
	service.LogInfo("mount", "ctrl", "Article", "action", "Create", "route", "POST /v1/project/:projectID/article", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteArticleContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleArticleOrigin(h)
	service.Mux.Handle("DELETE", "/v1/project/:projectID/article/:articleID", ctrl.MuxHandler("delete", h, nil))
	service.LogInfo("mount", "ctrl", "Article", "action", "Delete", "route", "DELETE /v1/project/:projectID/article/:articleID", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDownloadArticleContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Download(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleArticleOrigin(h)
	service.Mux.Handle("GET", "/v1/project/:projectID/article/download", ctrl.MuxHandler("download", h, nil))
	service.LogInfo("mount", "ctrl", "Article", "action", "Download", "route", "GET /v1/project/:projectID/article/download", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListArticleContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleArticleOrigin(h)
	service.Mux.Handle("GET", "/v1/project/:projectID/article/list", ctrl.MuxHandler("list", h, nil))
	service.LogInfo("mount", "ctrl", "Article", "action", "List", "route", "GET /v1/project/:projectID/article/list", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSnowballArticleContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Snowball(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleArticleOrigin(h)
	service.Mux.Handle("GET", "/v1/project/:projectID/article/snowball", ctrl.MuxHandler("snowball", h, nil))
	service.LogInfo("mount", "ctrl", "Article", "action", "Snowball", "route", "GET /v1/project/:projectID/article/snowball", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateArticleContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateArticlePayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleArticleOrigin(h)
	service.Mux.Handle("PUT", "/v1/project/:projectID/article/:articleID", ctrl.MuxHandler("update", h, unmarshalUpdateArticlePayload))
	service.LogInfo("mount", "ctrl", "Article", "action", "Update", "route", "PUT /v1/project/:projectID/article/:articleID", "security", "jwt")
}

// handleArticleOrigin applies the CORS response headers corresponding to the origin.
func handleArticleOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "X-List-Count")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Pin, X-Platform, content-type, X-Vendor-id, x-list-limit, x-list-page, x-list-filter, X-List-Count, Access-Control-Allow-Origin, accept")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// unmarshalCreateArticlePayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateArticlePayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createArticlePayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalUpdateArticlePayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateArticlePayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateArticlePayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// HealthController is the controller interface for the Health actions.
type HealthController interface {
	goa.Muxer
	Health(*HealthHealthContext) error
}

// MountHealthController "mounts" a Health resource controller on the given service.
func MountHealthController(service *goa.Service, ctrl HealthController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/v1/health", ctrl.MuxHandler("preflight", handleHealthOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewHealthHealthContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Health(rctx)
	}
	h = handleHealthOrigin(h)
	service.Mux.Handle("GET", "/v1/health", ctrl.MuxHandler("health", h, nil))
	service.LogInfo("mount", "ctrl", "Health", "action", "Health", "route", "GET /v1/health")
}

// handleHealthOrigin applies the CORS response headers corresponding to the origin.
func handleHealthOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "X-List-Count")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Pin, X-Platform, content-type, X-Vendor-id, x-list-limit, x-list-page, x-list-filter, X-List-Count, Access-Control-Allow-Origin, accept")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// JWTController is the controller interface for the JWT actions.
type JWTController interface {
	goa.Muxer
	Refresh(*RefreshJWTContext) error
	Signin(*SigninJWTContext) error
}

// MountJWTController "mounts" a JWT resource controller on the given service.
func MountJWTController(service *goa.Service, ctrl JWTController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/v1/jwt/refresh", ctrl.MuxHandler("preflight", handleJWTOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/jwt/signin", ctrl.MuxHandler("preflight", handleJWTOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewRefreshJWTContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Refresh(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleJWTOrigin(h)
	service.Mux.Handle("GET", "/v1/jwt/refresh", ctrl.MuxHandler("refresh", h, nil))
	service.LogInfo("mount", "ctrl", "JWT", "action", "Refresh", "route", "GET /v1/jwt/refresh", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSigninJWTContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Signin(rctx)
	}
	h = handleSecurity("SigninBasicAuth", h)
	h = handleJWTOrigin(h)
	service.Mux.Handle("POST", "/v1/jwt/signin", ctrl.MuxHandler("signin", h, nil))
	service.LogInfo("mount", "ctrl", "JWT", "action", "Signin", "route", "POST /v1/jwt/signin", "security", "SigninBasicAuth")
}

// handleJWTOrigin applies the CORS response headers corresponding to the origin.
func handleJWTOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "X-List-Count")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Pin, X-Platform, content-type, X-Vendor-id, x-list-limit, x-list-page, x-list-filter, X-List-Count, Access-Control-Allow-Origin, accept")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// ScreeningController is the controller interface for the Screening actions.
type ScreeningController interface {
	goa.Muxer
	Auto(*AutoScreeningContext) error
	Show(*ShowScreeningContext) error
	Update(*UpdateScreeningContext) error
}

// MountScreeningController "mounts" a Screening resource controller on the given service.
func MountScreeningController(service *goa.Service, ctrl ScreeningController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/screen/auto", ctrl.MuxHandler("preflight", handleScreeningOrigin(cors.HandlePreflight()), nil))
	service.Mux.Handle("OPTIONS", "/v1/project/:projectID/screen/:articleID", ctrl.MuxHandler("preflight", handleScreeningOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewAutoScreeningContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Auto(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleScreeningOrigin(h)
	service.Mux.Handle("POST", "/v1/project/:projectID/screen/auto", ctrl.MuxHandler("auto", h, nil))
	service.LogInfo("mount", "ctrl", "Screening", "action", "Auto", "route", "POST /v1/project/:projectID/screen/auto", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowScreeningContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleScreeningOrigin(h)
	service.Mux.Handle("GET", "/v1/project/:projectID/screen/:articleID", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Screening", "action", "Show", "route", "GET /v1/project/:projectID/screen/:articleID", "security", "jwt")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateScreeningContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateScreeningPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	h = handleSecurity("jwt", h, "api:access")
	h = handleScreeningOrigin(h)
	service.Mux.Handle("PUT", "/v1/project/:projectID/screen/:articleID", ctrl.MuxHandler("update", h, unmarshalUpdateScreeningPayload))
	service.LogInfo("mount", "ctrl", "Screening", "action", "Update", "route", "PUT /v1/project/:projectID/screen/:articleID", "security", "jwt")
}

// handleScreeningOrigin applies the CORS response headers corresponding to the origin.
func handleScreeningOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "X-List-Count")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Pin, X-Platform, content-type, X-Vendor-id, x-list-limit, x-list-page, x-list-filter, X-List-Count, Access-Control-Allow-Origin, accept")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// unmarshalUpdateScreeningPayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateScreeningPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateScreeningPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// UserController is the controller interface for the User actions.
type UserController interface {
	goa.Muxer
	Create(*CreateUserContext) error
}

// MountUserController "mounts" a User resource controller on the given service.
func MountUserController(service *goa.Service, ctrl UserController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/v1/user", ctrl.MuxHandler("preflight", handleUserOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateUserContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateUserPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	h = handleUserOrigin(h)
	service.Mux.Handle("POST", "/v1/user", ctrl.MuxHandler("create", h, unmarshalCreateUserPayload))
	service.LogInfo("mount", "ctrl", "User", "action", "Create", "route", "POST /v1/user")
}

// handleUserOrigin applies the CORS response headers corresponding to the origin.
func handleUserOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Expose-Headers", "X-List-Count")
			rw.Header().Set("Access-Control-Max-Age", "600")
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Pin, X-Platform, content-type, X-Vendor-id, x-list-limit, x-list-page, x-list-filter, X-List-Count, Access-Control-Allow-Origin, accept")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// unmarshalCreateUserPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateUserPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createUserPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}
