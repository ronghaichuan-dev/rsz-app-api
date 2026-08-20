// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsFamilyMember is the golang structure of table kids_family_member for DAO operations like Where/Data.
type KidsFamilyMember struct {
	g.Meta      `orm:"table:kids_family_member, do:true"`
	Id          any         // 家庭成员ID
	CircleId    any         // 所属圈子ID
	Name        any         // 显示名称
	Gender      any         // 性别：male/female
	Avatar      any         // 头像地址或预设标识
	AvatarStyle any         // 虚拟形象风格标识
	Relation    any         // 家庭关系
	Owner       any         // 是否家庭拥有者
	BindUserId  any         // 绑定用户ID
	BoundAt     *gtime.Time // 绑定时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
