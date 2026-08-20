// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsStarRecord is the golang structure of table kids_star_record for DAO operations like Where/Data.
type KidsStarRecord struct {
	g.Meta       `orm:"table:kids_star_record, do:true"`
	Id           any         // 星星流水ID
	KidId        any         // 儿童成员ID
	ChangeAmount any         // 星星变动数量
	Balance      any         // 变动后余额
	RecordType   any         // 流水类型：task/reward/adjustment
	Title        any         // 流水标题
	Remark       any         // 流水备注
	CreatedAt    *gtime.Time // 创建时间
}
