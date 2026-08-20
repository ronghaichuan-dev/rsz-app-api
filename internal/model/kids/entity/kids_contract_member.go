// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractMember is the golang structure for table kids_contract_member.
type KidsContractMember struct {
	Id             uint64      `json:"id"             orm:"id"               description:"主键"`     // 主键
	MemberId       string      `json:"memberId"       orm:"member_id"        description:"合同成员标识"` // 合同成员标识
	CircleId       string      `json:"circleId"       orm:"circle_id"        description:"合同圈子标识"` // 合同圈子标识
	BoundAccountId string      `json:"boundAccountId" orm:"bound_account_id" description:"绑定账号标识"` // 绑定账号标识
	DisplayName    string      `json:"displayName"    orm:"display_name"     description:"显示名称"`   // 显示名称
	Gender         string      `json:"gender"         orm:"gender"           description:"性别"`     // 性别
	Avatar         string      `json:"avatar"         orm:"avatar"           description:"头像视觉引用"` // 头像视觉引用
	Status         string      `json:"status"         orm:"status"           description:"成员状态"`   // 成员状态
	Version        uint64      `json:"version"        orm:"version"          description:"版本号"`    // 版本号
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"       description:"删除时间"`   // 删除时间
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:"创建时间"`   // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:"更新时间"`   // 更新时间
}
