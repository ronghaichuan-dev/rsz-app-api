// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsInviteCode is the golang structure for table kids_invite_code.
type KidsInviteCode struct {
	Id             uint64      `json:"id"             orm:"id"               description:"邀请码ID"`               // 邀请码ID
	Code           string      `json:"code"           orm:"code"             description:"六位邀请码"`               // 六位邀请码
	CircleId       uint64      `json:"circleId"       orm:"circle_id"        description:"圈子ID"`                // 圈子ID
	InviterUserId  uint64      `json:"inviterUserId"  orm:"inviter_user_id"  description:"邀请人用户ID"`             // 邀请人用户ID
	InviteRole     string      `json:"inviteRole"     orm:"invite_role"      description:"邀请加入角色：admin/member"` // 邀请加入角色：admin/member
	TargetMemberId uint64      `json:"targetMemberId" orm:"target_member_id" description:"目标家庭成员ID，0表示不指定"`     // 目标家庭成员ID，0表示不指定
	ExpiredAt      *gtime.Time `json:"expiredAt"      orm:"expired_at"       description:"过期时间"`                // 过期时间
	UsedAt         *gtime.Time `json:"usedAt"         orm:"used_at"          description:"使用时间"`                // 使用时间
	UsedByUserId   uint64      `json:"usedByUserId"   orm:"used_by_user_id"  description:"使用者用户ID"`             // 使用者用户ID
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:"创建时间"`                // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:"更新时间"`                // 更新时间
}
