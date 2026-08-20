// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsSyncCommit is the golang structure of table kids_sync_commit for DAO operations like Where/Data.
type KidsSyncCommit struct {
	g.Meta         `orm:"table:kids_sync_commit, do:true"`
	Id             any         // 主键
	CommitId       any         // 接口提交标识
	CircleId       any         // 圈子标识
	CommitSequence any         // 单调提交序列
	ChangePayload  any         // 完整变更集合
	CreatedAt      *gtime.Time // 创建时间
}
