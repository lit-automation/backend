package main

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/database"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

func preProcess() {
	for {
		database.Transact(DB.DB, func(tx *gorm.DB) error {
			artDB := models.NewArticleDB(tx)
			articles, err := artDB.ListNotPreProcessed(context.Background())
			if err != nil {
				log.Errorf("unable to list not preprocessed: %s", err)
			} else {
				for _, a := range articles {
					log.Infof("processing art: %s", a.ID)
					err = a.SetAbstractDoc()
					if err != nil {
						log.Errorf("unable to set abstract doc for: %s, err: %s", a.ID, err)
					}
					err = a.SetFullTextDoc()
					if err != nil {
						log.Errorf("unable to set text doc for: %s, err: %s", a.ID, err)
					}
					a.Preprocessed = true
					err = artDB.Update(context.Background(), a)
					if err != nil {
						log.Errorf("unable to update art: %s, err: %s", a.ID, err)
					}
				}
			}
			return nil
		})
		time.Sleep(time.Second)
	}
}
