package event

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/util"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/wkevent"
	"gorm.io/gorm"
	"time"
)

// DB 事件的db
type DB struct {
	db *gorm.DB
}

// NewDB 创建DB
func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// InsertTx 插入事件
func (d *DB) InsertTx(m *Model, tx *gorm.DB) (int64, error) {
	err := tx.Table("event").Create(m).Error
	if err != nil {
		return 0, err
	}
	return *m.Id, nil
}

// UpdateStatus 更新事件状态
func (d *DB) UpdateStatus(reason string, status int, versionLock int64, id int64) error {
	m := Model{
		Status: &status,
		Reason: &reason,
	}
	err := d.db.Table("event").Where("id=? and version_lock=?", id, versionLock).Updates(&m).Error
	return err
}

// QueryWithID 根据id查询事件
func (d *DB) QueryWithID(id int64) (*Model, error) {
	var model Model
	err := d.db.Table("event").Where("id=?", id).First(&model).Error
	return &model, err
}

// QueryAllWait 查询所有等待事件
func (d *DB) QueryAllWait(limit int) ([]*Model, error) {
	var models []*Model
	date := util.ToyyyyMMddHHmmss(time.Now().Add(-time.Second * 60))
	err := d.db.Table("event").Where("status=? and created_at<?", wkevent.Wait.Int(), date).Limit(limit).Find(&models).Error
	return models, err
}

// ---------- model ----------

// Model 数据库对象
type Model struct {
	Event       *string // 事件标示
	Type        *int    // 事件类型
	Data        *string // 事件数据
	Status      *int    // 事件状态 0.待发布 1.已发布 2.发布失败！
	Reason      *string // 原因 如果状态为2，则有发布失败的原因
	VersionLock *int64  // 乐观锁
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
