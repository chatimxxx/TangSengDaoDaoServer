package message

import (
	"fmt"
	"gorm.io/gorm"
	"hash/crc32"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
)

type channelOffsetDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newChannelOffsetDB(ctx *config.Context) *channelOffsetDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &channelOffsetDB{
		ctx: ctx,
		db:  db,
	}
}

func (c *channelOffsetDB) insertOrUpdate(m *channelOffsetModel) error {
	sq := fmt.Sprintf("INSERT INTO %s (uid,channel_id,channel_type,message_seq) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE message_seq=IF(message_seq<VALUES(message_seq),VALUES(message_seq),message_seq)", c.getTable(*m.UID))
	err := c.db.Exec(sq, m.UID, m.ChannelID, m.ChannelType, m.MessageSeq).Error
	return err
}

func (c *channelOffsetDB) delete(uid string, channelID string, channelType uint8, tx *gorm.DB) error {
	err := tx.Table(c.getTable(uid)).Where("uid=? and channel_id=? and channel_type=?", uid, channelID, channelType).Delete(nil).Error
	return err
}

func (c *channelOffsetDB) insertOrUpdateTx(m *channelOffsetModel, tx *gorm.DB) error {
	sq := fmt.Sprintf("INSERT INTO %s (uid,channel_id,channel_type,message_seq) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE  message_seq=IF(message_seq<VALUES(message_seq),VALUES(message_seq),message_seq)", c.getTable(*m.UID))
	err := tx.Exec(sq, m.UID, m.ChannelID, m.ChannelType, m.MessageSeq).Error
	return err
}

func (c *channelOffsetDB) queryWithUIDAndChannel(uid string, channelID string, channelType uint8) (*channelOffsetModel, error) {
	var m *channelOffsetModel
	err := c.db.Table(c.getTable(uid)).Where("(uid=? or uid='') and channel_id=? and channel_type=?", uid, channelID, channelType).Order("message_seq DESC").Limit(1).First(&m).Error
	return m, err
}

func (c *channelOffsetDB) queryWithUIDAndChannelIDs(uid string, channelIDs []string) ([]*channelOffsetModel, error) {
	var ms []*channelOffsetModel
	err := c.db.Select("channel_id,channel_type,max(message_seq) message_seq").Table(c.getTable(uid)).Where("(uid=? or uid='') and channel_id in ?", uid, channelIDs).Group("channel_id, channel_type").Find(&ms).Error
	return ms, err
}

func (c *channelOffsetDB) getTable(uid string) string {
	tableIndex := crc32.ChecksumIEEE([]byte(uid)) % uint32(c.ctx.GetConfig().TablePartitionConfig.ChannelOffsetTableCount)
	if tableIndex == 0 {
		return "channel_offset"
	}
	return fmt.Sprintf("channel_offset%d", tableIndex)
}

type channelOffsetModel struct {
	UID         *string
	ChannelID   *string
	ChannelType *uint8
	MessageSeq  *uint32
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
