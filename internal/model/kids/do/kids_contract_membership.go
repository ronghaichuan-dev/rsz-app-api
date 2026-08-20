// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractMembership is the golang structure of table kids_contract_membership for DAO operations like Where/Data.
type KidsContractMembership struct {
	g.Meta       `orm:"table:kids_contract_membership, do:true"`
	Id           any         // 主键
	MembershipId any         // 合同成员身份标识
	CircleId     any         // 合同圈子标识
	AccountId    any         // 合同账号标识
	ActorType    any         // 授权主体类型
	ActorId      any         // 授权主体标识
	Role         any         // 成员角色
	Permissions  any         // 权限集合
	Status       any         // 成员身份状态
	Version      any         // 版本号
	DeletedAt    *gtime.Time // 删除时间
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
