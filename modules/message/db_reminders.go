package message

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
)

type remindersDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newRemindersDB(ctx *config.Context) *remindersDB {
	db, err := ctx.DB()
	if err != nil {
		panic("服务初始化失败")
		return nil
	}
	return &remindersDB{
		ctx: ctx,
		db:  db,
	}
}

func (r *remindersDB) inserts(models []*remindersModel) error {
	tx := r.db.Begin()
	for _, m := range models {
		err := tx.Table("reminders").Create(m).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *remindersDB) deleteWithChannel(channelID string, channelType uint8, messageID int64, version int64) error {
	err := r.db.Table("reminders").Update("is_deleted", 1).Update("version", version).Where("channel_id=? and channel_type=? and message_id=?", channelID, channelType, messageID).Error
	return err
}

func (r *remindersDB) deleteWithChannelAndUIDTx(channelID string, channelType uint8, uid string, messageID int64, version int64, tx *gorm.DB) error {
	err := tx.Table("reminders").Update("is_deleted", 1).Update("version", version).Where("channel_id=? and channel_type=? and uid=? and message_id=?", channelID, channelType, uid, messageID).Error
	return err
}

/*
*
同步提醒项
@param uid 当前登录用户的uid
@param version 以uid为key的增量版本号
@param limit 数据限制
@param channelIDs 频道集合 查询以频道为目标的提醒项
*
*/
func (r *remindersDB) sync(uid string, version int64, limit int, channelIDs []string) ([]*remindersDetailModel, error) {
	var ms []*remindersDetailModel
	var err error
	if version == 0 {
		builder := r.db.Select("reminders.*,IF(reminder_done.id is null and reminders.is_deleted=0,0,1) done").Table("reminders").Joins(fmt.Sprintf("reminder_done on reminders.id=reminder_done.reminder_id and reminder_done.uid='%s'", uid))

		if len(channelIDs) == 0 {
			err = builder.Where("(reminders.uid=?  or   reminders.uid='')  and reminders.version>? and reminder_done.id is null", uid, version).Order("version ASC").Limit(limit).Find(&ms).Error
		} else {
			err = builder.Where("(reminders.uid=?  or  ( reminders.uid='' and reminders.channel_id in ?))  and reminders.version>? and reminder_done.id is null", uid, channelIDs, version).Order("version ASC").Limit(limit).Find(&ms).Error
		}
	} else {
		build := r.db.Select("reminders.*,IF(reminder_done.id is null and reminders.is_deleted=0,0,1) done").Table("reminders").Joins(fmt.Sprintf("reminder_done on reminders.id=reminder_done.reminder_id and reminder_done.uid='%s'", uid))
		if len(channelIDs) == 0 {
			err = build.Where("(reminders.uid=?  or  reminders.uid='')  and reminders.version>?", uid, version).Order("version ASC").Limit(limit).Find(&ms).Error
		} else {
			err = build.Where("(reminders.uid=?  or  ( reminders.uid='' and reminders.channel_id in ?))  and reminders.version>?", uid, channelIDs, version).Order("version ASC").Limit(limit).Find(&ms).Error
		}

	}
	return ms, err
}

func (r *remindersDB) insertDonesTx(ids []int64, uid string, tx *gorm.DB) error {
	for _, id := range ids {
		err := tx.Exec("insert ignore  into reminder_done(reminder_id,uid) values(?,?)", id, uid).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *remindersDB) updateVersionTx(version int64, id int64, tx *gorm.DB) error {
	err := tx.Table("reminders").Update("version", version).Where("id=?", id).Error
	return err
}

type remindersDetailModel struct {
	Done int
	remindersModel
}

type remindersModel struct {
	ChannelID    *string
	ChannelType  *uint8
	ClientMsgNo  *string
	MessageSeq   *uint32
	MessageID    *string
	ReminderType *int
	Publisher    *string
	UID          *string
	Text         *string
	Data         *string
	IsLocate     *int
	Version      *int64
	IsDeleted    *int
	Id           *int64
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type reminderDoneModel struct {
	ReminderID *int64
	UID        *string
	Id         *int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
