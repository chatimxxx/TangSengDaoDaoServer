package channel

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type channelSettingDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newChannelSettingDB(ctx *config.Context) *channelSettingDB {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &channelSettingDB{
		db:  db,
		ctx: ctx,
	}
}

func (c *channelSettingDB) queryWithChannel(channelID string, channelType uint8) (*channelSettingModel, error) {
	var m *channelSettingModel
	err := c.db.Table("channel_setting").Where("channel_id=? and channel_type=?", channelID, channelType).First(&m).Error
	return m, err
}

func (c *channelSettingDB) queryWithChannelIDs(channelIDs []string) ([]*channelSettingModel, error) {
	var ms []*channelSettingModel
	err := c.db.Table("channel_setting").Where("channel_id in ?", channelIDs).Find(&ms).Error
	return ms, err
}

func (c *channelSettingDB) insertOrAddMsgAutoDelete(channelID string, channelType uint8, msgAutoDelete int64) error {
	err := c.db.Exec("insert into channel_setting (channel_id, channel_type, msg_auto_delete) values (?, ?, ?) ON DUPLICATE KEY UPDATE msg_auto_delete=VALUES(msg_auto_delete)", channelID, channelType, msgAutoDelete).Error
	return err
}

type channelSettingModel struct {
	ChannelID         *string
	ChannelType       *uint8
	ParentChannelID   *string
	ParentChannelType *uint8
	MsgAutoDelete     *int64
	Id                *int64
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}
