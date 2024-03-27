package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type messageReactionDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newMessageReactionDB(ctx *config.Context) *messageReactionDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &messageReactionDB{
		ctx: ctx,
		db:  db,
	}
}

// 查询某个频道的回应数据
func (d *messageReactionDB) queryReactionWithChannelAndSeq(channelID string, channelType uint8, seq int64, limit int) ([]*reactionModel, error) {
	var ms []*reactionModel
	var err error
	if seq <= 0 { // TODO: 如果seq为0 不能去同步整个频道的 应该同步最新指定数量的回应数据（建议limit 100）
		err = d.db.Table("reaction_users").Where("channel_id=? and channel_type=?", channelID, channelType).Order("seq DESC").Limit(limit).Find(&ms).Error
	} else {
		err = d.db.Table("reaction_users").Where("channel_id=? and channel_type=? and seq>?", channelID, channelType, seq).Order("seq ASC").Limit(limit).Find(&ms).Error
	}
	return ms, err
}

func (d *messageReactionDB) queryWithMessageIDs(messageIDs []string) ([]*reactionModel, error) {
	if len(messageIDs) <= 0 {
		return nil, nil
	}
	var ms []*reactionModel
	err := d.db.Table("reaction_users").Where("message_id in ?", messageIDs).Find(&ms).Error
	return ms, err
}

// 查询某个用户的回应数据
func (d *messageReactionDB) queryReactionWithUIDAndMessageID(uid string, messageID string) (*reactionModel, error) {
	var model reactionModel
	err := d.db.Table("reaction_users").Where("uid=? and message_id=?", uid, messageID).First(&model).Error
	return &model, err
}

// 新增回应
func (d *messageReactionDB) insertReaction(model *reactionModel) error {
	err := d.db.Table("reaction_users").Create(model).Error
	return err
}

// 修改某条消息的回应
func (d *messageReactionDB) updateReactionStatus(model *reactionModel) error {
	err := d.db.Table("reaction_users").Updates(map[string]interface{}{
		"is_deleted": model.IsDeleted,
		"seq":        model.Seq,
		"emoji":      model.Emoji,
	}).Where("message_id=? and uid=?", model.MessageID, model.UID).Error
	return err
}
func (d *messageReactionDB) updateReactionText(model *reactionModel) error {
	err := d.db.Table("reaction_users").Updates(map[string]interface{}{
		"is_deleted": model.IsDeleted,
		"seq":        model.Seq,
	}).Where("message_id=? and uid=? and emoji=?", model.MessageID, model.UID, model.Emoji).Error
	return err
}

type reactionModel struct {
	MessageID   *string // 消息唯一ID
	Seq         *int64  // 回复序列号
	ChannelID   *string // 频道唯一ID
	ChannelType *uint8  // 频道类型
	UID         *string // 用户ID
	Name        *string // 用户名称
	Emoji       *string // 回应表情
	IsDeleted   *int    // 是否已删除
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
