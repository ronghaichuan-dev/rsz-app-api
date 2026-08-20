// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractMember is the golang structure of table kids_contract_member for DAO operations like Where/Data.
type KidsContractMember struct {
	g.Meta         `orm:"table:kids_contract_member, do:true"`
	Id             any         // 主键
	MemberId       any         // 合同成员标识
	CircleId       any         // 合同圈子标识
	BoundAccountId any         // 绑定账号标识
	DisplayName    any         // 显示名称
	Gender         any         // 性别
	Avatar         any         // 头像视觉引用
	Status         any         // 成员状态
	Version        any         // 版本号
	DeletedAt      *gtime.Time // 删除时间
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
