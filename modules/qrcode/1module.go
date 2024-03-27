package qrcode

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/register"
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
