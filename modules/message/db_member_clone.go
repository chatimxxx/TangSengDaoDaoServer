package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type memberCloneDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newMemberCloneDB(ctx *config.Context) *memberCloneDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &memberCloneDB{
		ctx: ctx,
		db:  db,
	}
}

func (m *memberCloneDB) insertTx(model *memberCloneModel, tx *gorm.DB) error {
	err := tx.Table("member_clone").Create(model).Error
	return err
}

func (m *memberCloneDB) queryWithCloneNo(cloneNo string) ([]*memberCloneModel, error) {
	var ms []*memberCloneModel
	err := m.db.Table("member_clone").Where("clone_no=?", cloneNo).Find(&ms).Error
	return ms, err
}

// 查询未读列表
func (m *memberCloneDB) queryUnreadWithMessageIDAndPage(cloneNo string, fromUID string, messageID int64, pIndex, pSize int) ([]*memberUnreadModel, error) {
	var ms []*memberUnreadModel
	err := m.db.Table("member_clone").Where("clone_no=? and uid<>? and uid not in (select member_readed.uid  from member_readed where message_id=?)", cloneNo, fromUID, messageID).Limit(pSize).Offset((pIndex - 1) * pSize).Find(&ms).Error
	return ms, err
}

type memberCloneModel struct {
	CloneNo     *string
	ChannelID   *string
	ChannelType *uint8
	UID         *string
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type memberUnreadModel struct {
	CloneNo     string
	ChannelID   string
	ChannelType uint8
	UID         string
}
