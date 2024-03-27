package message

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type deviceOffsetDB struct {
	db *gorm.DB
}

func newDeviceOffsetDB(db *gorm.DB) *deviceOffsetDB {
	return &deviceOffsetDB{
		db: db,
	}
}

func (d *deviceOffsetDB) insertOrUpdateTx(tx *gorm.DB, model *deviceOffsetModel) error {
	sq := fmt.Sprintf("INSERT INTO device_offset (uid,device_uuid,channel_id,channel_type,message_seq) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE message_seq=IF(message_seq<VALUES(message_seq),VALUES(message_seq),message_seq)")
	err := tx.Exec(sq, model.UID, model.DeviceUUID, model.ChannelID, model.ChannelType, model.MessageSeq).Error
	return err
}

func (d *deviceOffsetDB) queryWithUIDAndDeviceUUID(uid string, deviceUUID string) ([]*deviceOffsetModel, error) {
	var ms []*deviceOffsetModel
	err := d.db.Table("device_offset").Where("uid=? and device_uuid=?", uid, deviceUUID).Find(&ms).Error
	return ms, err
}

func (d *deviceOffsetDB) queryMessageSeq(uid string, deviceUUID string, channelID string, channelType uint8) (int64, error) {
	var messageSeq int64
	err := d.db.Select("IFNULL(message_seq,0)").Table("device_offset").Where("uid=? and device_uuid=? and channel_id=? and channel_type=?", uid, deviceUUID, channelID, channelType).Limit(1).First(&messageSeq).Error
	return messageSeq, err
}

type deviceOffsetModel struct {
	UID         *string
	DeviceUUID  *string
	ChannelID   *string
	ChannelType *uint8
	MessageSeq  *int64
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
