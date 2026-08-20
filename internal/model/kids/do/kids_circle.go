// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircle is the golang structure of table kids_circle for DAO operations like Where/Data.
type KidsCircle struct {
	g.Meta      `orm:"table:kids_circle, do:true"`
	Id          any         // 圈子ID
	Name        any         // 圈子名称
	Icon        any         // 圈子图标标识
	OwnerUserId any         // 创建者用户ID
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
