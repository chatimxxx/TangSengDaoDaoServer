package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type conversationExtraDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newConversationExtraDB(ctx *config.Context) *conversationExtraDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &conversationExtraDB{
		ctx: ctx,
		db:  db,
	}
}

func (c *conversationExtraDB) insertOrUpdate(model *conversationExtraModel) error {
	err := c.db.Exec("INSERT INTO conversation_extra (uid,channel_id,channel_type,browse_to,keep_message_seq,keep_offset_y,draft,version) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE browse_to=IF(VALUES(browse_to)>browse_to,VALUES(browse_to),browse_to),`keep_message_seq`=VALUES(`keep_message_seq`),keep_offset_y=VALUES(keep_offset_y),draft=VALUES(draft),version=VALUES(version)", model.UID, model.ChannelID, model.ChannelType, model.BrowseTo, model.KeepMessageSeq, model.KeepOffsetY, model.Draft, model.Version).Error
	return err
}

func (c *conversationExtraDB) sync(uid string, version int64) ([]*conversationExtraModel, error) {
	var ms []*conversationExtraModel
	err := c.db.Table("conversation_extra").Where("uid=? and version>?", uid, version).Find(&ms).Error
	return ms, err
}

func (c *conversationExtraDB) queryWithChannelIDs(uid string, channelIDs []string) ([]*conversationExtraModel, error) {
	if len(channelIDs) == 0 {
		return nil, nil
	}
	var ms []*conversationExtraModel
	err := c.db.Table("conversation_extra").Where("uid=? and channel_id in ?", uid, channelIDs).Find(&ms).Error
	return ms, err
}

type conversationExtraModel struct {
	UID            *string
	ChannelID      *string
	ChannelType    *uint8
	BrowseTo       *uint32
	KeepMessageSeq *uint32
	KeepOffsetY    *int
	Draft          *string // 草稿
	Version        *int64
	Id             *int64
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
