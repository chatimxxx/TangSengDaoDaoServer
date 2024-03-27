package user

import (
	"gorm.io/gorm"
	"time"
)

// SettingDB 设置db
type SettingDB struct {
	db *gorm.DB
}

// NewSettingDB NewDB
func NewSettingDB(db *gorm.DB) *SettingDB {
	return &SettingDB{
		db: db,
	}
}

// InsertUserSettingModel 插入用户设置
func (d *SettingDB) InsertUserSettingModel(m *SettingModel) error {
	err := d.db.Table("user_setting").Create(m).Error
	return err
}

// InsertUserSettingModelTx 插入用户设置
func (d *SettingDB) InsertUserSettingModelTx(m *SettingModel, tx *gorm.DB) error {
	err := tx.Table("user_setting").Create(m).Error
	return err
}

// QueryUserSettingModel 查询用户设置
func (d *SettingDB) QueryUserSettingModel(uid, loginUID string) (*SettingModel, error) {
	var m SettingModel
	err := d.db.Table("user_setting").Where("uid=? and to_uid=?", loginUID, uid).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// QueryTwoUserSettingModel 查询双方用户设置
func (d *SettingDB) QueryTwoUserSettingModel(uid, loginUID string) ([]*SettingModel, error) {
	var ms []*SettingModel
	err := d.db.Table("user_setting").Where("(uid=? and to_uid=?) or (uid=? and to_uid=?)", loginUID, uid, uid, loginUID).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *SettingDB) QueryWithUidsAndToUID(uids []string, toUID string) ([]*SettingModel, error) {
	var ms []*SettingModel
	err := d.db.Table("user_setting").Where("uid in ? and to_uid=?", uids, toUID).Find(&ms).Error
	return ms, err
}

func (d *SettingDB) QueryUserSettings(uids []string, loginUID string) ([]*SettingModel, error) {
	var ms []*SettingModel
	err := d.db.Table("user_setting").Where("uid=? and to_uid in ?", loginUID, uids).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}

// updateUserSettingModel 更新用户设置
func (d *SettingDB) updateUserSettingModelWithToUIDTx(setting *SettingModel, uid string, toUID string, tx *gorm.DB) error {
	err := tx.Table("user_setting").Updates(map[string]interface{}{
		"mute":          setting.Mute,
		"top":           setting.Top,
		"blacklist":     setting.Blacklist,
		"chat_pwd_on":   setting.ChatPwdOn,
		"screenshot":    setting.Screenshot,
		"revoke_remind": setting.RevokeRemind,
		"receipt":       setting.Receipt,
		"flame":         setting.Flame,
		"flame_second":  setting.FlameSecond,
		"remark":        setting.Remark,
	}).Where("uid=? and to_uid=?", uid, toUID).Error
	return err
}

// UpdateUserSettingModel 更新用户设置
func (d *SettingDB) UpdateUserSettingModel(setting *SettingModel) error {
	err := d.db.Table("user_setting").Updates(map[string]interface{}{
		"mute":          setting.Mute,
		"top":           setting.Top,
		"version":       setting.Version,
		"chat_pwd_on":   setting.ChatPwdOn,
		"screenshot":    setting.Screenshot,
		"revoke_remind": setting.RevokeRemind,
		"receipt":       setting.Receipt,
		"flame":         setting.Flame,
		"flame_second":  setting.FlameSecond,
		"remark":        setting.Remark,
	}).Where("id=?", setting.Id).Error
	return err
}

func (d *SettingDB) querySettingByUIDAndToUID(uid, toUID string) (*SettingModel, error) {
	var m SettingModel
	err := d.db.Table("user_setting").Where("uid=? and to_uid=?", uid, toUID).First(&m).Error
	return &m, err
}

// ------------ model ------------

// SettingModel 用户设置
type SettingModel struct {
	UID          *string // 用户UID
	ToUID        *string // 对方uid
	Mute         *int    // 免打扰
	Top          *int    // 置顶
	ChatPwdOn    *int    // 是否开启聊天密码
	Screenshot   *int    //截屏通知
	RevokeRemind *int    //撤回提醒
	Blacklist    *int    //黑名单
	Receipt      *int    //消息是否回执
	Flame        *int    // 是否开启阅后即焚
	FlameSecond  *int    // 阅后即焚秒数
	Version      *int64  // 版本
	Remark       *string // 备注
	Id           *int64
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func newDefaultSettingModel() *SettingModel {
	Screenshot := 1
	RevokeRemind := 1
	Receipt := 1
	return &SettingModel{
		Screenshot:   &Screenshot,
		RevokeRemind: &RevokeRemind,
		Receipt:      &Receipt,
	}
}
