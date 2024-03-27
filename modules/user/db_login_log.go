package user

import (
	"gorm.io/gorm"
	"time"
)

// LoginLogDB 登录日志DB
type LoginLogDB struct {
	db *gorm.DB
}

// NewLoginLogDB NewDB
func NewLoginLogDB(db *gorm.DB) *LoginLogDB {
	return &LoginLogDB{
		db: db,
	}
}

// insert 添加登录日志
func (l *LoginLogDB) insert(m *LoginLogModel) error {
	err := l.db.Table("login_log").Create(m).Error
	return err
}

// queryLastLoginIP 查询最后一次登录日志
func (l *LoginLogDB) queryLastLoginIP(uid string) (*LoginLogModel, error) {
	var m LoginLogModel
	err := l.db.Table("login_log").Where("uid=?", uid).Order("created_at DESC").Limit(1).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// LoginLogModel 登录日志
type LoginLogModel struct {
	LoginIP   *string //登录IP
	UID       *string
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
