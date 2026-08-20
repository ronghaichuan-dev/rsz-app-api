// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsDeviceNotification is the golang structure of table kids_device_notification for DAO operations like Where/Data.
type KidsDeviceNotification struct {
	g.Meta        `orm:"table:kids_device_notification, do:true"`
	Id            any         // 设备通知ID
	UserId        any         // 用户ID
	DeviceId      any         // 设备ID
	Platform      any         // 平台：ios/android/web
	DeviceToken   any         // 推送设备令牌
	Authorized    any         // 是否已授权通知
	TaskEnabled   any         // 是否开启任务提醒
	RewardEnabled any         // 是否开启奖励提醒
	MemberEnabled any         // 是否开启成员提醒
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
