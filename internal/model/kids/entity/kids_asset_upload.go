// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsAssetUpload is the golang structure for table kids_asset_upload.
type KidsAssetUpload struct {
	Id          uint64      `json:"id"          orm:"id"           description:"主键"`   // 主键
	UploadId    string      `json:"uploadId"    orm:"upload_id"    description:"上传标识"` // 上传标识
	AccountId   string      `json:"accountId"   orm:"account_id"   description:"账号标识"` // 账号标识
	CircleId    string      `json:"circleId"    orm:"circle_id"    description:"圈子标识"` // 圈子标识
	Purpose     string      `json:"purpose"     orm:"purpose"      description:"上传用途"` // 上传用途
	ContentType string      `json:"contentType" orm:"content_type" description:"内容类型"` // 内容类型
	ByteSize    uint64      `json:"byteSize"    orm:"byte_size"    description:"字节大小"` // 字节大小
	Sha256      string      `json:"sha256"      orm:"sha256"       description:"内容摘要"` // 内容摘要
	Version     uint64      `json:"version"     orm:"version"      description:"版本号"`  // 版本号
	Status      string      `json:"status"      orm:"status"       description:"上传状态"` // 上传状态
	ExpiresAt   *gtime.Time `json:"expiresAt"   orm:"expires_at"   description:"到期时间"` // 到期时间
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"` // 更新时间
}
