package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type deviceFlagDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newDeviceFlagDB(ctx *config.Context) *deviceFlagDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &deviceFlagDB{
		db:  d,
		ctx: ctx,
	}
}

func (d *deviceFlagDB) queryAll() ([]*deviceFlagModel, error) {
	var ms []*deviceFlagModel
	err := d.db.Table("device_flag").Find(&ms).Error
	return ms, err
}

type deviceFlagModel struct {
	DeviceFlag *uint8
	Weight     *int
	Remark     *string
	Id         *int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
