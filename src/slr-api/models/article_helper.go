// Code generated by goagen v1.4.1, DO NOT EDIT.
//
// API "SLR Automation": Model Helpers
//
// Command:
// $ goagen
// --design=github.com/wimspaargaren/slr-automation/src/slr-api/design
// --out=$(GOPATH)/src/github.com/wimspaargaren/slr-automation/src/slr-api
// --version=v1.4.3

package models

import (
	"context"
	"time"

	"github.com/goadesign/goa"
	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

// MediaType Retrieval Functions

// ListArticle returns an array of view: default.
func (m *ArticleDB) ListArticle(ctx context.Context) []*app.Article {
	defer goa.MeasureSince([]string{"goa", "db", "article", "listarticle"}, time.Now())

	var native []*Article
	var objs []*app.Article
	err := m.Db.Scopes().Table(m.TableName()).Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing Article", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.ArticleToArticle())
	}

	return objs
}

// ArticleToArticle loads a Article and builds the default view of media type Article.
func (m *Article) ArticleToArticle() *app.Article {
	article := &app.Article{}
	article.Abstract = &m.Abstract
	article.Authors = &m.Authors
	article.CitedAmount = &m.CitedAmount
	article.Comment = &m.Comment
	article.Doi = &m.Doi
	article.FullText = &m.FullText
	article.GotPdf = &m.GotPdf
	article.ID = m.ID
	article.Journal = &m.Journal
	article.Language = &m.Language
	tmp1 := int(m.Platform)
	article.Platform = &tmp1
	article.ProjectID = &m.ProjectID
	article.Publisher = &m.Publisher
	article.Query = &m.Query
	article.QueryPlatform = &m.QueryPlatform
	article.SearchResultNumber = &m.SearchResultNumber
	tmp2 := int(m.Status)
	article.Status = &tmp2
	article.Title = &m.Title
	article.Type = &m.Type
	article.URL = &m.URL
	article.Year = &m.Year

	return article
}

// OneArticle loads a Article and builds the default view of media type Article.
func (m *ArticleDB) OneArticle(ctx context.Context, id uuid.UUID) (*app.Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "onearticle"}, time.Now())

	var native Article
	err := m.Db.Scopes().Table(m.TableName()).Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting Article", "error", err.Error())
		return nil, err
	}

	view := *native.ArticleToArticle()
	return &view, err
}

// CRUD Functions

// Get returns a single Article as a Database Model
// This is more for use internally, and probably not what you want in your controllers
func (m *ArticleDB) Get(ctx context.Context, id uuid.UUID) (*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "get"}, time.Now())

	var native Article
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of Article
func (m *ArticleDB) List(ctx context.Context) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "list"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *ArticleDB) Add(ctx context.Context, model *Article) error {
	defer goa.MeasureSince([]string{"goa", "db", "article", "add"}, time.Now())

	err := m.Db.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding Article", "error", err.Error())
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *ArticleDB) Update(ctx context.Context, model *Article) error {
	defer goa.MeasureSince([]string{"goa", "db", "article", "update"}, time.Now())

	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating Article", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *ArticleDB) Delete(ctx context.Context, id uuid.UUID) error {
	defer goa.MeasureSince([]string{"goa", "db", "article", "delete"}, time.Now())

	err := m.Db.Where("id = ?", id).Delete(&Article{}).Error
	if err != nil {
		goa.LogError(ctx, "error deleting Article", "error", err.Error())
		return err
	}

	return nil
}
