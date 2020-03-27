package models

import (
	"context"
	"strings"
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
)

func (m *UserDB) GetOnEmail(ctx context.Context, email string) (*User, error) {
	defer goa.MeasureSince([]string{"goa", "db", "user", "get"}, time.Now())

	var native User
	err := m.Db.Table(m.TableName()).Where("lower(email) = ?", strings.ToLower(email)).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}
