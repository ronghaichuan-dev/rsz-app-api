// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsInvitation is the golang structure of table kids_invitation for DAO operations like Where/Data.
type KidsInvitation struct {
	g.Meta                `orm:"table:kids_invitation, do:true"`
	Id                    any         // 主键
	InviteId              any         // 接口邀请标识
	CircleId              any         // 接口圈子标识
	TargetRole            any         // 目标角色
	TargetAdministratorId any         // 目标管理员标识
	TargetMemberId        any         // 目标成员标识
	PermissionScope       any         // 邀请权限范围
	CodeHash              any         // 邀请码摘要
	SingleUse             any         // 是否单次使用
	Generation            any         // 邀请码代次
	Status                any         // 邀请状态
	ExpiresAt             *gtime.Time // 过期时间
	UsedAt                *gtime.Time // 使用时间
	RevokedAt             *gtime.Time // 撤销时间
	Version               any         // 版本号
	CreatedByAccountId    any         // 创建账号标识
	CreatedAt             *gtime.Time // 创建时间
	UpdatedAt             *gtime.Time // 更新时间
}
