package device

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

// Controller 处理设备通知设置接口。
type Controller struct{}

// New 创建设备控制器实例。
func New() *Controller {
	return &Controller{}
}

// SaveNotification 保存设备通知授权和提醒偏好。
func (c *Controller) SaveNotification(ctx context.Context, req *v1.DeviceNotificationSaveReq) (res *v1.DeviceNotificationSaveRes, err error) {
	out, err := service.Kids().SaveDeviceNotification(ctx, v1.DeviceNotificationSaveInput{UserId: authUserId(ctx), DeviceId: req.DeviceId, Platform: req.Platform, DeviceToken: req.DeviceToken, Authorized: req.Authorized, TaskEnabled: req.TaskEnabled, RewardEnabled: req.RewardEnabled, MemberEnabled: req.MemberEnabled})
	if err != nil {
		return nil, err
	}
	return &v1.DeviceNotificationSaveRes{DeviceId: out.DeviceId}, nil
}

// authUserId 从上下文读取当前登录用户ID。
func authUserId(ctx context.Context) uint64 {
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}
