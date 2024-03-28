package file

import (
	_ "embed"

	"github.com/xochat/xochat_im_server_lib/config"
	"github.com/xochat/xochat_im_server_lib/pkg/register"
)

//go:embed swagger/api.yaml
var swaggerContent string

func init() {

	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			Name: "file",
			SetupAPI: func() register.APIRouter {
				return New(ctx.(*config.Context))
			},
			Swagger: swaggerContent,
		}
	})
}
