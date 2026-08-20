// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractAccount is the golang structure for table kids_contract_account.
type KidsContractAccount struct {
	Id        uint64      `json:"id"        orm:"id"         description:"主键"`     // 主键
	AccountId string      `json:"accountId" orm:"account_id" description:"合同账号标识"` // 合同账号标识
	Status    string      `json:"status"    orm:"status"     description:"账号状态"`   // 账号状态
	Version   uint64      `json:"version"   orm:"version"    description:"版本号"`    // 版本号
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`   // 更新时间
}
