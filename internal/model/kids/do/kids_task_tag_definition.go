// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskTagDefinition is the golang structure of table kids_task_tag_definition for DAO operations like Where/Data.
type KidsTaskTagDefinition struct {
	g.Meta    `orm:"table:kids_task_tag_definition, do:true"`
	Id        any         // 主键
	TaskTagId any         // 接口任务标签标识
	CircleId  any         // 接口圈子标识
	Name      any         // 标签名称
	Status    any         // 标签状态
	Version   any         // 版本号
	DeletedAt *gtime.Time // 删除时间
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
