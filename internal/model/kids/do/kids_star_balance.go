// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsStarBalance is the golang structure of table kids_star_balance for DAO operations like Where/Data.
type KidsStarBalance struct {
	g.Meta               `orm:"table:kids_star_balance, do:true"`
	Id                   any         // 主键
	CircleId             any         // 接口圈子标识
	MemberId             any         // 接口成员标识
	Balance              any         // 星星余额
	Version              any         // 版本号
	SourceCommitId       any         // 来源提交标识
	SourceCommitSequence any         // 来源提交序列
	UpdatedAt            *gtime.Time // 更新时间
}
