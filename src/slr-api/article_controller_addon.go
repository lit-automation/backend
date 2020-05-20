package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/logfields"
)

func (c *ArticleController) ProjectIDFromContext(ctx context.Context, projectID uuid.UUID) (uuid.UUID, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, ErrUnauthorized("Should be logged in to retrieve projects")
	}
	project, err := DB.ProjectDB.PolicyScope(userID).Get(ctx, projectID)
	if err == gorm.ErrRecordNotFound {
		return uuid.Nil, ErrBadRequest("Project does not exist")
	} else if err != nil {
		log.WithError(err).WithField(logfields.UserID, userID).Error("unable to get project for user")
		return uuid.Nil, ErrInternal("Unable to retrieve projects")
	}
	return project.ID, nil
}

func (c *ArticleController) TryToFindAbstract(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if strings.Contains(url, "ieeexplore.ieee.org") {
		return c.processIEEE(resp)
	} else {
		return "", fmt.Errorf("unknown url")
	}
}

func (c *ArticleController) processIEEE(response *http.Response) (string, error) {
	// doc, err := goquery.NewDocumentFromReader(response.Body)
	// if err != nil {
	// 	log.Errorf("Could not create document from read, error: %s", err.Error())
	// 	return "", err
	// }
	return "", nil
}
