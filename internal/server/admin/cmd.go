package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"rslytics-app-api/internal/common/response"
	adminctrl "rslytics-app-api/internal/controller/admin"
	"rslytics-app-api/internal/middleware"
	serverconfig "rslytics-app-api/internal/server/config"
)

var Main = gcmd.Command{
	Name:  "admin",
	Usage: "admin [-e dev|test|prod]",
	Brief: "start microservice for admin application",
	Arguments: []gcmd.Argument{
		{
			Name:    "env",
			Short:   "e",
			Default: serverconfig.DefaultEnv,
			Brief:   "runtime environment: dev, test, prod",
		},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		env := serverconfig.ResolveEnv(parser)
		serverconfig.SetFile("admin", env)

		s := g.Server("admin")
		s.Use(middleware.Ctx, response.Middleware)
		s.Group("/", func(group *ghttp.RouterGroup) {
			adminctrl.Register(ctx, group)
		})
		s.Run()
		return nil
	},
}
