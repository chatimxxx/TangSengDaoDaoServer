package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/common"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/log"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/wkhttp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Setting 用户设置
type Setting struct {
	ctx *config.Context
	log.Log
	db *SettingDB
}

// NewSetting 创建
func NewSetting(ctx *config.Context) *Setting {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &Setting{ctx: ctx, Log: log.NewTLog("UserSetting"), db: NewSettingDB(db)}
}

// 用户设置
func (u *Setting) userSettingUpdate(c *wkhttp.Context) {
	loginUID := c.MustGet("uid").(string)
	toUID := c.Param("uid")
	var settingMap map[string]interface{}
	if err := c.BindJSON(&settingMap); err != nil {
		u.Error("数据格式有误！", zap.Error(err))
		c.ResponseError(errors.New("数据格式有误！"))
		return
	}
	model, err := u.db.QueryUserSettingModel(toUID, loginUID)
	if err != nil {
		u.Error("查询用户设置失败！", zap.Error(err))
		c.ResponseError(errors.New("查询用户设置失败！"))
		return
	}
	insert := false // 是否是插入操作
	if model == nil {
		insert = true // 是否是插入操作
		model = newDefaultSettingModel()
		model.UID = &loginUID
		model.ToUID = &toUID
	}
	for key, value := range settingMap {
		switch key {
		case "mute":
			mute := int(value.(float64))
			model.Mute = &mute
		case "top":
			top := int(value.(float64))
			model.Top = &top
		case "chat_pwd_on":
			chatPwdOn := int(value.(float64))
			model.ChatPwdOn = &chatPwdOn
		case "screenshot":
			screenshot := int(value.(float64))
			model.Screenshot = &screenshot
		case "revoke_remind":
			revokeRemind := int(value.(float64))
			model.RevokeRemind = &revokeRemind
		case "receipt":
			receipt := int(value.(float64))
			model.Receipt = &receipt
		case "flame":
			flame := int(value.(float64))
			model.Flame = &flame
		case "flame_second":
			flameSecond := int(value.(float64))
			model.FlameSecond = &flameSecond
		case "remark":
			remark := value.(string)
			model.Remark = &remark
		}
	}
	version, _ := u.ctx.GenSeq(common.UserSettingSeqKey)
	model.Version = &version
	if insert {
		err = u.db.InsertUserSettingModel(model)
		if err != nil {
			u.Error("添加设置失败！", zap.Error(err))
			c.ResponseError(errors.New("添加设置失败！"))
			return
		}
	} else {
		err = u.db.UpdateUserSettingModel(model)
		if err != nil {
			u.Error("修改设置失败！", zap.Error(err))
			c.ResponseError(errors.New("修改设置失败！"))
			return
		}
	}
	// 发送一个频道更新命令 发给自己的其他设备，如果其他设备在线的话
	err = u.ctx.SendCMD(config.MsgCMDReq{
		ChannelID:   loginUID,
		ChannelType: common.ChannelTypePerson.Uint8(),
		CMD:         common.CMDChannelUpdate,
		Param: map[string]interface{}{
			"channel_id":   toUID,
			"channel_type": common.ChannelTypePerson,
		},
	})
	if err != nil {
		u.Error("发送频道更新命令失败！", zap.Error(err))
		c.ResponseError(errors.New("发送频道更新命令失败！"))
		return
	}
	c.ResponseOK()
}
