package common

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	dbs "github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type shortnoDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newShortnoDB(ctx *config.Context) *shortnoDB {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &shortnoDB{
		ctx: ctx,
		db:  db,
	}
}

func (s *shortnoDB) inserts(shortnos []string) error {
	if len(shortnos) == 0 {
		return nil
	}
	tx := s.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	for _, st := range shortnos {
		err := tx.Exec("insert into shortno(shortno) values(?)", st).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

func (s *shortnoDB) queryVail() (*shortnoModel, error) {
	var m *shortnoModel
	err := s.db.Table("shortno").Where("used=0 and hold=0 and locked=0").Limit(1).Find(&m).Error
	return m, err
}

func (s *shortnoDB) updateLock(shortno string, lock int) error {
	err := s.db.Table("shortno").Update("locked", lock).Where("shortno=?", shortno).Error
	return err
}

func (s *shortnoDB) updateUsed(shortno string, used int, business string) error {
	err := s.db.Table("shortno").Update("used", used).Update("business", business).Where("shortno=?", shortno).Error
	return err
}
func (s *shortnoDB) updateHold(shortno string, hold int) error {
	err := s.db.Table("shortno").Update("hold", hold).Where("shortno=?", shortno).Error
	return err
}

func (s *shortnoDB) queryVailCount() (int64, error) {
	var cn int64
	err := s.db.Table("shortno").Where("used=0 and hold=0 and locked=0").Count(&cn).Error
	return cn, err
}

type shortnoModel struct {
	Shortno  string
	Used     int
	Hold     int
	Locked   int
	Business string
	dbs.BaseModel
}
