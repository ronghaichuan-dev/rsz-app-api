// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUser is the golang structure of table kids_user for DAO operations like Where/Data.
type KidsUser struct {
	g.Meta    `orm:"table:kids_user, do:true"`
	Id        any         // 用户ID
	DeviceId  any         // 最近登录设备ID
	Provider  any         // 当前登录方式：guest/google/apple
	Email     any         // 授权服务商返回的邮箱
	Nickname  any         // 昵称
	Avatar    any         // 头像地址
	IsGuest   any         // 是否游客账号
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
