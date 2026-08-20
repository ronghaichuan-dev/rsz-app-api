// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUserToken is the golang structure of table kids_user_token for DAO operations like Where/Data.
type KidsUserToken struct {
	g.Meta    `orm:"table:kids_user_token, do:true"`
	Id        any         // 令牌ID
	UserId    any         // kids用户ID
	Token     any         // 访问令牌
	ExpiredAt *gtime.Time // 过期时间
	CreatedAt *gtime.Time // 创建时间
}
