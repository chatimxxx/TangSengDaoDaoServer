package openapi

import (
	_ "embed"

	"github.com/xochat/xochat_im_server_lib/config"
	"github.com/xochat/xochat_im_server_lib/pkg/register"
)

//go:embed swagger/api.yaml
var swaggerContent string

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		x := ctx.(*config.Context)
		api := New(x)
		return register.Module{
			Name:    "openapi",
			Swagger: swaggerContent,
			SetupAPI: func() register.APIRouter {
				return api
			},
		}
	})
}
