// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsInvitation is the golang structure for table kids_invitation.
type KidsInvitation struct {
	Id                    uint64      `json:"id"                    orm:"id"                      description:"主键"`      // 主键
	InviteId              string      `json:"inviteId"              orm:"invite_id"               description:"接口邀请标识"`  // 接口邀请标识
	CircleId              string      `json:"circleId"              orm:"circle_id"               description:"接口圈子标识"`  // 接口圈子标识
	TargetRole            string      `json:"targetRole"            orm:"target_role"             description:"目标角色"`    // 目标角色
	TargetAdministratorId string      `json:"targetAdministratorId" orm:"target_administrator_id" description:"目标管理员标识"` // 目标管理员标识
	TargetMemberId        string      `json:"targetMemberId"        orm:"target_member_id"        description:"目标成员标识"`  // 目标成员标识
	PermissionScope       string      `json:"permissionScope"       orm:"permission_scope"        description:"邀请权限范围"`  // 邀请权限范围
	CodeHash              string      `json:"codeHash"              orm:"code_hash"               description:"邀请码摘要"`   // 邀请码摘要
	SingleUse             int         `json:"singleUse"             orm:"single_use"              description:"是否单次使用"`  // 是否单次使用
	Generation            uint        `json:"generation"            orm:"generation"              description:"邀请码代次"`   // 邀请码代次
	Status                string      `json:"status"                orm:"status"                  description:"邀请状态"`    // 邀请状态
	ExpiresAt             *gtime.Time `json:"expiresAt"             orm:"expires_at"              description:"过期时间"`    // 过期时间
	UsedAt                *gtime.Time `json:"usedAt"                orm:"used_at"                 description:"使用时间"`    // 使用时间
	RevokedAt             *gtime.Time `json:"revokedAt"             orm:"revoked_at"              description:"撤销时间"`    // 撤销时间
	Version               uint64      `json:"version"               orm:"version"                 description:"版本号"`     // 版本号
	CreatedByAccountId    string      `json:"createdByAccountId"    orm:"created_by_account_id"   description:"创建账号标识"`  // 创建账号标识
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"              description:"创建时间"`    // 创建时间
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"              description:"更新时间"`    // 更新时间
}
