package kids

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	kidsv1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	"rslytics-app-api/internal/common/response"
	"rslytics-app-api/internal/controller/kids"
	"rslytics-app-api/internal/middleware"
	serverconfig "rslytics-app-api/internal/server/config"
)

var Main = gcmd.Command{
	Name:  "kids",
	Usage: "kids [-e dev|test|prod]",
	Brief: "start microservice for kids application",
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
		serverconfig.SetFile("kids", env)
		commoni18n.Init()

		s := g.Server("kids")
		s.Use(middleware.Ctx, middleware.I18n, middleware.KidsJWT, middleware.V1Envelope, response.Middleware)
		if openAPIPath := s.GetOpenApiPath(); openAPIPath != "" {
			s.BindHookHandler(openAPIPath, ghttp.HookBeforeServe, func(r *ghttp.Request) {
				if r.Method != http.MethodGet {
					return
				}
				r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Response.Write(kidsv1.OpenAPISpec())
				r.ExitAll()
			})
		}
		s.Group("/", func(group *ghttp.RouterGroup) {
			kids.Register(ctx, group)
		})
		s.Run()
		return nil
	},
}
