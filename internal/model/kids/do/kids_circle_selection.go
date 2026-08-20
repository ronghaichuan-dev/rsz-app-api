// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircleSelection is the golang structure of table kids_circle_selection for DAO operations like Where/Data.
type KidsCircleSelection struct {
	g.Meta          `orm:"table:kids_circle_selection, do:true"`
	Id              any         // 主键
	SelectionId     any         // 接口选择标识
	AccountId       any         // 接口账号标识
	CurrentCircleId any         // 当前圈子标识
	Version         any         // 版本号
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
