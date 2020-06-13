package main

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/bibtex"
	"github.com/wimspaargaren/slr-automation/src/packages/crossref"
	"github.com/wimspaargaren/slr-automation/src/packages/logfields"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

var enhancementChan chan (ArticleEnhancement)

type ArticleEnhancement struct {
	ID    uuid.UUID
	Title string
}

func enhanceArticles() {
	enhancementChan = make(chan ArticleEnhancement, 500)
	go queueWorker()
	for {
		// 50 requests per second max
		if len(enhancementChan) < 50 {
			articles, err := DB.ArticleDB.ListNotChecked(context.Background())
			if err != nil {
				log.Errorf("unable to list not checked articles")
				continue
			}
			for _, article := range articles {
				article.CheckedByCrossref = true
				err = DB.ArticleDB.Update(context.Background(), article)
				if err != nil {
					log.WithError(err).WithField(logfields.ArticleID, article.ID).Error("unable to update article")
					continue
				}
				enhancementChan <- ArticleEnhancement{
					ID:    article.ID,
					Title: article.Title,
				}
			}
		}

		time.Sleep(time.Second * 2)
	}
}

type SafeCounter struct {
	sync.Mutex
	Counter int
}

func (s *SafeCounter) Get() int {
	s.Lock()
	val := s.Counter
	s.Unlock()
	return val
}

func (s *SafeCounter) Inc() {
	s.Lock()
	s.Counter++
	s.Unlock()
}

func (s *SafeCounter) Decr() {
	s.Lock()
	s.Counter--
	s.Unlock()
}

var counter = SafeCounter{}

// queueWorker processes article enhancement channel
func queueWorker() {
	crossRefClient := crossref.NewCrossRefClient(&http.Client{})
	bibClient := bibtex.NewDXDOIClient(&http.Client{})
	for {
		x := <-enhancementChan
		for counter.Get() > 40 {
			time.Sleep(time.Millisecond * 10)
		}
		counter.Inc()
		go gatherAdditionalInfo(crossRefClient, bibClient, x)

	}
}

func gatherAdditionalInfo(client crossref.Client, bibClient bibtex.Client, article ArticleEnhancement) {

	articleFromDB, err := DB.ArticleDB.Get(context.Background(), article.ID)
	if err != nil {
		log.WithError(err).WithField(logfields.ArticleID, err).Errorf("unable to retrieve article")
		return
	}

	citedByList, err := articleFromDB.GetCitedBy()
	if err != nil {
		log.WithError(err).WithField(logfields.ArticleID, article.ID).Errorf("err retrieving cited by")
	} else {
		for _, citedBy := range citedByList {
			if citedBy.Doi == "" {
				citedByRes, err := client.QueryWorks(citedBy.Title)
				if err == nil {
					if len(citedByRes.Message.Items) > 0 {
						curFound := citedByRes.Message.Items[0]
						if curFound.DOI != "" {
							citedBy.Doi = curFound.DOI
						} else {
							log.WithField(logfields.ArticleID, article.ID).Warning("no doi found for record")
						}
					}
				} else {
					log.WithError(err).WithField(logfields.ArticleID, article.ID).Errorf("err querying cited by record")
				}
			}
		}
	}
	err = articleFromDB.SetCitedBy(citedByList)
	if err != nil {
		log.WithField(logfields.ArticleID, article.ID).Warning("unable to set cited by")
	} else {
		err = DB.ArticleDB.Update(context.Background(), articleFromDB)
		if err != nil {
			log.WithError(err).WithField(logfields.ArticleID, article.ID).Error("unable to update article")
		}
	}

	var artInfo *crossref.ArticleInfo
	if articleFromDB.Doi != "" {
		res, err := client.GetOnDOI(articleFromDB.Doi)
		if err != nil {
			counter.Decr()
			log.WithError(err).WithFields(log.Fields{logfields.ArticleID: articleFromDB.ID, logfields.ArticleDOI: articleFromDB.Doi}).Errorf("err querying for work")
			return
		}
		counter.Decr()
		artInfo = &res.Message
	} else {
		res, err := client.QueryWorks(article.Title)
		if err != nil {
			counter.Decr()
			log.WithError(err).WithField(log.Fields{logfields.ArticleID: articleFromDB.ID, logfields.ArticleTitle: articleFromDB.Title}).Errorf("err querying for work")
			return
		}
		counter.Decr()

		if len(res.Message.Items) == 0 {
			log.WithField(logfields.ArticleID, article.ID).Warning("no results found")
			return
		}
		if len(res.Message.Items[0].Title) == 0 {
			log.WithField(logfields.ArticleID, article.ID).Warning("no titles found")
			return
		}
		titleFound := false
		for _, foundTitle := range res.Message.Items[0].Title {
			artTitLower := strings.ToLower(article.Title)
			foundTitLower := strings.ToLower(foundTitle)
			if strings.Contains(artTitLower, foundTitLower) ||
				strings.Contains(foundTitLower, artTitLower) {
				titleFound = true
				break
			}
			if strings.EqualFold(strings.TrimSpace(article.Title), strings.TrimSpace(foundTitle)) {
				titleFound = true
				break
			}
		}
		if !titleFound {
			log.Warningf("title not found for: %s expected: %s, got: %v", article.ID, article.Title, res.Message.Items[0].Title)
			return
		} else {
			artInfo = &res.Message.Items[0]
		}
	}
	if artInfo == nil {
		log.Error("Something went wrong artInfo not set.")
		return
	}
	updateArticleInfo(articleFromDB, artInfo)
	if articleFromDB.Doi != "" {
		updateBibTex(bibClient, articleFromDB)
	}
}

func updateBibTex(bibClient bibtex.Client, articleFromDB *models.Article) {
	bibTex, err := bibClient.GetBibTex(articleFromDB.Doi)
	if err == nil {
		articleFromDB.Bibtex = bibTex
		err = DB.ArticleDB.Update(context.Background(), articleFromDB)
		if err != nil {
			log.WithError(err).WithField(logfields.ArticleID, articleFromDB.ID).Error("unable to update article")
		}
	} else {
		log.WithError(err).WithField(logfields.ArticleID, articleFromDB.ID).Error("unable to retrieve bibtex for article")
	}
}

func updateArticleInfo(articleFromDB *models.Article, artInfo *crossref.ArticleInfo) {
	if articleFromDB.Doi == "" && artInfo.DOI != "" {
		articleFromDB.Doi = artInfo.DOI
	}
	if articleFromDB.CitedAmount == -1 && articleFromDB.CitedAmount < artInfo.IsReferencedByCount {
		articleFromDB.CitedAmount = artInfo.IsReferencedByCount
	}
	if articleFromDB.Authors == "" {
		authors := []string{}
		for _, auth := range artInfo.Author {
			authors = append(authors, auth.Given+" "+auth.Family)
		}
		articleFromDB.Authors = strings.Join(authors, " ; ")
	}
	if articleFromDB.Publisher == "" && artInfo.Publisher != "" {
		articleFromDB.Publisher = artInfo.Publisher
	}
	if articleFromDB.Abstract == "" && artInfo.Abstract != "" {
		articleFromDB.Abstract = artInfo.Abstract
	}
	if articleFromDB.Year == -1 {
		articleFromDB.Year = artInfo.Created.DateTime.Year()
	}
	if articleFromDB.URL == "" && len(artInfo.Link) > 0 {
		articleFromDB.URL = artInfo.Link[0].URL
	}
	if articleFromDB.Journal == "" && len(artInfo.ContainerTitle) > 0 {
		articleFromDB.Journal = artInfo.ContainerTitle[0]
	}

	articleFromDB.Type = artInfo.Type
	articleFromDB.Language = artInfo.Language

	err := DB.ArticleDB.Update(context.Background(), articleFromDB)
	if err != nil {
		log.WithError(err).WithField(logfields.ArticleID, articleFromDB.ID).Error("unable to update article")
	}
}
