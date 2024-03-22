package webhook

import (
	"embed"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/register"
)

//go:embed sql
var sqlFS embed.FS

func init() {

	register.AddModule(func(ctx interface{}) register.Module {
		wk := New(ctx.(*config.Context))
		return register.Module{
			SetupAPI: func() register.APIRouter {

				return wk
			},
			SQLDir: register.NewSQLFS(sqlFS),
			Start: func() error {
				return wk.Start()
			},
			Stop: func() error {
				return wk.Stop()
			},
		}
	})
}
