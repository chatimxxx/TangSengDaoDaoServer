package elastic

import (
	"github.com/gocraft/dbr/v2"
	"github.com/xochat/xochat_im_server_lib/pkg/db"
	"github.com/xochat/xochat_im_server_lib/pkg/util"
)

// DB DB
type DB struct {
	session *dbr.Session
}

// NewDB NewDB
func NewDB(session *dbr.Session) *DB {
	return &DB{
		session: session,
	}
}

// Insert Insert
func (d *DB) Insert(model *IndexerErrorModel) error {
	_, err := d.session.InsertInto("indexer_error").Columns(util.AttrToUnderscore(model)...).Record(model).Exec()
	return err
}

// IndexerErrorModel IndexerErrorModel
type IndexerErrorModel struct {
	Index      string
	Action     string
	DocumentID string
	Body       string
	Error      string
	db.BaseModel
}
