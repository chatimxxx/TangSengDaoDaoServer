package workplace

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	dba "github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type managerDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newManagerDB(ctx *config.Context) *managerDB {
	return &managerDB{
		ctx: ctx,
		db:  ctx.DB(),
	}
}

func (d *managerDB) queryCateogryWithName(name string) (*categoryModel, error) {
	var m categoryModel
	err := d.db.Table("workplace_category").Where("name=?", name).First(&m).Error
	return &m, err
}

func (d *managerDB) insertCategory(m *categoryModel) error {
	err := d.db.Table("workplace_category").Create(m).Error
	return err
}
func (d *managerDB) queryCategoryWithNo(categoryNo string) (*categoryModel, error) {
	var m categoryModel
	err := d.db.Table("workplace_category").Where("category_no=?", categoryNo).First(&m).Error
	return &m, err
}
func (d *managerDB) queryMaxSortNumCategory() (*categoryModel, error) {
	var m categoryModel
	err := d.db.Table("workplace_category").Order("sort_num DESC").Limit(1).First(&m).Error
	return &m, err
}

func (d *managerDB) queryAppWithAppNameAndCategoryNo(appName string) (*appModel, error) {
	var m appModel
	err := d.db.Table("workplace_app").Where("name=?", appName).First(&m).Error
	return &m, err
}

func (d *managerDB) insertAPP(app *appModel) error {
	err := d.db.Table("workplace_app").Create(app).Error
	return err
}

func (d *managerDB) insertBanner(banner *bannerModel) error {
	err := d.db.Table("workplace_banner").Create(banner).Error
	return err
}

func (d *managerDB) updateApp(app *appModel) error {
	err := d.db.Table("workplace_app").Updates(map[string]interface{}{
		"app_category": app.AppCategory,
		"icon":         app.Icon,
		"name":         app.Name,
		"description":  app.Description,
		"status":       app.Status,
		"jump_type":    app.JumpType,
		"app_route":    app.AppRoute,
		"web_route":    app.WebRoute,
		"is_paid_app":  app.IsPaidApp,
	}).Where("app_id=?", app.AppID).Error
	return err
}
func (d *managerDB) updateBanner(banner *bannerModel) error {
	err := d.db.Table("workplace_banner").Updates(map[string]interface{}{
		"cover":       banner.Cover,
		"title":       banner.Title,
		"description": banner.Description,
		"jump_type":   banner.JumpType,
		"route":       banner.Route,
	}).Where("banner_no=?", banner.BannerNo).Error
	return err
}
func (d *managerDB) updateCategory(category *categoryModel) error {
	err := d.db.Table("workplace_category").Updates(map[string]interface{}{
		"name": category.Name,
	}).Where("category_no=?", category.CategoryNo).Error
	return err
}

func (d *managerDB) updateCategorySortNumWithTx(categoryNo string, sortNum int, tx *gorm.DB) error {
	err := tx.Table("workplace_category").Updates(map[string]interface{}{
		"sort_num": sortNum,
	}).Where("category_no=?", categoryNo).Error
	return err
}

func (d *managerDB) updateBannerSortNumWithTx(bannerNo string, sortNum int, tx *gorm.DB) error {
	err := tx.Table("workplace_banner").Updates(map[string]interface{}{
		"sort_num": sortNum,
	}).Where("banner_no=?", bannerNo).Error
	return err
}

func (d *managerDB) deleteAppTx(appId string, tx *gorm.DB) error {
	err := tx.Table("workplace_app").Where("app_id=?", appId).Delete(nil).Error
	return err
}
func (d *managerDB) deleteCategoryAppTx(appId string, tx *gorm.DB) error {
	err := tx.Table("workplace_category_app").Where("app_id=?", appId).Delete(nil).Error
	return err
}

func (d *managerDB) deleteUserAppTx(appId string, tx *gorm.DB) error {
	err := tx.Table("workplace_user_app").Where("app_id=?", appId).Delete(nil).Error
	return err
}

func (d *managerDB) deleteUserRecordAppTx(appId string, tx *gorm.DB) error {
	err := tx.Table("workplace_app_user_record").Where("app_id=?", appId).Delete(nil).Error
	return err
}

func (d *managerDB) queryCategory() ([]*categoryModel, error) {
	var ms []*categoryModel
	err := d.db.Table("workplace_category").Order("sort_num DESC").Find(&ms).Error
	return ms, err
}

func (d *managerDB) deleteBanner(bannerNo string) error {
	err := d.db.Table("workplace_banner").Where("banner_no=?", bannerNo).Error
	return err
}
func (d *managerDB) updateCategoryAppSortNumWithTx(categoryNo string, appId string, sortNum int, tx *gorm.DB) error {
	err := tx.Table("workplace_category_app").Updates(map[string]interface{}{
		"sort_num": sortNum,
	}).Where("category_no=? and app_id=?", categoryNo, appId).Error
	return err
}
func (d *managerDB) insertCategoryAppWithTx(m *categoryAppModel, tx *gorm.DB) error {
	err := tx.Table("workplace_category_app").Create(m).Error
	return err
}

func (d *managerDB) insertCategoryApp(m *categoryAppModel) error {
	err := d.db.Table("workplace_category_app").Create(m).Error
	return err
}

func (d *managerDB) deleteCategoryApp(appId, categoryNo string) error {
	err := d.db.Table("workplace_category_app").Where("app_id=? and category_no=?", appId, categoryNo).Delete(nil).Error
	return err
}
func (d *managerDB) deleteCategory(categoryNo string) error {
	err := d.db.Table("workplace_category").Where("category_no=?", categoryNo).Delete(nil).Error
	return err
}

func (d *managerDB) queryAppWithPage(pageSize, page int) ([]*appModel, error) {
	var ms []*appModel
	err := d.db.Table("workplace_app").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

func (d *managerDB) searchApp(keyword string, pageSize, page int) ([]*appModel, error) {
	var ms []*appModel
	err := d.db.Table("workplace_app").Where("name like ?", "%"+keyword+"%").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

// 通过关键字查询app总数
func (m *managerDB) queryAppCountWithKeyWord(keyword string) (int64, error) {
	var count int64
	err := m.db.Table("workplace_app").Where("name like ?", "%"+keyword+"%").Count(&count).Error
	return count, err
}

// 查询app总数
func (d *managerDB) queryAppCount() (int64, error) {
	var count int64
	err := d.db.Table("workplace_app").Count(&count).Error
	return count, err
}

func (d *managerDB) queryMaxSortNumCategoryApp(categoryNo string) (*categoryAppModel, error) {
	var m categoryAppModel
	err := d.db.Table("workplace_category_app").Where("category_no=?", categoryNo).Order("sort_num DESC").Limit(1).First(&m).Error
	return &m, err
}

type categoryAppModel struct {
	CategoryNo string
	AppId      string
	SortNum    int // 排序号
	dba.BaseModel
}
