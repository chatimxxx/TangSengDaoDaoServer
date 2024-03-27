package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type memberReadedDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newMemberReadedDB(ctx *config.Context) *memberReadedDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &memberReadedDB{
		ctx: ctx,
		db:  db,
	}
}

func (m *memberReadedDB) insertOrUpdateTx(model *memberReadedModel, tx *gorm.DB) error {
	err := m.db.Exec("INSERT INTO member_readed (message_id,clone_no,channel_id,channel_type,uid) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `message_id`=VALUES(`message_id`),`clone_no`=VALUES(`clone_no`),uid=VALUES(uid)", model.MessageID, model.CloneNo, model.ChannelID, model.ChannelType, model.UID).Error
	return err
}

// 查询消息已读数量
func (m *memberReadedDB) queryCountWithMessageIDs(channelID string, channelType uint8, messageIDs []string) (map[int64]int, error) {
	if len(messageIDs) <= 0 {
		return nil, nil
	}
	var ms []struct {
		MessageID int64
		Num       int
	}
	err := m.db.Select("member_readed.message_id,count(*) num").Table("member_readed").Where("member_readed.channel_id=? and member_readed.channel_type=? and member_readed.message_id in ?", channelID, channelType, messageIDs).Group("member_readed.message_id, member_readed.channel_id, member_readed.channel_type").Find(&ms).Error
	if err != nil {
		return nil, err
	}
	resultMap := map[int64]int{}
	if len(ms) > 0 {
		for _, m := range ms {
			resultMap[m.MessageID] = m.Num
		}
	}
	return resultMap, nil
}

// 查询已读列表
func (m *memberReadedDB) queryWithMessageIDAndPage(messageID string, pIndex, pSize int) ([]*memberReadedModel, error) {
	var ms []*memberReadedModel
	err := m.db.Table("member_readed").Where("member_readed.message_id=?", messageID).Order("created_at DESC").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&ms).Error
	return ms, err
}

type memberReadedModel struct {
	CloneNo     *string // TODO: 此字段作废
	MessageID   *int64
	ChannelID   *string
	ChannelType *uint8
	UID         *string
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
