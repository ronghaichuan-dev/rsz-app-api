// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsAsset is the golang structure of table kids_asset for DAO operations like Where/Data.
type KidsAsset struct {
	g.Meta      `orm:"table:kids_asset, do:true"`
	Id          any         // 主键
	AssetId     any         // 资产标识
	UploadId    any         // 上传标识
	CircleId    any         // 圈子标识
	Purpose     any         // 资产用途
	ContentType any         // 内容类型
	ByteSize    any         // 字节大小
	Sha256      any         // 内容摘要
	State       any         // 资产状态
	Version     any         // 版本号
	CommittedAt *gtime.Time // 提交时间
	CreatedAt   *gtime.Time // 创建时间
}
