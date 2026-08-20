// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsAssetUploadDao is the data access object for the table kids_asset_upload.
type KidsAssetUploadDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  KidsAssetUploadColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// KidsAssetUploadColumns defines and stores column names for the table kids_asset_upload.
type KidsAssetUploadColumns struct {
	Id          string // 主键
	UploadId    string // 上传标识
	AccountId   string // 账号标识
	CircleId    string // 圈子标识
	Purpose     string // 上传用途
	ContentType string // 内容类型
	ByteSize    string // 字节大小
	Sha256      string // 内容摘要
	Version     string // 版本号
	Status      string // 上传状态
	ExpiresAt   string // 到期时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// kidsAssetUploadColumns holds the columns for the table kids_asset_upload.
var kidsAssetUploadColumns = KidsAssetUploadColumns{
	Id:          "id",
	UploadId:    "upload_id",
	AccountId:   "account_id",
	CircleId:    "circle_id",
	Purpose:     "purpose",
	ContentType: "content_type",
	ByteSize:    "byte_size",
	Sha256:      "sha256",
	Version:     "version",
	Status:      "status",
	ExpiresAt:   "expires_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewKidsAssetUploadDao creates and returns a new DAO object for table data access.
func NewKidsAssetUploadDao(handlers ...gdb.ModelHandler) *KidsAssetUploadDao {
	return &KidsAssetUploadDao{
		group:    "kids",
		table:    "kids_asset_upload",
		columns:  kidsAssetUploadColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsAssetUploadDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsAssetUploadDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsAssetUploadDao) Columns() KidsAssetUploadColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsAssetUploadDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsAssetUploadDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsAssetUploadDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
