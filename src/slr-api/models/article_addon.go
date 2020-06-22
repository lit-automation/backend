package models

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/jdkato/prose.v2"

	uuid "github.com/gofrs/uuid"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

type ScreeningData struct {
	Sentences []prose.Sentence
	Tokens    []prose.Token
}

// UpdateArticlePayload is the article update action payload.
type ArticlePayload struct {
	Abstract           *string `form:"abstract,omitempty" json:"abstract,omitempty" yaml:"abstract,omitempty" xml:"abstract,omitempty"`
	FullText           *string `form:"full_text,omitempty" json:"full_text,omitempty" yaml:"full_text,omitempty" xml:"full_text,omitempty"`
	Authors            *string `form:"authors,omitempty" json:"authors,omitempty" yaml:"authors,omitempty" xml:"authors,omitempty"`
	CitedAmount        *int    `form:"cited_amount,omitempty" json:"cited_amount,omitempty" yaml:"cited_amount,omitempty" xml:"cited_amount,omitempty"`
	Comment            *string `form:"comment,omitempty" json:"comment,omitempty" yaml:"comment,omitempty" xml:"comment,omitempty"`
	Doi                *string `form:"doi,omitempty" json:"doi,omitempty" yaml:"doi,omitempty" xml:"doi,omitempty"`
	GotPdf             *bool   `form:"got_pdf,omitempty" json:"got_pdf,omitempty" yaml:"got_pdf,omitempty" xml:"got_pdf,omitempty"`
	Journal            *string `form:"journal,omitempty" json:"journal,omitempty" yaml:"journal,omitempty" xml:"journal,omitempty"`
	Language           *string `form:"language,omitempty" json:"language,omitempty" yaml:"language,omitempty" xml:"language,omitempty"`
	Platform           *int    `form:"platform,omitempty" json:"platform,omitempty" yaml:"platform,omitempty" xml:"platform,omitempty"`
	Status             *int    `form:"status,omitempty" json:"status,omitempty" yaml:"status,omitempty" xml:"status,omitempty"`
	Publisher          *string `form:"publisher,omitempty" json:"publisher,omitempty" yaml:"publisher,omitempty" xml:"publisher,omitempty"`
	Query              *string `form:"query,omitempty" json:"query,omitempty" yaml:"query,omitempty" xml:"query,omitempty"`
	QueryPlatform      *string `form:"query_platform,omitempty" json:"query_platform,omitempty" yaml:"query_platform,omitempty" xml:"query_platform,omitempty"`
	SearchResultNumber *int    `form:"search_result_number,omitempty" json:"search_result_number,omitempty" yaml:"search_result_number,omitempty" xml:"search_result_number,omitempty"`
	Title              *string `form:"title,omitempty" json:"title,omitempty" yaml:"title,omitempty" xml:"title,omitempty"`
	URL                *string `form:"url,omitempty" json:"url,omitempty" yaml:"url,omitempty" xml:"url,omitempty"`
	Year               *int    `form:"year,omitempty" json:"year,omitempty" yaml:"year,omitempty" xml:"year,omitempty"`
}

func FromCreateArticlePayload(p *app.CreateArticlePayload) (*Article, error) {
	if p == nil {
		return nil, fmt.Errorf("payload is empty")
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal create article payload: %w", err)
	}
	return ArticleFromBytes(b, false)
}

func FromUpdateArticlePayload(p *app.UpdateArticlePayload) (*Article, error) {
	if p == nil {
		return nil, fmt.Errorf("payload is empty")
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal create article payload: %w", err)
	}
	return ArticleFromBytes(b, true)
}

func ArticleFromBytes(b []byte, update bool) (*Article, error) {
	log.Infof(string(b))
	var articlePayload *ArticlePayload
	err := json.Unmarshal(b, &articlePayload)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal create article payload into article: %w", err)
	}
	article := &Article{}
	if articlePayload.Abstract != nil {
		article.Abstract = *articlePayload.Abstract
	}
	if articlePayload.FullText != nil {
		article.FullText = *articlePayload.FullText
	}
	if articlePayload.Authors != nil {
		article.Authors = *articlePayload.Authors
	}
	if articlePayload.CitedAmount != nil {
		article.CitedAmount = *articlePayload.CitedAmount
	} else if !update {
		article.CitedAmount = -1
	}
	if articlePayload.Comment != nil {
		article.Comment = *articlePayload.Comment
	}
	if articlePayload.Doi != nil {
		article.Doi = *articlePayload.Doi
	}
	if articlePayload.GotPdf != nil {
		article.GotPdf = *articlePayload.GotPdf
	}
	if articlePayload.Journal != nil {
		article.Journal = *articlePayload.Journal
	}
	if articlePayload.Language != nil {
		article.Language = *articlePayload.Language
	}
	if articlePayload.Platform != nil {
		article.Platform = Platform(*articlePayload.Platform)
	}
	if articlePayload.Publisher != nil {
		article.Publisher = *articlePayload.Publisher
	}
	if articlePayload.Query != nil {
		article.Query = *articlePayload.Query
	}
	if articlePayload.QueryPlatform != nil {
		article.QueryPlatform = *articlePayload.QueryPlatform
	}
	if articlePayload.SearchResultNumber != nil {
		article.SearchResultNumber = *articlePayload.SearchResultNumber
	}
	if articlePayload.Status != nil {
		article.Status = ArticleStatus(*articlePayload.Status)
	} else if !update {
		article.Status = ArticleStatusUnprocessed
	}
	if articlePayload.Title != nil {
		article.Title = *articlePayload.Title
	}
	if articlePayload.URL != nil {
		article.URL = *articlePayload.URL
	}
	if articlePayload.Year != nil {
		article.Year = *articlePayload.Year
	} else if !update {
		article.Year = -1
	}

	article.Keywords = []byte("{}")
	article.Metadata = []byte("{}")
	article.DocAbstract = []byte("{}")
	article.DocFullText = []byte("{}")
	article.CitedBy = []byte("[]")
	return article, nil
}

// PolicyScope returns a project db which scope for given user id
func (m *ArticleDB) PolicyScope(projectID uuid.UUID) *ArticleDB {
	return &ArticleDB{Db: m.Db.Scopes(ProjectPolicyScope(projectID))}
}

func (m *ArticleDB) CountForProject(projectID uuid.UUID) (int, error) {
	var count int
	err := m.Db.Table(m.TableName()).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// Count number of total articles
type Count struct {
	Count int
}

// List returns an array of Article
func (m *ArticleDB) ListPaginated(ctx *app.ListArticleContext, projectID uuid.UUID, page int) ([]*Article, int, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "list"}, time.Now())

	query := m.Db.Table(m.TableName()).
		Where("project_id = ?", projectID)
		// Type        *string
		// Year        *intj
	if ctx.Abstract != nil {
		query = query.Where("LOWER(abstract) LIKE '%" + strings.ToLower(*ctx.Abstract) + "%'")
	}
	if ctx.AmountCited != nil {
		query = query.Where("cited_amount > ?", *ctx.AmountCited)
	}
	if ctx.Doi != nil {
		query = query.Where("LOWER(doi) LIKE '%" + strings.ToLower(*ctx.Doi) + "%'")
	}
	if ctx.Status != nil {
		query = query.Where("status = ?", *ctx.Status)
	}
	if ctx.Type != nil {
		query = query.Where("type = ?", *ctx.Type)
	}
	if ctx.Year != nil {
		query = query.Where("year > ?", *ctx.Year)
	}
	if ctx.Title != nil {
		query = query.Where("LOWER(title) LIKE '%" + strings.ToLower(*ctx.Title) + "%'")
	}
	var count Count
	err := query.
		Select("COUNT(id)").
		Find(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	var objs []*Article
	limit := 20
	err = query.
		Limit(limit).
		Offset(page * limit).
		Order("doi DESC").
		Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return objs, count.Count, nil
}

// SetDuplicates adjust status of article id list to duplicates
func (m *ArticleDB) SetDuplicates(ctx context.Context, articleIDs []uuid.UUID) error {
	defer goa.MeasureSince([]string{"goa", "db", "article", "setduplicates"}, time.Now())

	return m.Db.Exec("UPDATE articles SET status = ? WHERE articles.id IN (?)", ArticleStatusDuplicate, articleIDs).Error
}

// List returns an array of Article
func (m *ArticleDB) ListForProject(ctx context.Context, projectID uuid.UUID) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "list"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("project_id = ?", projectID).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

func (m *ArticleDB) CountOnStatusList(ctx context.Context, projectID uuid.UUID, statuses []ArticleStatus) (int, error) {
	var count Count
	err := m.Db.Table(m.TableName()).Where("project_id = ? AND status IN (?)", projectID, statuses).
		Select("COUNT(id)").
		Find(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count.Count, nil
}

// ListNonDuplicatesForProject returns an array of Article
func (m *ArticleDB) ListNonDuplicatesForProject(ctx context.Context, projectID uuid.UUID) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "list"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("project_id = ? AND status != ?", projectID, ArticleStatusDuplicate).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// ListNotPreProcessed returns an array of Articles which are not yet preprocessed
func (m *ArticleDB) ListNotPreProcessed(ctx context.Context) ([]*Article, error) {

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("preprocessed = false OR preprocessed IS NULL").Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// ListOnStatus returns an array of Article
func (m *ArticleDB) ListOnStatus(ctx context.Context, projectID uuid.UUID, status ArticleStatus) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "list"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("project_id = ? AND status = ?", projectID, status).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

func (m *ArticleDB) ListNotChecked(ctx context.Context) ([]*Article, error) {
	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("checked_by_crossref = false OR checked_by_crossref IS NULL").Limit(50).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return objs, nil
}

type CitedBy struct {
	// ID          uuid.UUID `json:"id"`
	CitedAmount int    `form:"cited_amount,omitempty" json:"cited_amount,omitempty" yaml:"cited_amount,omitempty" xml:"cited_amount,omitempty"`
	Doi         string `form:"doi,omitempty" json:"doi,omitempty" yaml:"doi,omitempty" xml:"doi,omitempty"`
	Title       string `form:"title,omitempty" json:"title,omitempty" yaml:"title,omitempty" xml:"title,omitempty"`
	URL         string `form:"url,omitempty" json:"url,omitempty" yaml:"url,omitempty" xml:"url,omitempty"`
}

func (a *Article) SetCitedByFromString(citedBy string) error {
	res, err := base64.StdEncoding.DecodeString(citedBy)
	if err != nil {
		return err
	}
	a.CitedBy = res
	return nil
}

func (a *Article) GetCitedBy() ([]*CitedBy, error) {
	var res []*CitedBy
	err := json.Unmarshal(a.CitedBy, &res)
	return res, err
}

func (a *Article) SetCitedBy(c []*CitedBy) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	a.CitedBy = b
	return nil
}

func (m *ArticleDB) RetrieveNotSnowBalled(ctx context.Context, projectID uuid.UUID) (*Article, error) {
	var obj Article

	err := m.Db.Table(m.TableName()).Where("(backward_snowball = false OR backward_snowball IS NULL) AND (status = ? OR status = ?) AND project_id = ? AND doi != ''", ArticleStatusIncludedOnAbstract, ArticleStatusIncluded, projectID).Order("cited_amount DESC").Limit(1).Find(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &obj, nil
}

func (a *Article) SetAbstractDoc() error {
	doc, err := prose.NewDocument(a.Title + " " + a.Abstract) //, prose.WithExtraction(false), prose.WithTokenization(false))
	if err != nil {
		return err
	}
	screeningData := ScreeningData{
		Tokens:    doc.Tokens(),
		Sentences: doc.Sentences(),
	}
	b, err := json.Marshal(screeningData)
	if err != nil {
		return err
	}
	a.DocAbstract = b
	return nil
}

func (a *Article) GetAbstractDoc() (*ScreeningData, error) {
	var doc ScreeningData
	err := json.Unmarshal(a.DocAbstract, &doc)
	return &doc, err
}

func (a *Article) SetFullTextDoc() error {
	doc, err := prose.NewDocument(a.Title+" "+a.Abstract, prose.WithExtraction(false), prose.WithTokenization(false))
	if err != nil {
		return err
	}
	screeningData := ScreeningData{
		Tokens:    doc.Tokens(),
		Sentences: doc.Sentences(),
	}
	b, err := json.Marshal(screeningData)
	if err != nil {
		return err
	}
	a.DocFullText = b
	return nil
}

func (a *Article) GetFullTextDoc() (*ScreeningData, error) {
	var doc ScreeningData
	err := json.Unmarshal(a.DocFullText, &doc)
	return &doc, err
}
