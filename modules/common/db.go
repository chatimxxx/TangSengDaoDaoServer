package common

import (
	dbs "github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func newDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// 添加版本升级
func (d *DB) insertAppVersion(m *appVersionModel) (int64, error) {
	err := d.db.Table("app_version").Create(m).Error
	if err != nil {
		return 0, err
	}
	return *m.Id, err
}

// 查询某个系统的最新版本
func (d *DB) queryNewVersion(os string) (*appVersionModel, error) {
	var model *appVersionModel
	err := d.db.Table("app_version").Where("os=?", os).Order("created_at").Limit(1).First(&model).Error
	return model, err
}

// 查询版本升级列表
func (d *DB) queryAppVersionListWithPage(pageSize, page int) ([]*appVersionModel, error) {
	var models []*appVersionModel
	err := d.db.Table("app_version").Offset((page - 1) * pageSize).Limit(pageSize).Order("updated_at").Find(&models).Error
	return models, err
}

// 模糊查询用户数量
func (d *DB) queryCount() (int64, error) {
	var count int64
	err := d.db.Table("app_version").Count(&count).Error
	return count, err
}

// 查询所有背景图片
func (d *DB) queryChatBgs() ([]*chatBgModel, error) {
	var models []*chatBgModel
	err := d.db.Table("chat_bg").Find(&models).Error
	return models, err
}

// 查询app模块
func (d *DB) queryAppModule() ([]*appModuleModel, error) {
	var list []*appModuleModel
	err := d.db.Table("app_module").Order("created_at ASC").Find(&list).Error
	return list, err
}

// 查询某个app模块
func (d *DB) queryAppModuleWithSid(sid string) (*appModuleModel, error) {
	var m appModuleModel
	err := d.db.Table("app_module").Where("sid=?", sid).First(&m).Error
	return &m, err
}

// 新增app模块
func (d *DB) insertAppModule(m *appModuleModel) (int64, error) {
	err := d.db.Table("app_module").Create(m).Error
	if err != nil {
		return 0, err
	}
	return *m.Id, err
}

// 修改app模块
func (d *DB) updateAppModule(m *appModuleModel) error {
	err := d.db.Table("app_module").Updates(&m).Where("id=?", m.Id).Error
	return err
}

// 删除模块
func (d *DB) deleteAppModule(sid string) error {
	err := d.db.Table("app_module").Where("sid=?", sid).Delete(nil).Error
	return err
}

type chatBgModel struct {
	Cover string // 封面
	Url   string // 图片地址
	IsSvg int    // 1 svg图片 0 普通图片
	dbs.BaseModel
}

type appVersionModel struct {
	AppVersion  string // app版本
	OS          string // android | ios
	IsForce     int    // 是否强制更新 1:是
	UpdateDesc  string // 更新说明
	DownloadURL string // 下载地址
	Signature   string // 安装包签名
	dbs.BaseModel
}

type appModuleModel struct {
	SID    string // 模块ID
	Name   string // 模块名称
	Desc   string // 介绍
	Status int    // 状态
	dbs.BaseModel
}
