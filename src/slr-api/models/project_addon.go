package models

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/goadesign/goa"
	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type ScrapeStatus string

const (
	ScrapeStatusActive = "active"
	ScrapeStatusPaused = "paused"
)

type ScrapeState struct {
	Status          ScrapeStatus `json:"status"`
	Platforms       []Platform   `json:"platforms"`
	URLs            []string     `json:"urls"`
	ProcessingIndex int64        `json:"processing_index"`
	PlatformCounter int64        `json:"platform_counter"`
	NextURL         string       `json:"next_url"`
}

func (p *Project) SetScrapeStateFromString(state string) error {
	res, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return err
	}
	p.ScrapeState = res
	return nil
}

func (p *Project) SetDefaultScrapeState() error {
	return p.SetScrapeState(ScrapeState{
		Status: ScrapeStatusPaused,
		Platforms: []Platform{PlatformGoogleScholar,
			PlatformACM,
			PlatformSpringer,
			PlatformIEEE,
			PlatformWebOfScience,
			PlatformScienceDirect},
		ProcessingIndex: 0,
		PlatformCounter: 0,
		NextURL:         "",
	})
}

func (p *Project) SetScrapeState(state ScrapeState) error {
	b, err := json.Marshal(state)
	if err != nil {
		return err
	}
	p.ScrapeState = b
	return nil
}

func (p *Project) GetScrapeState() (*ScrapeState, error) {
	var res ScrapeState
	err := json.Unmarshal(p.ScrapeState, &res)
	return &res, err
}

// PolicyScope returns a project db which scope for given user id
func (m *ProjectDB) PolicyScope(userID uuid.UUID) *ProjectDB {
	return &ProjectDB{Db: m.Db.Scopes(UserPolicyScope(userID))}
}

func (m *ProjectDB) GetLatest(ctx context.Context) (*Project, error) {
	defer goa.MeasureSince([]string{"goa", "db", "project", "get"}, time.Now())

	var native Project
	err := m.Db.Table(m.TableName()).Order("created_at desc").Limit(1).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

func (m *ProjectDB) ListOrderCreatedAt(ctx context.Context) ([]*Project, error) {
	defer goa.MeasureSince([]string{"goa", "db", "project", "list"}, time.Now())

	var objs []*Project
	err := m.Db.Table(m.TableName()).Order("created_at desc").Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}
