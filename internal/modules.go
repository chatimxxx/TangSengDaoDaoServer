package modules

// 引入模块
import (
	_ "github.com/xochat/xochat_im_server/modules/base"
	_ "github.com/xochat/xochat_im_server/modules/channel"
	_ "github.com/xochat/xochat_im_server/modules/common"
	_ "github.com/xochat/xochat_im_server/modules/file"
	_ "github.com/xochat/xochat_im_server/modules/group"
	_ "github.com/xochat/xochat_im_server/modules/message"
	_ "github.com/xochat/xochat_im_server/modules/openapi"
	_ "github.com/xochat/xochat_im_server/modules/qrcode"
	_ "github.com/xochat/xochat_im_server/modules/report"
	_ "github.com/xochat/xochat_im_server/modules/robot"
	_ "github.com/xochat/xochat_im_server/modules/statistics"
	_ "github.com/xochat/xochat_im_server/modules/user"
	_ "github.com/xochat/xochat_im_server/modules/webhook"
)
