package report

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	dba "github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type db struct {
	db  *gorm.DB
	ctx *config.Context
}

func newDB(ctx *config.Context) *db {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &db{
		ctx: ctx,
		db:  d,
	}
}

func (d *db) queryCategoryAll() ([]*categoryModel, error) {
	var ms []*categoryModel
	err := d.db.Table("report_category").Find(&ms).Error
	return ms, err
}

func (d *db) insertCategory(m *categoryModel) error {
	err := d.db.Table("report_category").Create(m).Error
	return err
}

func (d *db) insert(m *model) error {
	err := d.db.Table("report").Create(m).Error
	return err
}

type categoryModel struct {
	CategoryNo       string
	CategoryName     string
	CategoryEname    string // 英文分类名称
	ParentCategoryNo string
	dba.BaseModel
}

type model struct {
	UID         string
	CategoryNo  string
	ChannelID   string
	ChannelType uint8
	Imgs        string
	Remark      string
	dba.BaseModel
}
