package message

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type userLastOffsetDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newUserLastOffsetDB(ctx *config.Context) *userLastOffsetDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &userLastOffsetDB{
		ctx: ctx,
		db:  db,
	}
}

func (d *userLastOffsetDB) insertOrUpdateTx(tx *gorm.DB, model *userLastOffsetModel) error {
	sq := "INSERT INTO user_last_offset (uid,channel_id,channel_type,message_seq) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE message_seq=IF(message_seq<VALUES(message_seq),VALUES(message_seq),message_seq)"
	err := tx.Exec(sq, model.UID, model.ChannelID, model.ChannelType, model.MessageSeq).Error
	return err
}

func (d *userLastOffsetDB) queryWithUID(uid string) ([]*userLastOffsetModel, error) {
	var ms []*userLastOffsetModel
	err := d.db.Table("user_last_offset").Where("uid=?", uid).Find(&ms).Error
	return ms, err
}

type userLastOffsetModel struct {
	UID         *string
	ChannelID   *string
	ChannelType *uint8
	MessageSeq  *int64
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
