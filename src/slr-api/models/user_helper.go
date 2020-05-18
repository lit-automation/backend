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

// ListUser returns an array of view: default.
func (m *UserDB) ListUser(ctx context.Context) []*app.User {
	defer goa.MeasureSince([]string{"goa", "db", "user", "listuser"}, time.Now())

	var native []*User
	var objs []*app.User
	err := m.Db.Scopes().Table(m.TableName()).Find(&native).Error

	if err != nil {
		goa.LogError(ctx, "error listing User", "error", err.Error())
		return objs
	}

	for _, t := range native {
		objs = append(objs, t.UserToUser())
	}

	return objs
}

// UserToUser loads a User and builds the default view of media type User.
func (m *User) UserToUser() *app.User {
	user := &app.User{}
	user.Email = m.Email
	user.FamilyName = m.FamilyName
	user.FirstName = m.FirstName
	user.ID = m.ID
	user.MiddleName = &m.MiddleName

	return user
}

// OneUser loads a User and builds the default view of media type User.
func (m *UserDB) OneUser(ctx context.Context, id uuid.UUID) (*app.User, error) {
	defer goa.MeasureSince([]string{"goa", "db", "user", "oneuser"}, time.Now())

	var native User
	err := m.Db.Scopes().Table(m.TableName()).Preload("Projects").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		goa.LogError(ctx, "error getting User", "error", err.Error())
		return nil, err
	}

	view := *native.UserToUser()
	return &view, err
}

// CRUD Functions

// Get returns a single User as a Database Model
// This is more for use internally, and probably not what you want in your controllers
func (m *UserDB) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	defer goa.MeasureSince([]string{"goa", "db", "user", "get"}, time.Now())

	var native User
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of User
func (m *UserDB) List(ctx context.Context) ([]*User, error) {
	defer goa.MeasureSince([]string{"goa", "db", "user", "list"}, time.Now())

	var objs []*User
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *UserDB) Add(ctx context.Context, model *User) error {
	defer goa.MeasureSince([]string{"goa", "db", "user", "add"}, time.Now())

	err := m.Db.Create(model).Error
	if err != nil {
		goa.LogError(ctx, "error adding User", "error", err.Error())
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *UserDB) Update(ctx context.Context, model *User) error {
	defer goa.MeasureSince([]string{"goa", "db", "user", "update"}, time.Now())

	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating User", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *UserDB) Delete(ctx context.Context, id uuid.UUID) error {
	defer goa.MeasureSince([]string{"goa", "db", "user", "delete"}, time.Now())

	err := m.Db.Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		goa.LogError(ctx, "error deleting User", "error", err.Error())
		return err
	}

	return nil
}
