// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsSyncSequence is the golang structure of table kids_sync_sequence for DAO operations like Where/Data.
type KidsSyncSequence struct {
	g.Meta             `orm:"table:kids_sync_sequence, do:true"`
	Id                 any         // 固定序列标识
	NextCommitSequence any         // 下一个提交序列
	UpdatedAt          *gtime.Time // 更新时间
}
