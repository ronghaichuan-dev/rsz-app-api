// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractAccountBinding is the golang structure of table kids_contract_account_binding for DAO operations like Where/Data.
type KidsContractAccountBinding struct {
	g.Meta          `orm:"table:kids_contract_account_binding, do:true"`
	Id              any         // 主键
	BindingId       any         // 合同绑定标识
	AccountId       any         // 合同账号标识
	Environment     any         // 绑定环境
	MigrationPolicy any         // 迁移策略
	Version         any         // 版本号
	IssuedAt        *gtime.Time // 签发时间
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
