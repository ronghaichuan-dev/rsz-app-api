// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskTag is the golang structure of table kids_task_tag for DAO operations like Where/Data.
type KidsTaskTag struct {
	g.Meta    `orm:"table:kids_task_tag, do:true"`
	Id        any         // 任务标签ID
	Name      any         // 标签名称
	Color     any         // 标签颜色
	SortOrder any         // 排序值
	DeletedAt *gtime.Time // 删除时间
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
