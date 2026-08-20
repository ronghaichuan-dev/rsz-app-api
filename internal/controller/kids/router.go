package kids

import (
	"context"

	"rslytics-app-api/internal/controller/kids/analytics"
	"rslytics-app-api/internal/controller/kids/index"

	"github.com/gogf/gf/v2/net/ghttp"

	"rslytics-app-api/internal/controller/kids/circle"
	"rslytics-app-api/internal/controller/kids/device"
	"rslytics-app-api/internal/controller/kids/family"
	"rslytics-app-api/internal/controller/kids/health"
	"rslytics-app-api/internal/controller/kids/notification"
	"rslytics-app-api/internal/controller/kids/profile"
	"rslytics-app-api/internal/controller/kids/ranking"
	"rslytics-app-api/internal/controller/kids/reward"
	"rslytics-app-api/internal/controller/kids/star"
	"rslytics-app-api/internal/controller/kids/task"
	"rslytics-app-api/internal/controller/kids/upload"
	"rslytics-app-api/internal/controller/kids/user"
)

func Register(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			health.New(),
			analytics.NewV1(),
		)
		group.Group("/kids", func(group *ghttp.RouterGroup) {
			group.Bind(
				index.New(),
				analytics.New(),
				user.New(),
				circle.New(),
				profile.New(),
				family.New(),
				task.New(),
				reward.New(),
				star.New(),
				notification.New(),
				ranking.New(),
				upload.New(),
				device.New(),
			)
		})
	})
}
