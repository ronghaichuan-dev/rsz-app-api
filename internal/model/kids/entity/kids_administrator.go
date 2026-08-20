// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsAdministrator is the golang structure for table kids_administrator.
type KidsAdministrator struct {
	Id              uint64      `json:"id"              orm:"id"               description:"主键"`      // 主键
	AdministratorId string      `json:"administratorId" orm:"administrator_id" description:"接口管理员标识"` // 接口管理员标识
	CircleId        string      `json:"circleId"        orm:"circle_id"        description:"接口圈子标识"`  // 接口圈子标识
	BoundAccountId  string      `json:"boundAccountId"  orm:"bound_account_id" description:"绑定账号标识"`  // 绑定账号标识
	DisplayName     string      `json:"displayName"     orm:"display_name"     description:"显示名称"`    // 显示名称
	Avatar          string      `json:"avatar"          orm:"avatar"           description:"头像视觉引用"`  // 头像视觉引用
	Role            string      `json:"role"            orm:"role"             description:"管理员角色"`   // 管理员角色
	Permissions     string      `json:"permissions"     orm:"permissions"      description:"权限集合"`    // 权限集合
	Status          string      `json:"status"          orm:"status"           description:"管理员状态"`   // 管理员状态
	Version         uint64      `json:"version"         orm:"version"          description:"版本号"`     // 版本号
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"       description:"删除时间"`    // 删除时间
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`    // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:"更新时间"`    // 更新时间
}
