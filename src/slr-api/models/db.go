package models

import (
	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/slr-automation/src/packages/database"
)

func InitDB() (*gorm.DB, error) {
	m := []interface{}{
		Article{},
		Project{},
		User{},
	}
	return database.ConnectDefault("slr", m)
}
