package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/packages/logfields"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

func (c *ProjectController) createProjectResponse(project *models.Project) (*app.Project, error) {
	count, err := DB.ArticleDB.CountForProject(project.ID)
	if err != nil {
		log.WithError(err).WithField(logfields.ProjectID, project.ID).Error("unable to get amount of articles for project")
		return nil, ErrInternal("Unable to update project")
	}
	appProject := project.ProjectToProject()
	appProject.AmountOfArticles = &count
	return appProject, nil
}

func (c *ProjectController) ProcessCSV(tx *gorm.DB, projectID uuid.UUID, input string) error {
	r := csv.NewReader(strings.NewReader(input))
	articleDB := models.NewArticleDB(tx)
	columnMapper := make(map[string]int)
	firstRow := true
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		if firstRow {
			columnMapper, err = c.VerifyFirstRow(record)
			if err != nil {
				return err
			}
			firstRow = false
			continue
		}
		article := models.Article{
			ProjectID:   projectID,
			CitedAmount: -1,
			Status:      models.ArticleStatusUnprocessed,
			Keywords:    []byte("[]"),
			Metadata:    []byte("{}"),
			CitedBy:     []byte("[]"),
		}
		for k, v := range columnMapper {
			switch k {
			case "title":
				article.Title = record[v]
			case "abstract":
				article.Abstract = record[v]
			case "doi":
				article.Doi = record[v]
			case "url":
				article.URL = record[v]
			case "full_text":
				article.FullText = record[v]
			case "year":
				year, err := strconv.Atoi(record[v])
				if err == nil {
					article.Year = year
				} else {
					article.Year = -1
				}
			default:
				log.Warningf("mapper got wrong key: %s", k)
			}
		}

		err = articleDB.Add(context.Background(), &article)
		if err != nil {
			log.WithError(err).Error("unable to create article")
			return ErrInternal("Unable to create article")
		}
	}
	return nil
}

func (c *ProjectController) VerifyFirstRow(elements []string) (map[string]int, error) {
	columnMapper := make(map[string]int)
	for i, e := range elements {
		switch strings.ToLower(e) {
		case "title":
			columnMapper["title"] = i
		case "year":
			columnMapper["year"] = i
		case "abstract":
			columnMapper["abstract"] = i
		case "url":
			columnMapper["url"] = i
		case "doi":
			columnMapper["doi"] = i
		case "full_text":
			columnMapper["full_text"] = i
		default:
			return nil, fmt.Errorf("incorrect column name provided: %s", e)
		}
	}
	return columnMapper, nil
}

func (c *ProjectController) BuildGraphIDList(articles []*models.Article) (map[string]int, error) {
	idMap := make(map[string]int)

	counter := 1
	for _, article := range articles {
		if article.Doi == "" {
			continue
		}
		_, ok := idMap[article.Doi]
		if !ok {
			idMap[article.Doi] = counter
			counter++
		}
		citedByList, err := article.GetCitedBy()
		if err != nil {
			log.WithError(err).WithField(logfields.ArticleID, article.ID).Error("unable to get cited by list for article")
			return nil, err
		}
		for _, citedBy := range citedByList {
			if citedBy.Doi == "" {
				continue
			}
			_, ok := idMap[citedBy.Doi]
			if !ok {
				idMap[citedBy.Doi] = counter
				counter++
			}
		}
	}
	return idMap, nil
}
