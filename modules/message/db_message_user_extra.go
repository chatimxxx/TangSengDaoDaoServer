package message

import (
	"fmt"
	"gorm.io/gorm"
	"hash/crc32"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
)

type messageUserExtraDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newMessageUserExtraDB(ctx *config.Context) *messageUserExtraDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &messageUserExtraDB{ctx: ctx, db: db}
}

// 插入或更新消息为已删除
func (m *messageUserExtraDB) insertOrUpdateDeleted(md *messageUserExtraModel) error {
	sq := fmt.Sprintf("INSERT INTO %s (uid,message_id,message_seq,channel_id,channel_type,message_is_deleted) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE  message_is_deleted=VALUES(message_is_deleted)", m.getTable(*md.UID))
	err := m.db.Exec(sq, md.UID, md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.MessageIsDeleted).Error
	return err
}
func (m *messageUserExtraDB) insertOrUpdateDeletedTx(md *messageUserExtraModel, tx *gorm.DB) error {
	sq := fmt.Sprintf("INSERT INTO %s (uid,message_id,message_seq,channel_id,channel_type,message_is_deleted) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE  message_is_deleted=VALUES(message_is_deleted)", m.getTable(*md.UID))
	err := tx.Exec(sq, md.UID, md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.MessageIsDeleted).Error
	return err
}

// 插入或更新消息语音已读状态
func (m *messageUserExtraDB) insertOrUpdateVoiceRead(md *messageUserExtraModel) error {
	sq := fmt.Sprintf("INSERT INTO %s (uid,message_id,message_seq,channel_id,channel_type,voice_readed) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE  voice_readed=VALUES(voice_readed)", m.getTable(*md.UID))
	err := m.db.Exec(sq, md.UID, md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.VoiceReaded).Error
	return err
}

// 通过消息id集合和消息拥有者uid查询编辑消息
func (m *messageUserExtraDB) queryWithMessageIDsAndUID(messageIDs []string, uid string) ([]*messageUserExtraModel, error) {
	if len(messageIDs) == 0 {
		return nil, nil
	}
	var ms []*messageUserExtraModel
	err := m.db.Table(m.getTable(uid)).Where("uid=? and message_id in ?", uid, messageIDs).Find(&ms).Error
	return ms, err
}

func (m *messageUserExtraDB) getTable(uid string) string {
	tableIndex := crc32.ChecksumIEEE([]byte(uid)) % uint32(m.ctx.GetConfig().TablePartitionConfig.MessageUserEditTableCount)
	if tableIndex == 0 {
		return "message_user_extra"
	}
	return fmt.Sprintf("message_user_extra%d", tableIndex)
}

type messageUserExtraModel struct {
	UID              *string
	MessageID        *string
	MessageSeq       *uint32
	ChannelID        *string
	ChannelType      *uint8
	VoiceReaded      *int
	MessageIsDeleted *int
	Id               *int64
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}
