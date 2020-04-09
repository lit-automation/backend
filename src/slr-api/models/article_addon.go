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

	uuid "github.com/gofrs/uuid"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

// UpdateArticlePayload is the article update action payload.
type ArticlePayload struct {
	Abstract           *string `form:"abstract,omitempty" json:"abstract,omitempty" yaml:"abstract,omitempty" xml:"abstract,omitempty"`
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
	return ArticleFromBytes(b)
}

func FromUpdateArticlePayload(p *app.UpdateArticlePayload) (*Article, error) {
	if p == nil {
		return nil, fmt.Errorf("payload is empty")
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal create article payload: %w", err)
	}
	return ArticleFromBytes(b)
}

func ArticleFromBytes(b []byte) (*Article, error) {
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
	if articlePayload.Authors != nil {
		article.Authors = *articlePayload.Authors
	}
	if articlePayload.CitedAmount != nil {
		article.CitedAmount = *articlePayload.CitedAmount
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
	}
	if articlePayload.Title != nil {
		article.Title = *articlePayload.Title
	}
	if articlePayload.URL != nil {
		article.URL = *articlePayload.URL
	}
	if articlePayload.Year != nil {
		article.Year = *articlePayload.Year
	}

	article.Keywords = []byte("{}")
	article.Metadata = []byte("{}")
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
	// TODO add filtering
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

// List returns an array of Article
func (m *ArticleDB) ListForProject(ctx context.Context, projectID uuid.UUID) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "list"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("project_id = ? AND id = 'B215BE2B-E860-4900-8E3B-A872B93C3687'", projectID).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	fmt.Println("LEN OBJS", len(objs))

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

	err := m.Db.Table(m.TableName()).Where("(backward_snowball = false OR backward_snowball IS NULL) AND project_id = ? AND doi != ''", projectID).Limit(1).Find(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &obj, nil
}
