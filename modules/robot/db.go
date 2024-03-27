package robot

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type robotDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newBotDB(ctx *config.Context) *robotDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &robotDB{
		ctx: ctx,
		db:  d,
	}
}
func (d *robotDB) queryRobotWithRobtID(robotID string) (*robot, error) {
	var m robot
	err := d.db.Table("robot").Where("robot_id=?", robotID).First(&m).Error
	return &m, err
}
func (d *robotDB) queryVaildRobotWithRobtID(robotID string) (*robot, error) {
	var m robot
	err := d.db.Table("robot").Where("robot_id=? and status=1", robotID).First(&m).Error
	return &m, err
}

func (d *robotDB) exist(robotID string) (bool, error) {
	var cn int64
	err := d.db.Table("robot").Where("robot_id=? and status=1", robotID).Count(&cn).Error
	return cn > 0, err
}

func (d *robotDB) insert(m *robot) error {
	err := d.db.Table("robot").Create(m).Error
	return err
}

func (d *robotDB) insertTx(m *robot, tx *gorm.DB) error {
	err := tx.Table("robot").Create(m).Error
	return err
}
func (d *robotDB) insertMenuTx(m *menu, tx *gorm.DB) error {
	err := tx.Table("robot_menu").Create(m).Error
	return err
}

func (d *robotDB) queryWithIDs(robotIDs []string) ([]*robot, error) {
	var ms []*robot
	err := d.db.Table("robot").Where("robot_id in ?", robotIDs).Find(&ms).Error
	return ms, err
}
func (d *robotDB) queryWithUsernames(usernames []string) ([]*robot, error) {
	var ms []*robot
	err := d.db.Table("robot").Where("username in ?", usernames).Find(&ms).Error
	return ms, err
}
func (d *robotDB) queryWithUsername(username string) (*robot, error) {
	var rb robot
	err := d.db.Table("robot").Where("username = ?", username).First(&rb).Error
	return &rb, err
}

func (d *robotDB) queryVaildRobotIDs(robotIDs []string) ([]string, error) {
	var ms []string
	err := d.db.Select("robot_id").Table("robot").Where("robot_id in ?", robotIDs).Find(&ms).Error
	return ms, err
}

// 同步机器人菜单
func (d *robotDB) queryMenusWithRobotIDs(uids []string) ([]*menu, error) {
	var ms []*menu
	err := d.db.Table("robot_menu").Where("robot_id in ?", uids).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

// 修改机器人信息
func (d *robotDB) updateRobotTx(m *robot, tx *gorm.DB) error {
	err := tx.Table("robot").Updates(map[string]interface{}{
		"version": m.Version,
	}).Where("robot_id=?", m.RobotID).Error
	return err
}
func (d *robotDB) updateRobot(m *robot) error {
	err := d.db.Table("robot").Updates(map[string]interface{}{
		"version": m.Version,
		"status":  m.Status,
	}).Where("robot_id=?", m.RobotID).Error
	return err
}
func (d *robotDB) queryMenusWithRobotID(robotID string) ([]*menu, error) {
	var ms []*menu
	err := d.db.Table("robot_menu").Where("robot_id=?", robotID).Order("created_at DESC").Find(&ms).Error
	return ms, err
}
func (d *robotDB) deleteMenuWithID(robotID string, id int64, tx *gorm.DB) error {
	err := tx.Table("robot_menu").Where("robot_id=? and id=?", robotID, id).Delete(nil).Error
	return err
}

type menu struct {
	RobotID   *string // 机器人ID
	CMD       *string // 命令
	Remark    *string // 命令说明
	Type      *string // 命令类型
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
type robot struct {
	AppID       *string
	RobotID     *string // 机器人唯一ID
	Username    *string // 机器人用户名
	InlineOn    *int    // 是否开启行内搜索
	Placeholder *string // 输入框占位符，开启行内搜索有效
	Token       *string
	Version     *int64
	Status      *int
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
