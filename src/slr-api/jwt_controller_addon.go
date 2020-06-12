package main

import (
	"context"
	"fmt"
	"time"

	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/goadesign/goa"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/slr-automation/src/packages/jwt"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"golang.org/x/crypto/bcrypt"
)

func (c *JWTController) createJWTAndAddToHeader(responseData *goa.ResponseData, userID uuid.UUID) error {
	in1Month := time.Now().Add(time.Hour * 24 * 7 * time.Duration(4)).Unix()
	token, err := jwt.CreateJWTToken(in1Month, userID, []string{"api:access"})
	if err != nil {
		log.WithError(err).Error("unable to create JWT")
		return ErrInternal("Unable to create JWT token")
	}
	// Set auth header for client retrieval
	responseData.Header().Set("Authorization", "Bearer "+token)
	return nil
}

/*
	Middleware functions
*/

// NewJWTMiddleware creates a middleware that checks for the presence of a JWT Authorization header
// and validates its content. A real app would probably use goa's JWT security middleware instead.
func NewJWTMiddleware() (goa.Middleware, error) {
	err := jwt.LoadKeys()
	if err != nil {
		panic(err)
	}
	pub := jwt.PublicKey()

	forceFail := func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return h(ctx, rw, req)
		}
	}
	fm, err := goa.NewMiddleware(forceFail)
	if err != nil {
		panic(err)
	}
	return goajwt.New(pub, fm,
		app.NewJWTSecurity()), nil
}

// NewBasicAuthMiddleware creates a middleware that checks for the presence of a basic auth header
// and validates its content.
func NewBasicAuthMiddleware() goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			// Retrieve and log basic auth info
			email, pass, ok := req.BasicAuth()
			if !ok || email == "" || pass == "" {
				log.Warningf("Failed basic auth, user: %s", email)
				return ErrUnauthorized("Failed basic authentication")
			}

			// Normal basic auth
			user, err := DB.UserDB.GetOnEmail(ctx, email)
			if err != nil || err == gorm.ErrRecordNotFound {
				log.Warningf("User trying to sign in not found in DB, error: %s", err.Error())
				return ErrUnauthorized("Authorization failed")
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
			if err != nil {
				return ErrUnauthorized("Wrong username or password")
			}

			// Proceed
			rw.Header().Set("Access-Control-Expose-Headers", "Authorization")
			return h(ctx, rw, req)
		}
	}
}

func userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	token := goajwt.ContextJWT(ctx)
	if token == nil {
		return uuid.Nil, fmt.Errorf("forcing failure because token is missing")
	}
	claims := token.Claims.(jwtgo.MapClaims)

	userIDString := claims["sub"].(string)
	if userIDString == "" {
		return uuid.Nil, fmt.Errorf("no user id specified in token")
	}
	return uuid.FromString(userIDString)
}
