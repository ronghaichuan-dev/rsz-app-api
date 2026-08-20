// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsAccount is the golang structure of table kids_account for DAO operations like Where/Data.
type KidsAccount struct {
	g.Meta    `orm:"table:kids_account, do:true"`
	Id        any         // 主键
	AccountId any         // 接口账号标识
	Status    any         // 账号状态
	Version   any         // 版本号
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
