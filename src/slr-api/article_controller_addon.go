package main

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/logfields"
)

func ProjectIDFromContext(ctx context.Context, projectID uuid.UUID) (uuid.UUID, error) {
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
