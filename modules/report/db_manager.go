package report

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	dba "github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type managerDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newManagerDB(ctx *config.Context) *managerDB {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &managerDB{
		ctx: ctx,
		db:  db,
	}
}

// 查询举报列表
func (m *managerDB) list(pageSize, page int, channelType int) ([]*managerReportModel, error) {
	var ms []*managerReportModel
	err := m.db.Select("report.*,report_category.category_name").Table("report").Joins("report_category on report.category_no=report_category.category_no").Where("report.channel_type=?", channelType).Offset((page - 1) * pageSize).Limit(pageSize).Order("report.created_at DESC").Find(&ms).Error
	return ms, err
}

// 查询总用户
func (m *managerDB) queryReportCount(channelType int) (int64, error) {
	var count int64
	err := m.db.Table("report").Where("channel_type=?", channelType).Count(&count).Error
	return count, err
}

type managerReportModel struct {
	UID          string
	CategoryNo   string
	ChannelID    string
	ChannelType  uint8
	Imgs         string
	Remark       string
	CategoryName string
	dba.BaseModel
}
