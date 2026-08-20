package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	UploadUsageAvatar    = "avatar"
	UploadUsageTaskPhoto = "task_photo"
	UploadUsageReward    = "reward"
)

type UploadFileReq struct {
	g.Meta      `path:"/uploads" method:"post" tags:"儿童端上传" summary:"登记上传文件"`
	MemberId    uint64 `json:"memberId" dc:"绑定成员ID"`
	UsageType   string `json:"usageType" v:"required|in:avatar,task_photo,reward" dc:"文件用途：avatar/task_photo/reward"`
	FileName    string `json:"fileName" dc:"文件名称"`
	FileUrl     string `json:"fileUrl" v:"required|max-length:512" dc:"文件访问地址"`
	ContentType string `json:"contentType" dc:"文件类型"`
	FileSize    uint64 `json:"fileSize" dc:"文件大小字节数"`
}

type UploadFileInput struct {
	UserId      uint64
	MemberId    uint64
	UsageType   string
	FileName    string
	FileUrl     string
	ContentType string
	FileSize    uint64
}

type UploadFileOutput struct {
	Id      uint64
	FileUrl string
}

type UploadFileRes struct {
	Id      uint64 `json:"id" dc:"上传文件ID"`
	FileUrl string `json:"fileUrl" dc:"文件访问地址"`
}
