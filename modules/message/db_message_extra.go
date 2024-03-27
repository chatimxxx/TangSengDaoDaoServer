package message

import (
	"gorm.io/gorm"
	"sort"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/gocraft/dbr/v2"
)

type messageExtraDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newMessageExtraDB(ctx *config.Context) *messageExtraDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &messageExtraDB{
		ctx: ctx,
		db:  db,
	}
}

func (m *messageExtraDB) insertOrUpdateRevoke(md *messageExtraModel) error {
	err := m.db.Exec("INSERT INTO message_extra (message_id,message_seq,channel_id,channel_type,`revoke`,revoker,version) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `revoke`=VALUES(`revoke`),revoker=VALUES(revoker),version=VALUES(version)", md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.Revoke, md.Revoker, md.Version).Error
	return err
}

func (m *messageExtraDB) insertOrUpdateRevokeTx(md *messageExtraModel, tx *gorm.DB) error {
	err := tx.Exec("INSERT INTO message_extra (message_id,message_seq,channel_id,channel_type,`revoke`,revoker,version) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `revoke`=VALUES(`revoke`),revoker=VALUES(revoker),version=VALUES(version)", md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.Revoke, md.Revoker, md.Version).Error
	return err
}

// 更新已读数量
func (m *messageExtraDB) insertOrUpdateReadedCount(md *messageExtraModel) error {
	err := m.db.Exec("INSERT INTO message_extra (clone_no,message_id,message_seq,from_uid,channel_id,channel_type,readed_count,version) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE clone_no=IF(clone_no='',VALUES(clone_no),clone_no),readed_count=VALUES(readed_count),version=VALUES(version)", md.CloneNo, md.MessageID, md.MessageSeq, md.FromUID, md.ChannelID, md.ChannelType, md.ReadedCount, md.Version).Error
	return err
}

// 更新已读数量
func (m *messageExtraDB) insertOrUpdateReadedCountTx(md *messageExtraModel, tx *gorm.DB) error {
	err := tx.Exec("INSERT INTO message_extra (clone_no,message_id,message_seq,from_uid,channel_id,channel_type,readed_count,version) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE clone_no=IF(clone_no='',VALUES(clone_no),clone_no),readed_count=VALUES(readed_count),version=VALUES(version)", md.CloneNo, md.MessageID, md.MessageSeq, md.FromUID, md.ChannelID, md.ChannelType, md.ReadedCount, md.Version).Error
	return err
}

func (m *messageExtraDB) insertOrUpdateContentEdit(md *messageExtraModel) error {
	err := m.db.Exec("INSERT INTO message_extra (message_id,message_seq,channel_id,channel_type,content_edit,content_edit_hash,edited_at,version) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE content_edit=VALUES(content_edit),content_edit_hash=VALUES(content_edit_hash),edited_at=VALUES(edited_at),version=VALUES(version)", md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.ContentEdit, md.ContentEditHash, md.EditedAt, md.Version).Error
	return err
}

// 是否存在相同编辑内容
func (m *messageExtraDB) existContentEdit(messageID string, contentEditHash string) (bool, error) {
	var count int64
	err := m.db.Table("message_extra").Where("message_id=? and content_edit_hash=?", messageID, contentEditHash).Count(&count).Error
	return count > 0, err
}

func (m *messageExtraDB) queryWithMessageIDs(messageIDs []string, loginUID string) ([]*messageExtraDetailModel, error) {
	if len(messageIDs) <= 0 {
		return nil, nil
	}
	var ms []*messageExtraDetailModel
	err := m.db.Select("message_extra.*,(select count(*) from member_readed where member_readed.message_id=message_extra.message_id and member_readed.uid='"+loginUID+"') readed,(select created_at from member_readed where member_readed.message_id=message_extra.message_id and member_readed.uid='"+loginUID+"') readed_at").Table("message_extra").Where("message_id in ?", messageIDs).Find(&ms).Error
	return ms, err
}

func (m *messageExtraDB) queryWithMessageID(messageID int64) (*messageExtraModel, error) {
	var model messageExtraModel
	err := m.db.Table("message_extra").Where("message_id=?", messageID).First(&model).Error
	return &model, err
}

func (m *messageExtraDB) sync(version int64, channelID string, channelType uint8, limit int, loginUID string) ([]*messageExtraDetailModel, error) {
	var models []*messageExtraDetailModel
	selectSql := "message_extra.*,(select count(*) from member_readed where member_readed.message_id=message_extra.message_id and member_readed.uid='" + loginUID + "') readed,(select created_at from member_readed where member_readed.message_id=message_extra.message_id and member_readed.uid='" + loginUID + "') readed_at"
	builder := m.db.Select(selectSql).Table("message_extra")
	var err error
	if version == 0 {
		builder = builder.Where("channel_id=? and channel_type=?", channelID, channelType).Order("version DESC").Limit(limit)
		err = builder.Find(&models).Error
		newModels := messageExtraDetailModelSlice(models)
		sort.Sort(newModels)
		models = newModels
	} else {
		builder = builder.Where("channel_id=? and channel_type=? and version>?", channelID, channelType, version).Order("version ASC").Limit(limit)
		err = builder.Find(&models).Error
	}

	return models, err
}

type messageExtraDetailModelSlice []*messageExtraDetailModel

func (m messageExtraDetailModelSlice) Len() int {
	return len(m)
}
func (m messageExtraDetailModelSlice) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m messageExtraDetailModelSlice) Less(i, j int) bool { return *m[i].Version < *m[j].Version }

type messageExtraDetailModel struct {
	messageExtraModel
	Readed   int          // 是否已读（针对于自己）
	ReadedAt dbr.NullTime // 已读时间

}

type messageExtraModel struct {
	MessageID       *string
	MessageSeq      *uint32
	FromUID         *string
	ChannelID       *string
	ChannelType     *uint8
	Revoke          *int
	Revoker         *string // 消息撤回者的uid
	CloneNo         *string
	ReadedCount     *int    // 已读数量
	ContentEdit     *string // 编辑后的正文
	ContentEditHash *string
	EditedAt        *int // 编辑时间 时间戳（秒）
	IsDeleted       *int
	Version         *int64 // 数据版本
	Id              *int64
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
