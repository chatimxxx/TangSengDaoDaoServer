package workplace

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	dba "github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type db struct {
	db  *gorm.DB
	ctx *config.Context
}

func newDB(ctx *config.Context) *db {
	return &db{
		ctx: ctx,
		db:  ctx.DB(),
	}
}
func (d *db) updateUserAppSortNumWithTx(uid, appId string, sortNum int, tx *gorm.DB) error {
	err := tx.Table("workplace_user_app").Updates(map[string]interface{}{
		"sort_num": sortNum,
	}).Where("uid=? and app_id=?", uid, appId).Error
	return err
}

func (d *db) insertUserApp(app *userAppModel) error {
	err := d.db.Table("workplace_user_app").Create(app).Error
	return err
}

func (d *db) queryUserAppMaxSortNumWithUID(uid string) (*userAppModel, error) {
	var m userAppModel
	err := d.db.Table("workplace_user_app").Where("uid=?", uid).Order("sort_num DESC").Limit(1).Find(&m).Error
	return &m, err
}

func (d *db) deleteUserAppWithAppId(uid, appId string) error {
	err := d.db.Table("workplace_user_app").Where("app_id=? and uid=?", appId, uid).Delete(nil).Error
	return err
}

func (d *db) queryCategory() ([]*categoryModel, error) {
	var ms []*categoryModel
	err := d.db.Table("workplace_category").Order("sort_num DESC").Find(&ms).Error
	return ms, err
}

func (d *db) queryAppWithAppIds(ids []string) ([]*appModel, error) {
	var ms []*appModel
	err := d.db.Table("workplace_app").Where("app_id in ?", ids).Find(&ms).Error
	return ms, err
}

func (d *db) queryAppWithUid(uid string) ([]*appModel, error) {
	var ms []*appModel
	err := d.db.Table("workplace_user_app").Where("uid=?", uid).Find(&ms).Error
	return ms, err
}

func (d *db) queryAppWithAppId(appId string) (*appModel, error) {
	var app appModel
	err := d.db.Table("workplace_app").Where("app_id=?", appId).First(&app).Error
	return &app, err
}

func (d *db) queryAppWithCategroyNo(categoryNo string) ([]*cAppModel, error) {
	var ms []*cAppModel
	err := d.db.Select("workplace_category_app.sort_num,workplace_app.app_id,workplace_app.icon,workplace_app.name,workplace_app.description,workplace_app.app_category,workplace_app.jump_type,workplace_app.status,workplace_app.app_route,workplace_app.web_route,workplace_app.is_paid_app,workplace_app.created_at").Table("workplace_category_app").Joins("workplace_app on workplace_category_app.app_id=workplace_app.app_id").Where("workplace_category_app.category_no=?", categoryNo).Order("workplace_category_app.sort_num DESC").Find(&ms).Error
	return ms, err
}

func (d *db) queryUserAppWithAPPId(uid string, appId string) (*userAppModel, error) {
	var app userAppModel
	err := d.db.Table("workplace_user_app").Where("uid=? and app_id=?", uid, appId).First(&app).Error
	return &app, err
}

func (d *db) queryUserApp(uid string) ([]*userAppModel, error) {
	var ms []*userAppModel
	err := d.db.Table("workplace_user_app").Where("uid=?", uid).Order("sort_num DESC").Find(&ms).Error
	return ms, err
}

func (d *db) queryBanner() ([]*bannerModel, error) {
	var ms []*bannerModel
	err := d.db.Table("workplace_banner").Order("sort_num DESC").Find(&ms).Error
	return ms, err
}

func (d *db) insertRecord(record *recordModel) error {
	err := d.db.Table("workplace_app_user_record").Create(record).Error
	return err
}

func (d *db) queryRecordWithUid(uid string) ([]*recordModel, error) {
	var ms []*recordModel
	err := d.db.Table("workplace_app_user_record").Where("uid=?", uid).Order("count DESC").Find(&ms).Error
	return ms, err
}

func (d *db) queryRecordWithUidAndAppId(uid, appId string) (*recordModel, error) {
	var record recordModel
	err := d.db.Table("workplace_app_user_record").Where("uid=? and app_id=?", uid, appId).First(&record).Error
	return &record, err
}

func (d *db) updateRecordCount(record *recordModel) error {
	err := d.db.Table("workplace_app_user_record").Updates(map[string]interface{}{
		"count": record.Count,
	}).Where("uid=? and app_id=?", record.Uid, record.AppId).Error
	return err
}
func (d *db) deleteRecord(uid, appId string) error {
	err := d.db.Table("workplace_app_user_record").Where("app_id=? and uid=?", appId, uid).Delete(nil).Error
	return err
}

type recordModel struct {
	Count int // 使用次数
	Uid   string
	AppId string
	dba.BaseModel
}
type categoryModel struct {
	CategoryNo string //  分类编号
	Name       string // 分类名称
	SortNum    int    //  排序编号
	dba.BaseModel
}

type bannerModel struct {
	BannerNo    string // 封面编号
	Cover       string // 封面
	Title       string // 标题
	Description string // 介绍
	JumpType    int    // 打开方式 0.网页 1.原生
	Route       string // 打开地址
	SortNum     int    //  排序编号
	dba.BaseModel
}
type userAppModel struct {
	AppID   string // 分类项唯一id
	SortNum int    // 排序编号
	Uid     string // 所属用户uid
	dba.BaseModel
}

type cAppModel struct {
	SortNum int // 排序编号
	appModel
}

type appModel struct {
	AppID       string // 应用ID
	Icon        string // 应用icon
	Name        string // 应用名称
	Description string // 应用介绍
	AppCategory string // 应用分类 [‘机器人’ ‘客服’]
	Status      int    // 是否可用 0.禁用 1.可用
	JumpType    int    // 打开方式 0.网页 1.原生
	AppRoute    string // app打开地址
	WebRoute    string // web打开地址
	IsPaidApp   int    // 是否为付费应用 0.否 1.是
	dba.BaseModel
}
