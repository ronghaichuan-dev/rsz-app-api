// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsAdministrator is the golang structure of table kids_administrator for DAO operations like Where/Data.
type KidsAdministrator struct {
	g.Meta          `orm:"table:kids_administrator, do:true"`
	Id              any         // 主键
	AdministratorId any         // 接口管理员标识
	CircleId        any         // 接口圈子标识
	BoundAccountId  any         // 绑定账号标识
	DisplayName     any         // 显示名称
	Avatar          any         // 头像视觉引用
	Role            any         // 管理员角色
	Permissions     any         // 权限集合
	Status          any         // 管理员状态
	Version         any         // 版本号
	DeletedAt       *gtime.Time // 删除时间
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
