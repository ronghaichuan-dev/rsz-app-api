// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsStarRecord is the golang structure for table kids_star_record.
type KidsStarRecord struct {
	Id           uint64      `json:"id"           orm:"id"            description:"星星流水ID"`                      // 星星流水ID
	KidId        uint64      `json:"kidId"        orm:"kid_id"        description:"儿童成员ID"`                      // 儿童成员ID
	ChangeAmount int         `json:"changeAmount" orm:"change_amount" description:"星星变动数量"`                      // 星星变动数量
	Balance      int         `json:"balance"      orm:"balance"       description:"变动后余额"`                       // 变动后余额
	RecordType   string      `json:"recordType"   orm:"record_type"   description:"流水类型：task/reward/adjustment"` // 流水类型：task/reward/adjustment
	Title        string      `json:"title"        orm:"title"         description:"流水标题"`                        // 流水标题
	Remark       string      `json:"remark"       orm:"remark"        description:"流水备注"`                        // 流水备注
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`                        // 创建时间
}
