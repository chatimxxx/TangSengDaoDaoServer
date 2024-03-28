package webhook

import (
	"embed"

	"github.com/xochat/xochat_im_server_lib/config"
	"github.com/xochat/xochat_im_server_lib/pkg/register"
)

//go:embed sql
var sqlFS embed.FS

func init() {

	register.AddModule(func(ctx interface{}) register.Module {
		xo := New(ctx.(*config.Context))
		return register.Module{
			SetupAPI: func() register.APIRouter {

				return xo
			},
			SQLDir: register.NewSQLFS(sqlFS),
			Start: func() error {
				return xo.Start()
			},
			Stop: func() error {
				return xo.Stop()
			},
		}
	})
}
