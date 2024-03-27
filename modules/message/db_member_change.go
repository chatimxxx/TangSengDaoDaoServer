package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type memberChangeDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newMemberChangeDB(ctx *config.Context) *memberChangeDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &memberChangeDB{
		ctx: ctx,
		db:  db,
	}
}

// 查询频道成员最大版本号
func (m *memberChangeDB) queryMaxVersion(channelID string, channelType uint8) (*memberChangeModel, error) {
	var model memberChangeModel
	err := m.db.Table("member_change").Where("channel_id=? and channel_type=?", channelID, channelType).Order("max_version DESC").Limit(1).First(&model).Error
	return &model, err
}

func (m *memberChangeDB) insertTx(model *memberChangeModel, tx *gorm.DB) error {
	err := tx.Table("member_change").Create(model).Error
	return err
}

func (m *memberChangeDB) updateMaxVersion(maxVersion int64, channelID string, channelType uint8) error {
	err := m.db.Table("member_change").Update("max_version", maxVersion).Where("channel_id=? and channel_type=?", channelID, channelType).Error
	return err
}

type memberChangeModel struct {
	CloneNo     *string
	ChannelID   *string
	ChannelType *uint8
	MaxVersion  *int64
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
