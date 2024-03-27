package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type onetimePrekeysDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newOnetimePrekeysDB(ctx *config.Context) *onetimePrekeysDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &onetimePrekeysDB{
		db:  d,
		ctx: ctx,
	}
}

func (o *onetimePrekeysDB) insertTx(m *onetimePrekeysModel, tx *gorm.DB) error {
	err := tx.Table("signal_onetime_prekeys").Create(m).Error
	return err
}

func (o *onetimePrekeysDB) delete(uid string, keyID int) error {
	err := o.db.Table("signal_onetime_prekeys").Where("uid=? and key_id=?", uid, keyID).Delete(nil).Error
	return err
}

func (o *onetimePrekeysDB) deleteWithUID(uid string) error {
	err := o.db.Table("signal_onetime_prekeys").Where("uid=?", uid).Delete(nil).Error
	return err
}

// 查询用户最小的onetimePreKey
func (o *onetimePrekeysDB) queryMinWithUID(uid string) (*onetimePrekeysModel, error) {
	var m onetimePrekeysModel
	err := o.db.Table("signal_onetime_prekeys").Where("uid=?", uid).Order("key_id ASC").Limit(1).Find(&m).Error
	return &m, err
}

func (o *onetimePrekeysDB) queryCount(uid string) (int64, error) {
	var cn int64
	err := o.db.Table("signal_onetime_prekeys").Where("uid=?", uid).Count(&cn).Error
	return cn, err
}

type onetimePrekeysModel struct {
	UID       *string
	KeyID     *int
	Pubkey    *string
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
