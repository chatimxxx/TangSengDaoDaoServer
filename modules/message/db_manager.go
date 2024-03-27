package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

// managerDB 管理员代发消息记录
type managerDB struct {
	db *gorm.DB
	t  *DB
}

// newManagerDB newManagerDB
func newManagerDB(ctx *config.Context) *managerDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &managerDB{
		db: db,
		t:  NewDB(ctx),
	}
}

// 添加一条发送消息记录
func (m *managerDB) insertMsgHistory(message *managerMsgModel) error {
	err := m.db.Table("send_history").Create(message).Error
	return err
}

// 查询代发消息记录
func (m *managerDB) queryMsgWithPage(pageSize, page int) ([]*managerMsgModel, error) {
	var ms []*managerMsgModel
	err := m.db.Table("send_history").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

// 查询消息数量
func (m *managerDB) queryMsgCount() (int64, error) {
	var count int64
	err := m.db.Table("send_history").Count(&count).Error
	return count, err
}

func (m *managerDB) queryWithChannelID(channelID string, page, pageSize int) ([]*messageModel, error) {
	var ms []*messageModel
	var table = m.t.getTable(channelID)
	err := m.db.Table(table).Where("channel_id=?", channelID).Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

func (m *managerDB) queryRecordCount(channelID string) (int64, error) {
	var count int64
	err := m.db.Table(m.t.getTable(channelID)).Where("channel_id=?", channelID).Count(&count).Error
	return count, err
}

func (m *managerDB) queryMsgExtrWithMsgIds(msgIds []int64) ([]*messageExtraModel, error) {
	var ms []*messageExtraModel
	err := m.db.Table("message_extra").Where("message_id in ?", msgIds).Find(&ms).Error
	return ms, err
}

func (m *managerDB) queryProhibitWordsWithContent(content string) (*prohibitWordsModel, error) {
	var model prohibitWordsModel
	err := m.db.Table("prohibit_words").Where("content=?", content).First(&model).Error
	return &model, err
}

func (m *managerDB) queryProhibitWordsWithID(id int) (*prohibitWordsModel, error) {
	var model prohibitWordsModel
	err := m.db.Table("prohibit_words").Where("id=?", id).Find(&model).Error
	return &model, err
}

func (m *managerDB) updateProhibitWord(word *prohibitWordsModel) error {
	err := m.db.Table("prohibit_words").Updates(map[string]interface{}{
		"version":    word.Version,
		"is_deleted": word.IsDeleted,
	}).Where("content=?", word.Content).Error
	return err
}
func (m *managerDB) insertProhibitWord(word *prohibitWordsModel) error {
	err := m.db.Table("prohibit_words").Create(word).Error
	return err
}
func (m *managerDB) queryProhibitWords(pageIndex, pageSize int) ([]*prohibitWordsModel, error) {
	var ms []*prohibitWordsModel
	err := m.db.Table("prohibit_words").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

func (m *managerDB) queryProhibitWordsWithContentAndPage(content string, pageIndex, pageSize int) ([]*prohibitWordsModel, error) {
	var ms []*prohibitWordsModel
	err := m.db.Table("prohibit_words").Where("content like ?", "%"+content+"%").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

func (m *managerDB) queryProhibitWordsCount() (int64, error) {
	var count int64
	err := m.db.Table("prohibit_words").Count(&count).Error
	return count, err
}

func (m *managerDB) queryProhibitWordsCountWithContent(content string) (int64, error) {
	var count int64
	err := m.db.Table("prohibit_words").Where("content like ?", "%"+content+"%").Count(&count).Error
	return count, err
}

func (m *managerDB) updateMsgExtraVersionAndDeletedTx(md *messageExtraModel, tx *gorm.DB) error {
	err := tx.Exec("INSERT INTO message_extra (message_id,message_seq,channel_id,channel_type,is_deleted,version) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE is_deleted=VALUES(is_deleted),version=VALUES(version)", md.MessageID, md.MessageSeq, md.ChannelID, md.ChannelType, md.IsDeleted, md.Version).Error
	return err
}

type prohibitWordsModel struct {
	Content   *string
	IsDeleted *int
	Version   *int64
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// 管理员代发消息记录
type managerMsgModel struct {
	Receiver            *string // 接受者uid
	ReceiverName        *string // 接受者名字
	ReceiverChannelType *int    // 接受者频道类型
	Sender              *string // 发送者uid
	SenderName          *string // 发送者名字
	HandlerUID          *string // 操作者uid
	HandlerName         *string // 操作者名字
	Content             *string // 发送内容
	Id                  *int64
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}
