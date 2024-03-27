package elastic

import (
	"gorm.io/gorm"
)

// DB DB
type DB struct {
	db *gorm.DB
}

// NewDB NewDB
func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// Insert Insert
func (d *DB) Insert(model *IndexerErrorModel) error {
	err := d.db.Table("indexer_error").Create(model).Error
	return err
}

// IndexerErrorModel IndexerErrorModel
type IndexerErrorModel struct {
	Index      *string
	Action     *string
	DocumentID *string
	Body       *string
	Error      *string
	Id         *int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
