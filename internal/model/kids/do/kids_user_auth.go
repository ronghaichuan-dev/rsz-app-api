// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUserAuth is the golang structure of table kids_user_auth for DAO operations like Where/Data.
type KidsUserAuth struct {
	g.Meta    `orm:"table:kids_user_auth, do:true"`
	Id        any         // 授权记录ID
	UserId    any         // kids用户ID
	Provider  any         // 授权服务商：google/apple
	OpenId    any         // 服务商开放ID或主体标识
	Email     any         // 服务商邮箱
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
