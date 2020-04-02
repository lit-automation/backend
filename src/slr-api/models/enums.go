// Code generated by goagen v1.4.1, DO NOT EDIT.
//
// API "SLR Automation": Enums
//
// Command:
// $ goagen
// --design=github.com/wimspaargaren/slr-automation/src/slr-api/design
// --out=$(GOPATH)/src/github.com/wimspaargaren/slr-automation/src/slr-api
// --version=v1.4.1

package models

import (
	"database/sql/driver"
	"fmt"
)

// ArticleStatus Enum
type ArticleStatus int64

const (
	ArticleStatusUnprocessed ArticleStatus = 1
	ArticleStatusNotUseful   ArticleStatus = 2
	ArticleStatusUseful      ArticleStatus = 3
	ArticleStatusUnknown     ArticleStatus = 4
	ArticleStatusDuplicate   ArticleStatus = 5
)

var articleStatusStrings = map[int64]string{
	1: "Unprocessed",
	2: "NotUseful",
	3: "Useful",
	4: "Unknown",
	5: "Duplicate",
}

var articleStatusStringMap = map[string]ArticleStatus{
	"Unprocessed": ArticleStatusUnprocessed,
	"NotUseful":   ArticleStatusNotUseful,
	"Useful":      ArticleStatusUseful,
	"Unknown":     ArticleStatusUnknown,
	"Duplicate":   ArticleStatusDuplicate,
}

func (u *ArticleStatus) Scan(value interface{}) error {
	i, ok := value.(int64)
	if !ok {
		return fmt.Errorf("could not cast value %v to int64", value)
	}
	*u = ArticleStatus(i)
	return nil
}
func (u *ArticleStatus) ScanFromString(name string) error {
	var ok bool
	*u, ok = articleStatusStringMap[name]
	if !ok {
		return fmt.Errorf("%s is not a valid name for a value of ArticleStatus", name)
	}
	return nil
}
func (u ArticleStatus) Value() (driver.Value, error) { return int64(u), nil }
func (u ArticleStatus) String() string {
	if u == 0 {
		return "undefined"
	}
	return articleStatusStrings[int64(u)]
}
func (u ArticleStatus) AllStrings() map[int64]string { return articleStatusStrings }

// Platform Enum
type Platform int64

const (
	PlatformGoogleScholar Platform = 1
	PlatformACM           Platform = 2
	PlatformSpringer      Platform = 3
	PlatformIEEE          Platform = 4
	PlatformWebOfScience  Platform = 5
	PlatformScienceDirect Platform = 6
)

var platformStrings = map[int64]string{
	1: "GoogleScholar",
	2: "ACM",
	3: "Springer",
	4: "IEEE",
	5: "WebOfScience",
	6: "ScienceDirect",
}

var platformStringMap = map[string]Platform{
	"GoogleScholar": PlatformGoogleScholar,
	"ACM":           PlatformACM,
	"Springer":      PlatformSpringer,
	"IEEE":          PlatformIEEE,
	"WebOfScience":  PlatformWebOfScience,
	"ScienceDirect": PlatformScienceDirect,
}

func (u *Platform) Scan(value interface{}) error {
	i, ok := value.(int64)
	if !ok {
		return fmt.Errorf("could not cast value %v to int64", value)
	}
	*u = Platform(i)
	return nil
}
func (u *Platform) ScanFromString(name string) error {
	var ok bool
	*u, ok = platformStringMap[name]
	if !ok {
		return fmt.Errorf("%s is not a valid name for a value of Platform", name)
	}
	return nil
}
func (u Platform) Value() (driver.Value, error) { return int64(u), nil }
func (u Platform) String() string {
	if u == 0 {
		return "undefined"
	}
	return platformStrings[int64(u)]
}
func (u Platform) AllStrings() map[int64]string { return platformStrings }

// ProjectStatus Enum
type ProjectStatus int64

const (
	ProjectStatusInitPhase           ProjectStatus = 1
	ProjectStatusConductingSearch    ProjectStatus = 2
	ProjectStatusArticlesGathered    ProjectStatus = 3
	ProjectStatusScreeningPhase      ProjectStatus = 4
	ProjectStatusDataExtractionPhase ProjectStatus = 5
	ProjectStatusDone                ProjectStatus = 6
)

var projectStatusStrings = map[int64]string{
	1: "InitPhase",
	2: "ConductingSearch",
	3: "ArticlesGathered",
	4: "ScreeningPhase",
	5: "DataExtractionPhase",
	6: "Done",
}

var projectStatusStringMap = map[string]ProjectStatus{
	"InitPhase":           ProjectStatusInitPhase,
	"ConductingSearch":    ProjectStatusConductingSearch,
	"ArticlesGathered":    ProjectStatusArticlesGathered,
	"ScreeningPhase":      ProjectStatusScreeningPhase,
	"DataExtractionPhase": ProjectStatusDataExtractionPhase,
	"Done":                ProjectStatusDone,
}

func (u *ProjectStatus) Scan(value interface{}) error {
	i, ok := value.(int64)
	if !ok {
		return fmt.Errorf("could not cast value %v to int64", value)
	}
	*u = ProjectStatus(i)
	return nil
}
func (u *ProjectStatus) ScanFromString(name string) error {
	var ok bool
	*u, ok = projectStatusStringMap[name]
	if !ok {
		return fmt.Errorf("%s is not a valid name for a value of ProjectStatus", name)
	}
	return nil
}
func (u ProjectStatus) Value() (driver.Value, error) { return int64(u), nil }
func (u ProjectStatus) String() string {
	if u == 0 {
		return "undefined"
	}
	return projectStatusStrings[int64(u)]
}
func (u ProjectStatus) AllStrings() map[int64]string { return projectStatusStrings }

// ScopeType Enum
type ScopeType int64

const (
	ScopeTypeAccess ScopeType = 1
)

var scopeTypeStrings = map[int64]string{
	1: "access",
}

var scopeTypeStringMap = map[string]ScopeType{
	"access": ScopeTypeAccess,
}

func (u *ScopeType) Scan(value interface{}) error {
	i, ok := value.(int64)
	if !ok {
		return fmt.Errorf("could not cast value %v to int64", value)
	}
	*u = ScopeType(i)
	return nil
}
func (u *ScopeType) ScanFromString(name string) error {
	var ok bool
	*u, ok = scopeTypeStringMap[name]
	if !ok {
		return fmt.Errorf("%s is not a valid name for a value of ScopeType", name)
	}
	return nil
}
func (u ScopeType) Value() (driver.Value, error) { return int64(u), nil }
func (u ScopeType) String() string {
	if u == 0 {
		return "undefined"
	}
	return scopeTypeStrings[int64(u)]
}
func (u ScopeType) AllStrings() map[int64]string { return scopeTypeStrings }
