package base

import (
	"embed"

	"github.com/xochat/xochat_im_server/modules/base/app"
	"github.com/xochat/xochat_im_server_lib/config"
	"github.com/xochat/xochat_im_server_lib/pkg/register"
)

//go:embed sql
var sqlFS embed.FS

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			SetupAPI: func() register.APIRouter {
				return app.New(ctx.(*config.Context))
			},
			SQLDir: register.NewSQLFS(sqlFS),
		}
	})
}
