package group

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

// DB DB
type settingDB struct {
	ctx *config.Context
	db  *gorm.DB
}

func newSettingDB(ctx *config.Context) *settingDB {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &settingDB{
		ctx: ctx,
		db:  db,
	}
}

// QuerySetting 查询设置
func (s *settingDB) QuerySetting(groupNo, uid string) (*Setting, error) {
	var m Setting
	err := s.db.Table("group_setting").Where("group_no=? and uid=?", groupNo, uid).First(&m).Error
	return &m, err
}

func (s *settingDB) querySettingWithTx(groupNo, uid string, tx *gorm.DB) (*Setting, error) {
	var m Setting
	err := tx.Table("group_setting").Where("group_no=? and uid=?", groupNo, uid).First(&m).Error
	return &m, err
}

func (s *settingDB) QuerySettings(groupNos []string, uid string) ([]*Setting, error) {
	var ms []*Setting
	err := s.db.Table("group_setting").Where("group_no in ? and uid=?", groupNos, uid).Find(&ms).Error
	return ms, err
}
func (s *settingDB) QuerySettingsWithUIDs(groupNo string, uids []string) ([]*Setting, error) {
	var ms []*Setting
	err := s.db.Table("group_setting").Where("group_no=? and uid in ?", groupNo, uids).Find(&ms).Error
	return ms, err
}

// InsertSetting 添加设置
func (s *settingDB) InsertSetting(setting *Setting) error {
	err := s.db.Table("group_setting").Create(setting).Error
	return err
}

// InsertSettingTx 添加设置
func (s *settingDB) InsertSettingTx(setting *Setting, tx *gorm.DB) error {
	err := tx.Table("group_setting").Create(setting).Error
	return err
}

// UpdateSetting 更新设置
func (s *settingDB) UpdateSetting(setting *Setting) error {
	err := s.db.Table("group_setting").Updates(map[string]interface{}{
		"chat_pwd_on":       setting.ChatPwdOn,
		"mute":              setting.Mute,
		"top":               setting.Top,
		"save":              setting.Save,
		"show_nick":         setting.ShowNick,
		"group_no":          setting.GroupNo,
		"uid":               setting.UID,
		"version":           setting.Version,
		"revoke_remind":     setting.RevokeRemind,
		"join_group_remind": setting.JoinGroupRemind,
		"screenshot":        setting.Screenshot,
		"receipt":           setting.Receipt,
		"flame":             setting.Flame,
		"flame_second":      setting.FlameSecond,
		"remark":            setting.Remark,
	}).Where("id=?", setting.Id).Error
	return err
}

// UpdateSetting 更新设置
func (s *settingDB) UpdateSettingWithTx(setting *Setting, tx *gorm.DB) error {
	err := tx.Table("group_setting").Updates(map[string]interface{}{
		"chat_pwd_on":       setting.ChatPwdOn,
		"mute":              setting.Mute,
		"top":               setting.Top,
		"save":              setting.Save,
		"show_nick":         setting.ShowNick,
		"group_no":          setting.GroupNo,
		"uid":               setting.UID,
		"version":           setting.Version,
		"revoke_remind":     setting.RevokeRemind,
		"join_group_remind": setting.JoinGroupRemind,
		"screenshot":        setting.Screenshot,
		"receipt":           setting.Receipt,
		"flame":             setting.Flame,
		"flame_second":      setting.FlameSecond,
		"remark":            setting.Remark,
	}).Where("id=?", setting.Id).Error
	return err
}

// Setting 群设置
type Setting struct {
	UID             *string // 用户uid
	GroupNo         *string // 群编号
	Mute            *int    // 免打扰
	Top             *int    // 置顶
	ShowNick        *int    // 显示昵称
	Save            *int    // 是否保存
	ChatPwdOn       *int    //是否开启聊天密码
	Screenshot      *int    //截屏通知
	RevokeRemind    *int    //撤回通知
	JoinGroupRemind *int    //进群提醒
	Receipt         *int    //消息是否回执
	Flame           *int    // 是否开启阅后即焚
	FlameSecond     *int    // 阅后即焚秒数
	Remark          *string // 群备注
	Version         *int64  // 版本
	Id              *int64
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func newDefaultSetting() *Setting {
	RevokeRemind := 1
	Screenshot := 1
	Receipt := 1
	return &Setting{
		RevokeRemind: &RevokeRemind,
		Screenshot:   &Screenshot,
		Receipt:      &Receipt,
	}
}
