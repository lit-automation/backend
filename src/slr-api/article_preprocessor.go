package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func preProcess() {
	for {
		articles, err := DB.ArticleDB.ListNotPreProcessed(context.Background())
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
				err = DB.ArticleDB.Update(context.Background(), a)
				if err != nil {
					log.Errorf("unable to update art: %s, err: %s", a.ID, err)
				}
			}
		}
		time.Sleep(time.Second)
	}
}
