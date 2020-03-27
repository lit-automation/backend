package models

import (
	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// UserPolicyScope returns a gorm scope for user id
func UserPolicyScope(userID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

// ProjectPolicyScope returns a gorm scope for project id
func ProjectPolicyScope(projectID uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("project_id = ?", projectID)
	}
}
