package base

import (
	"embed"

	"github.com/TangSengDaoDao/TangSengDaoDaoServer/modules/base/app"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/register"
)

//go:embed sql
var sqlFS embed.FS

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			SetupAPI: func() register.APIRouter {
				return app.New(ctx.(*config.Context))
			},
			SQLDir: register.NewSqlFS(sqlFS),
		}
	})
}
