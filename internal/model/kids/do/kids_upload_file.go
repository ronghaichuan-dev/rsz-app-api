// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUploadFile is the golang structure of table kids_upload_file for DAO operations like Where/Data.
type KidsUploadFile struct {
	g.Meta      `orm:"table:kids_upload_file, do:true"`
	Id          any         // 上传文件ID
	UserId      any         // 上传用户ID
	MemberId    any         // 绑定成员ID
	UsageType   any         // 文件用途：avatar/task_photo/reward
	FileName    any         // 文件名称
	FileUrl     any         // 文件访问地址
	ContentType any         // 文件类型
	FileSize    any         // 文件大小字节数
	CreatedAt   *gtime.Time // 创建时间
}
