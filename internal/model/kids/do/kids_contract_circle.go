// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractCircle is the golang structure of table kids_contract_circle for DAO operations like Where/Data.
type KidsContractCircle struct {
	g.Meta               `orm:"table:kids_contract_circle, do:true"`
	Id                   any         // 主键
	CircleId             any         // 合同圈子标识
	Name                 any         // 圈子名称
	Icon                 any         // 圈子视觉引用
	OwnerAdministratorId any         // 所有者管理员标识
	Status               any         // 圈子状态
	Version              any         // 版本号
	DeletedAt            *gtime.Time // 删除时间
	CreatedAt            *gtime.Time // 创建时间
	UpdatedAt            *gtime.Time // 更新时间
}
