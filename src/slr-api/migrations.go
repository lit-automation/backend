//nolint
package main

import (
	"golang.org/x/net/context"
	"gopkg.in/gormigrate.v1"

	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func runMigrations(db *DBManager) {
	options := gormigrate.DefaultOptions
	options.UseTransaction = true
	ctx := context.Background()
	m := gormigrate.New(db.DB, options, []*gormigrate.Migration{
		// extract pricing information and baselocation from timeslot
		{
			ID: "2018-05-11_1",
			Migrate: func(tx *gorm.DB) error {
				articles, err := db.ArticleDB.ListForProject(ctx, uuid.FromStringOrNil("DA69A44A-F022-46E1-BA88-A1E608793F14"))
				if err != nil {
					return err
				}
				for _, article := range articles {
					if article.URL == "" {
						continue
					}
				}
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.WithError(err).Fatal("Could not migrate")
	}
	log.Info("Migration did run successfully")
}
