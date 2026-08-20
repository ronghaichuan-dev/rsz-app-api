// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUser is the golang structure for table kids_user.
type KidsUser struct {
	Id        uint64      `json:"id"        orm:"id"         description:"用户ID"`                      // 用户ID
	DeviceId  string      `json:"deviceId"  orm:"device_id"  description:"最近登录设备ID"`                  // 最近登录设备ID
	Provider  string      `json:"provider"  orm:"provider"   description:"当前登录方式：guest/google/apple"` // 当前登录方式：guest/google/apple
	Email     string      `json:"email"     orm:"email"      description:"授权服务商返回的邮箱"`                // 授权服务商返回的邮箱
	Nickname  string      `json:"nickname"  orm:"nickname"   description:"昵称"`                        // 昵称
	Avatar    string      `json:"avatar"    orm:"avatar"     description:"头像地址"`                      // 头像地址
	IsGuest   uint        `json:"isGuest"   orm:"is_guest"   description:"是否游客账号"`                    // 是否游客账号
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                      // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                      // 更新时间
}
