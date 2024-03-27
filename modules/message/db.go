package message

import (
	"fmt"
	"gorm.io/gorm"
	"hash/crc32"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
)

// DB DB
type DB struct {
	db  *gorm.DB
	ctx *config.Context
}

// NewDB NewDB
func NewDB(ctx *config.Context) *DB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &DB{
		db:  db,
		ctx: ctx,
	}
}

// InsertTx 添加消息
// func (d *DB) InsertTx(m *Model, tx *gorm.DB) error {
// 	err := tx.InsertInto("message").Create(m).Error
// 	return err
// }

func (d *DB) queryMessageWithMessageID(channelID string, channelType uint8, messageID string) (*messageModel, error) {
	var m *messageModel
	err := d.db.Table(d.getTable(channelID)).Where("message_id=?", messageID).First(&m).Error
	return m, err
}

func (d *DB) queryMessagesWithMessageIDs(channelID string, channelType uint8, messageIDs []string) ([]*messageModel, error) {
	if len(messageIDs) <= 0 {
		return nil, nil
	}
	var ms []*messageModel
	err := d.db.Table(d.getTable(channelID)).Where("message_id in ?", messageIDs).Find(&ms).Error
	return ms, err
}

func (d *DB) queryMaxMessageSeq(channelID string, channelType uint8) (uint32, error) {
	var maxMessageSeq uint32
	err := d.db.Select("IFNULL(max(message_seq),0)").Table(d.getTable(channelID)).Where("channel_id=? and channel_type=?", channelID, channelType).First(&maxMessageSeq).Error
	return maxMessageSeq, err
}

func (d *DB) queryMessagesWithChannelClientMsgNo(channelID string, channelType uint8, clientMsgNo string) ([]*messageModel, error) {
	var ms []*messageModel
	err := d.db.Table(d.getTable(channelID)).Where("channel_id=? and channel_type=? and client_msg_no=?", channelID, channelType, clientMsgNo).Find(&ms).Error
	return ms, err
}
func (d *DB) queryProhibitWordsWithVersion(version int64) ([]*ProhibitWordModel, error) {
	var ms []*ProhibitWordModel
	err := d.db.Table("prohibit_words").Where("`version` > ?", version).Find(&ms).Error
	return ms, err
}

// 通过频道ID获取表
func (d *DB) getTable(channelID string) string {
	tableIndex := crc32.ChecksumIEEE([]byte(channelID)) % uint32(d.ctx.GetConfig().TablePartitionConfig.MessageTableCount)
	if tableIndex == 0 {
		return "message"
	}
	return fmt.Sprintf("message%d", tableIndex)
}

// ProhibitWordModel 违禁词model
type ProhibitWordModel struct {
	Content   *string
	IsDeleted *int
	Version   *int64
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// Model 消息model
type messageModel struct {
	MessageID   *int64
	MessageSeq  *uint32
	ClientMsgNo *string
	Header      *string
	Setting     *uint8
	FromUID     *string
	ChannelID   *string
	ChannelType *uint8
	Timestamp   *int64
	Type        *int
	Payload     *[]byte
	IsDeleted   *int
	Signal      *int
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
