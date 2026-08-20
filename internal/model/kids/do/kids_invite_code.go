// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsInviteCode is the golang structure of table kids_invite_code for DAO operations like Where/Data.
type KidsInviteCode struct {
	g.Meta         `orm:"table:kids_invite_code, do:true"`
	Id             any         // 邀请码ID
	Code           any         // 六位邀请码
	CircleId       any         // 圈子ID
	InviterUserId  any         // 邀请人用户ID
	InviteRole     any         // 邀请加入角色：admin/member
	TargetMemberId any         // 目标家庭成员ID，0表示不指定
	ExpiredAt      *gtime.Time // 过期时间
	UsedAt         *gtime.Time // 使用时间
	UsedByUserId   any         // 使用者用户ID
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
