package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type managerDB struct {
	db  *gorm.DB
	ctx *config.Context
}

// newManagerDB
func newManagerDB(ctx *config.Context) *managerDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &managerDB{
		ctx: ctx,
		db:  d,
	}
}

// 通过账号和密码查询用户信息
func (d managerDB) queryUserInfoWithNameAndPwd(username string) (*managerLoginModel, error) {
	var m managerLoginModel
	err := d.db.Table("user").Where("username=?", username).First(&m).Error
	return &m, err
}

// 获取用户列表
func (d managerDB) queryUserListWithPage(pageSize, page int, onelineStatus int) ([]*managerUserModel, error) {
	// var ms []*managerUserModel
	// err := d.db.Table("user").Offset((page-1)*pageSize).Limit(pageSize).OrderDir("created_at", false).Find(&ms).Error
	// return ms, err

	var ms []*managerUserModel
	selectStm := d.db.Select("user.uid,user.name,user.username,user.status,user.phone,user.short_no,user.sex,user.is_destroy,user.created_at,user.gitee_uid,user.github_uid,user.wx_openid,max(user_online.online) online").Table("user").Joins("user_online on user.uid=user_online.uid")
	if onelineStatus != -1 {
		selectStm = selectStm.Where("user_online.online=?", onelineStatus)
	}
	selectStm = selectStm.Group("user.uid,user.name,user.username,user.status,user.phone,user.short_no,user.sex,user.is_destroy,user.created_at,user.gitee_uid,user.github_uid,user.wx_openid")

	// select  from user left join user_online on user.uid=user_online.uid where user_online.online=1  group by user.uid,user.name,user.status,user.phone,user.short_no,user.sex,user.is_destroy,user.created_at  limit 100
	err := selectStm.Offset((page - 1) * pageSize).Limit(pageSize).Order("user.created_at DESC").Find(&ms).Error
	return ms, err
}

// 模糊查询用户列表
// onelineStatus 在线状态 -1 为所有 0. 离线 1. 在线
func (d managerDB) queryUserListWithPageAndKeyword(keyword string, onelineStatus int, pageSize, page int) ([]*managerUserModel, error) {
	var ms []*managerUserModel
	selectStm := d.db.Select("user.uid,user.name,user.username,user.status,user.phone,user.short_no,user.sex,user.is_destroy,user.created_at,user.gitee_uid,user.github_uid,user.wx_openid,max(user_online.online) online").Table("user").Joins("user_online on user.uid=user_online.uid").Where("user.name like ? or user.uid like ? or user.phone like ? or user.short_no like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	if onelineStatus != -1 {
		selectStm = selectStm.Where("user_online.online=?", onelineStatus)
	}
	selectStm = selectStm.Group("user.uid,user.name,user.username,user.status,user.phone,user.short_no,user.sex,user.is_destroy,user.created_at,user.gitee_uid,user.github_uid,user.wx_openid")

	// select  from user left join user_online on user.uid=user_online.uid where user_online.online=1  group by user.uid,user.name,user.status,user.phone,user.short_no,user.sex,user.is_destroy,user.created_at  limit 100
	err := selectStm.Offset((page - 1) * pageSize).Limit(pageSize).Order("user.created_at DESC").Find(&ms).Error
	return ms, err
}

// 模糊查询用户数量
func (d managerDB) queryUserCountWithKeyWord(keyword string) (int64, error) {
	var count int64
	err := d.db.Table("user").Where("name like ? or uid like ? or phone like ? or short_no like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Count(&count).Error
	return count, err
}

// queryUserBlacklist 查询某个用户的黑名单
func (d managerDB) queryUserBlacklists(uid string) ([]*managerUserBlacklistModel, error) {
	var ms []*managerUserBlacklistModel
	err := d.db.Select("`user`.*,IFNULL(user_setting.updated_at,'') ").Table("`user`").Joins("`user_setting` on user.uid=user_setting.to_uid and user_setting.blacklist=1").Where("`user_setting`.uid=?", uid).Find(&ms).Error
	return ms, err
}

// 通过status查询用户列表
func (d managerDB) queryUserListWithStatus(status int, pageSize, page int) ([]*managerUserModel, error) {
	var ms []*managerUserModel
	err := d.db.Table("user").Where("status=?", status).Offset((page - 1) * pageSize).Limit(pageSize).Order("updated_at DESC").Find(&ms).Error
	return ms, err
}

// 通过status查询用户数量
func (d managerDB) queryUserCountWithStatus(status int) (int64, error) {
	var count int64
	err := d.db.Table("user").Where("status=?", status).Count(&count).Error
	return count, err
}

func (d managerDB) queryUserOnline(uid string) ([]*userOnline, error) {
	var ms []*userOnline
	err := d.db.Table("user_online").Where("uid=?", uid).Find(&ms).Error
	return ms, err
}

func (d managerDB) queryUserWithNameAndRole(username string, role string) (*managerUserModel, error) {
	var m managerUserModel
	err := d.db.Table("user").Where("username=? and role=?", username, role).First(&m).Error
	return &m, err
}

func (d managerDB) queryUsersWithRole(role string) ([]*managerUserModel, error) {
	var ms []*managerUserModel
	err := d.db.Table("user").Where("role=?", role).Find(&ms).Error
	return ms, err
}
func (d managerDB) deleteUserWithUIDAndRole(uid, role string) error {
	err := d.db.Table("user").Where("uid=? and role=?", uid, role).Delete(nil).Error
	return err
}

type managerLoginModel struct {
	Username string
	UID      string
	Name     string
	Password string
	Role     string
}

type managerUserModel struct {
	Username  *string
	Name      *string
	UID       *string
	Status    *int
	Phone     *string
	ShortNo   *string
	WXOpenid  *string // 微信openid
	GiteeUID  *string // gitee uid
	GithubUID *string // github uid
	Sex       *int
	IsDestroy *int
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type managerUserBlacklistModel struct {
	Name      *string
	UID       *string
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type userOnline struct {
	UID         *string
	DeviceFlag  *uint8 // 设备标记 0. APP 1.web
	LastOnline  *int   // 最后一次在线时间
	LastOffline *int   // 最后一次离线时间
	Online      *int
	Version     *int64 // 数据版本
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
