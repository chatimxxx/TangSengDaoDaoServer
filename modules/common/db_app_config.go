package common

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type appConfigDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newAppConfigDB(ctx *config.Context) *appConfigDB {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &appConfigDB{
		db:  db,
		ctx: ctx,
	}
}

func (a *appConfigDB) query() (*appConfigModel, error) {
	var m appConfigModel
	err := a.db.Table("app_config").Order("created_at DESC").Find(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, err
}

func (a *appConfigDB) insert(m *appConfigModel) error {
	err := a.db.Table("app_config").Create(m).Error
	return err
}
func (a *appConfigDB) updateWithMap(RevokeSecond int, WelcomeMessage string, NewUserJoinSystemGroup int, SearchByPhone int, id int64) error {
	m := appConfigModel{
		RevokeSecond:           &RevokeSecond,
		WelcomeMessage:         &WelcomeMessage,
		NewUserJoinSystemGroup: &NewUserJoinSystemGroup,
		SearchByPhone:          &SearchByPhone,
	}
	err := a.db.Table("app_config").Updates(&m).Where("id=?", id).Error
	return err
}

type appConfigModel struct {
	RSAPrivateKey          *string
	RSAPublicKey           *string
	Version                *int
	SuperToken             *string
	SuperTokenOn           *int
	RevokeSecond           *int    // 消息可撤回时长
	WelcomeMessage         *string // 登录欢迎语
	NewUserJoinSystemGroup *int    // 新用户是否加入系统群聊
	SearchByPhone          *int    // 是否可通过手机号搜索
	Id                     *int64
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
}
