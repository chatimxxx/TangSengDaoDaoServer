package user

import (
	"github.com/gocraft/dbr/v2"
	"github.com/xochat/xochat_im_server_lib/config"
	"github.com/xochat/xochat_im_server_lib/pkg/db"
)

type deviceFlagDB struct {
	session *dbr.Session
	ctx     *config.Context
}

func newDeviceFlagDB(ctx *config.Context) *deviceFlagDB {
	return &deviceFlagDB{
		session: ctx.DB(),
		ctx:     ctx,
	}
}

func (d *deviceFlagDB) queryAll() ([]*deviceFlagModel, error) {
	var deviceFlags []*deviceFlagModel
	_, err := d.session.Select("*").From("device_flag").Load(&deviceFlags)
	return deviceFlags, err
}

type deviceFlagModel struct {
	DeviceFlag uint8
	Weight     int
	Remark     string
	db.BaseModel
}
