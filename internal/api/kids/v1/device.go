package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	DevicePlatformIOS     = "ios"
	DevicePlatformAndroid = "android"
	DevicePlatformWeb     = "web"
)

type DeviceNotificationSaveReq struct {
	g.Meta        `path:"/devices/notification" method:"put" tags:"儿童端设备" summary:"保存设备通知设置"`
	DeviceId      string `json:"deviceId" v:"required|max-length:128" dc:"设备ID"`
	Platform      string `json:"platform" v:"in:ios,android,web" dc:"平台：ios/android/web"`
	DeviceToken   string `json:"deviceToken" dc:"推送设备令牌"`
	Authorized    bool   `json:"authorized" dc:"是否已授权通知"`
	TaskEnabled   bool   `json:"taskEnabled" dc:"是否开启任务提醒"`
	RewardEnabled bool   `json:"rewardEnabled" dc:"是否开启奖励提醒"`
	MemberEnabled bool   `json:"memberEnabled" dc:"是否开启成员提醒"`
}

type DeviceNotificationSaveInput struct {
	UserId        uint64
	DeviceId      string
	Platform      string
	DeviceToken   string
	Authorized    bool
	TaskEnabled   bool
	RewardEnabled bool
	MemberEnabled bool
}

type DeviceNotificationSaveOutput struct {
	DeviceId string
}

type DeviceNotificationSaveRes struct {
	DeviceId string `json:"deviceId" dc:"设备ID"`
}
