package main

import (
	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

type DBManager struct {
	DB        *gorm.DB
	ArticleDB *models.ArticleDB
	ProjectDB *models.ProjectDB
	UserDB    *models.UserDB
}

func NewDBManagaer(db *gorm.DB) *DBManager {
	return &DBManager{
		DB:        db,
		ArticleDB: models.NewArticleDB(db),
		ProjectDB: models.NewProjectDB(db),
		UserDB:    models.NewUserDB(db),
	}
}
