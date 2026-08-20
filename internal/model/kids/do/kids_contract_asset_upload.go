// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractAssetUpload is the golang structure of table kids_contract_asset_upload for DAO operations like Where/Data.
type KidsContractAssetUpload struct {
	g.Meta      `orm:"table:kids_contract_asset_upload, do:true"`
	Id          any         // 主键
	UploadId    any         // 上传标识
	AccountId   any         // 账号标识
	CircleId    any         // 圈子标识
	Purpose     any         // 上传用途
	ContentType any         // 内容类型
	ByteSize    any         // 字节大小
	Sha256      any         // 内容摘要
	Version     any         // 版本号
	Status      any         // 上传状态
	ExpiresAt   *gtime.Time // 到期时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
