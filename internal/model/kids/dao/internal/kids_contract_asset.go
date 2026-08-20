// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractAssetDao is the data access object for the table kids_contract_asset.
type KidsContractAssetDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  KidsContractAssetColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// KidsContractAssetColumns defines and stores column names for the table kids_contract_asset.
type KidsContractAssetColumns struct {
	Id          string // 主键
	AssetId     string // 资产标识
	UploadId    string // 上传标识
	CircleId    string // 圈子标识
	Purpose     string // 资产用途
	ContentType string // 内容类型
	ByteSize    string // 字节大小
	Sha256      string // 内容摘要
	Version     string // 版本号
	CreatedAt   string // 创建时间
}

// kidsContractAssetColumns holds the columns for the table kids_contract_asset.
var kidsContractAssetColumns = KidsContractAssetColumns{
	Id:          "id",
	AssetId:     "asset_id",
	UploadId:    "upload_id",
	CircleId:    "circle_id",
	Purpose:     "purpose",
	ContentType: "content_type",
	ByteSize:    "byte_size",
	Sha256:      "sha256",
	Version:     "version",
	CreatedAt:   "created_at",
}

// NewKidsContractAssetDao creates and returns a new DAO object for table data access.
func NewKidsContractAssetDao(handlers ...gdb.ModelHandler) *KidsContractAssetDao {
	return &KidsContractAssetDao{
		group:    "kids",
		table:    "kids_contract_asset",
		columns:  kidsContractAssetColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractAssetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractAssetDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractAssetDao) Columns() KidsContractAssetColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractAssetDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractAssetDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractAssetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
