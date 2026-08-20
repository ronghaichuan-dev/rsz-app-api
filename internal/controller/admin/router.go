package controller

import (
	"context"
	kids2 "rslytics-app-api/internal/controller/admin/kids"
	"rslytics-app-api/internal/controller/admin/system"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Register(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Bind(
			kids2.New(),
			system.New(),
		)
	})
}
