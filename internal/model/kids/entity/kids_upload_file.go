// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUploadFile is the golang structure for table kids_upload_file.
type KidsUploadFile struct {
	Id          uint64      `json:"id"          orm:"id"           description:"上传文件ID"`                        // 上传文件ID
	UserId      uint64      `json:"userId"      orm:"user_id"      description:"上传用户ID"`                        // 上传用户ID
	MemberId    uint64      `json:"memberId"    orm:"member_id"    description:"绑定成员ID"`                        // 绑定成员ID
	UsageType   string      `json:"usageType"   orm:"usage_type"   description:"文件用途：avatar/task_photo/reward"` // 文件用途：avatar/task_photo/reward
	FileName    string      `json:"fileName"    orm:"file_name"    description:"文件名称"`                          // 文件名称
	FileUrl     string      `json:"fileUrl"     orm:"file_url"     description:"文件访问地址"`                        // 文件访问地址
	ContentType string      `json:"contentType" orm:"content_type" description:"文件类型"`                          // 文件类型
	FileSize    uint64      `json:"fileSize"    orm:"file_size"    description:"文件大小字节数"`                       // 文件大小字节数
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`                          // 创建时间
}
