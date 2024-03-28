package qrcode

import (
	"github.com/xochat/xochat_im_server_lib/config"
	"github.com/xochat/xochat_im_server_lib/pkg/register"
)

func init() {

	register.AddModule(func(ctx interface{}) register.Module {

		return register.Module{
			SetupAPI: func() register.APIRouter {
				return New(ctx.(*config.Context))
			},
		}
	})
}
