// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractAccountBinding is the golang structure for table kids_contract_account_binding.
type KidsContractAccountBinding struct {
	Id              uint64      `json:"id"              orm:"id"               description:"主键"`     // 主键
	BindingId       string      `json:"bindingId"       orm:"binding_id"       description:"合同绑定标识"` // 合同绑定标识
	AccountId       string      `json:"accountId"       orm:"account_id"       description:"合同账号标识"` // 合同账号标识
	Environment     string      `json:"environment"     orm:"environment"      description:"绑定环境"`   // 绑定环境
	MigrationPolicy string      `json:"migrationPolicy" orm:"migration_policy" description:"迁移策略"`   // 迁移策略
	Version         uint64      `json:"version"         orm:"version"          description:"版本号"`    // 版本号
	IssuedAt        *gtime.Time `json:"issuedAt"        orm:"issued_at"        description:"签发时间"`   // 签发时间
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`   // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:"更新时间"`   // 更新时间
}
