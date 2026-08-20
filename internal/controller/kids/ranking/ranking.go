package ranking

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

// Controller 处理排行榜接口。
type Controller struct{}

// New 创建排行榜控制器实例。
func New() *Controller {
	return &Controller{}
}

// Stars 查询群组星星排行榜。
func (c *Controller) Stars(ctx context.Context, req *v1.StarRankingReq) (res *v1.StarRankingRes, err error) {
	out, err := service.Kids().GetStarRanking(ctx, v1.StarRankingInput{UserId: authUserId(ctx), CircleId: req.CircleId})
	if err != nil {
		return nil, err
	}
	return &v1.StarRankingRes{List: out.List}, nil
}

// authUserId 从上下文读取当前登录用户ID。
func authUserId(ctx context.Context) uint64 {
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}
