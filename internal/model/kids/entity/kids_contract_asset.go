// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractAsset is the golang structure for table kids_contract_asset.
type KidsContractAsset struct {
	Id          uint64      `json:"id"          orm:"id"           description:"主键"`   // 主键
	AssetId     string      `json:"assetId"     orm:"asset_id"     description:"资产标识"` // 资产标识
	UploadId    string      `json:"uploadId"    orm:"upload_id"    description:"上传标识"` // 上传标识
	CircleId    string      `json:"circleId"    orm:"circle_id"    description:"圈子标识"` // 圈子标识
	Purpose     string      `json:"purpose"     orm:"purpose"      description:"资产用途"` // 资产用途
	ContentType string      `json:"contentType" orm:"content_type" description:"内容类型"` // 内容类型
	ByteSize    uint64      `json:"byteSize"    orm:"byte_size"    description:"字节大小"` // 字节大小
	Sha256      string      `json:"sha256"      orm:"sha256"       description:"内容摘要"` // 内容摘要
	Version     uint64      `json:"version"     orm:"version"      description:"版本号"`  // 版本号
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"` // 创建时间
}
