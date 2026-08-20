package kids

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// SaveDeviceNotification 保存设备通知授权和偏好。
func (s *sKids) SaveDeviceNotification(ctx context.Context, in v1.DeviceNotificationSaveInput) (*v1.DeviceNotificationSaveOutput, error) {
	deviceId := strings.TrimSpace(in.DeviceId)
	if deviceId == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "deviceId is required")
	}
	if !validDevicePlatform(in.Platform) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "platform must be ios, android or web")
	}
	_, err := utils.KidsDB(ctx).Model(consts.KidsDeviceNotificationTable).Ctx(ctx).Data(map[string]any{
		"user_id":        in.UserId,
		"device_id":      deviceId,
		"platform":       strings.TrimSpace(in.Platform),
		"device_token":   strings.TrimSpace(in.DeviceToken),
		"authorized":     utils.BoolToInt(in.Authorized),
		"task_enabled":   utils.BoolToInt(in.TaskEnabled),
		"reward_enabled": utils.BoolToInt(in.RewardEnabled),
		"member_enabled": utils.BoolToInt(in.MemberEnabled),
	}).Save()
	if err != nil {
		return nil, err
	}
	return &v1.DeviceNotificationSaveOutput{DeviceId: deviceId}, nil
}

// validDevicePlatform 校验设备平台枚举，空值允许兼容旧客户端。
func validDevicePlatform(platform string) bool {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "", v1.DevicePlatformIOS, v1.DevicePlatformAndroid, v1.DevicePlatformWeb:
		return true
	default:
		return false
	}
}
