// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskPreset is the golang structure of table kids_task_preset for DAO operations like Where/Data.
type KidsTaskPreset struct {
	g.Meta      `orm:"table:kids_task_preset, do:true"`
	Id          any         // 任务预设ID
	Title       any         // 预设任务标题
	Icon        any         // 图标标识
	Star        any         // 默认星星数量
	NeedPhoto   any         // 是否建议照片凭证
	Description any         // 预设描述
	Enabled     any         // 是否启用预设
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
