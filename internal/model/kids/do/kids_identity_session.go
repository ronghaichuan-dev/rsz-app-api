// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsIdentitySession is the golang structure of table kids_identity_session for DAO operations like Where/Data.
type KidsIdentitySession struct {
	g.Meta                `orm:"table:kids_identity_session, do:true"`
	Id                    any         // 主键
	SessionId             any         // 接口会话标识
	AccountId             any         // 接口账号标识
	PrincipalKind         any         // 主体类型
	AccessTokenHash       any         // 访问令牌摘要
	RefreshTokenHash      any         // 刷新令牌摘要
	GuestUpgradeGrantHash any         // 游客升级凭据摘要
	ExpiresAt             *gtime.Time // 访问令牌到期时间
	RefreshExpiresAt      *gtime.Time // 刷新令牌到期时间
	RevokedAt             *gtime.Time // 撤销时间
	CreatedAt             *gtime.Time // 创建时间
	UpdatedAt             *gtime.Time // 更新时间
}
