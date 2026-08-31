package kids

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"rslytics-app-api/internal/controller/kids/analytics"
	"rslytics-app-api/internal/controller/kids/health"
	"rslytics-app-api/internal/controller/kids/v1"
)

func Register(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			health.New(),
			// 接口中不属于圈子领域的根路径由专用控制器统一注册。
			v1.New(),
			// 统计接口由独立控制器承载，必须与其余 v1 路由一同注册。
			analytics.NewV1(),
		)
		//group.Group("/kids", func(group *ghttp.RouterGroup) {
		//	group.Bind(
		//		index.New(),
		//		analytics.New(),
		//		user.New(),
		//		circle.New(),
		//		profile.New(),
		//		family.New(),
		//		task.New(),
		//		reward.New(),
		//		star.New(),
		//		notification.New(),
		//		ranking.New(),
		//		upload.New(),
		//		device.New(),
		//	)
		//})
	})
}
