// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircleUser is the golang structure of table kids_circle_user for DAO operations like Where/Data.
type KidsCircleUser struct {
	g.Meta    `orm:"table:kids_circle_user, do:true"`
	Id        any         // 圈子用户关系ID
	CircleId  any         // 圈子ID
	UserId    any         // 用户ID
	Role      any         // 圈子角色：admin/member
	MemberId  any         // 绑定的家庭成员ID
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
	LeftAt    *gtime.Time // 退出时间
}
