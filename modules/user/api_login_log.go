package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/log"
	"go.uber.org/zap"
)

// LoginLog 用户设置
type LoginLog struct {
	ctx *config.Context
	log.Log
	loginLogDB *LoginLogDB
}

// NewLoginLog 创建
func NewLoginLog(ctx *config.Context) *LoginLog {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &LoginLog{ctx: ctx, Log: log.NewTLog("loginLog"), loginLogDB: NewLoginLogDB(db)}
}

// add 添加登录日志
func (l *LoginLog) add(uid string, publicIP string) {
	err := l.loginLogDB.insert(&LoginLogModel{
		UID:     &uid,
		LoginIP: &publicIP,
	})
	if err != nil {
		l.Error("添加登录日志错误", zap.Error(err))
	}
}

// getLastLoginIp 获取最后一次登录ip
func (l *LoginLog) getLastLoginIP(uid string) *loginLogResp {
	model, err := l.loginLogDB.queryLastLoginIP(uid)
	if err != nil {
		l.Error("查询登录日志错误", zap.Error(err))
		return nil
	}
	if model != nil {
		return &loginLogResp{
			UID:      *model.UID,
			CreateAt: model.CreatedAt.String(),
			LoginIP:  *model.LoginIP,
		}
	}
	return nil
}

// loginLogResp 登录日志
type loginLogResp struct {
	UID      string
	CreateAt string
	LoginIP  string
}
