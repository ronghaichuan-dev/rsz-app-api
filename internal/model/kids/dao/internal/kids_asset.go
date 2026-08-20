// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsAssetDao is the data access object for the table kids_asset.
type KidsAssetDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KidsAssetColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KidsAssetColumns defines and stores column names for the table kids_asset.
type KidsAssetColumns struct {
	Id          string // 主键
	AssetId     string // 资产标识
	UploadId    string // 上传标识
	CircleId    string // 圈子标识
	Purpose     string // 资产用途
	ContentType string // 内容类型
	ByteSize    string // 字节大小
	Sha256      string // 内容摘要
	State       string // 资产状态
	Version     string // 版本号
	CommittedAt string // 提交时间
	CreatedAt   string // 创建时间
}

// kidsAssetColumns holds the columns for the table kids_asset.
var kidsAssetColumns = KidsAssetColumns{
	Id:          "id",
	AssetId:     "asset_id",
	UploadId:    "upload_id",
	CircleId:    "circle_id",
	Purpose:     "purpose",
	ContentType: "content_type",
	ByteSize:    "byte_size",
	Sha256:      "sha256",
	State:       "state",
	Version:     "version",
	CommittedAt: "committed_at",
	CreatedAt:   "created_at",
}

// NewKidsAssetDao creates and returns a new DAO object for table data access.
func NewKidsAssetDao(handlers ...gdb.ModelHandler) *KidsAssetDao {
	return &KidsAssetDao{
		group:    "kids",
		table:    "kids_asset",
		columns:  kidsAssetColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsAssetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsAssetDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsAssetDao) Columns() KidsAssetColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsAssetDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsAssetDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *KidsAssetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
