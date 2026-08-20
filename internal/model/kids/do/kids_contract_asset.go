// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractAsset is the golang structure of table kids_contract_asset for DAO operations like Where/Data.
type KidsContractAsset struct {
	g.Meta      `orm:"table:kids_contract_asset, do:true"`
	Id          any         // 主键
	AssetId     any         // 资产标识
	UploadId    any         // 上传标识
	CircleId    any         // 圈子标识
	Purpose     any         // 资产用途
	ContentType any         // 内容类型
	ByteSize    any         // 字节大小
	Sha256      any         // 内容摘要
	Version     any         // 版本号
	CreatedAt   *gtime.Time // 创建时间
}
