package app

import (
	"gorm.io/gorm"
	"time"
)

// DB DB
type DB struct {
	db *gorm.DB
}

func newDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

func (d *DB) queryWithAppID(appID string) (*model, error) {
	var m *model
	err := d.db.Table("app").Where("app_id=?", appID).First(&m).Error
	return m, err
}

func (d *DB) existWithAppID(appID string) (bool, error) {
	var count int64
	err := d.db.Table("app").Where("app_id=?", appID).Count(&count).Error
	return count > 0, err
}

func (d *DB) insert(m *model) error {
	return d.db.Table("app").Create(m).Error
}

type model struct {
	AppID     *string
	AppKey    *string
	AppName   *string
	AppLogo   *string
	Status    *int
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
