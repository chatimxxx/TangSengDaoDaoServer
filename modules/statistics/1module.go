package statistics

import (
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/register"
)

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		x := ctx.(*config.Context)
		return register.Module{
			Name: "statistics",
			SetupAPI: func() register.APIRouter {
				return NewStatistics(x)
			},
		}
	})
}
