// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsUploadFileDao is the data access object for the table kids_upload_file.
type KidsUploadFileDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  KidsUploadFileColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// KidsUploadFileColumns defines and stores column names for the table kids_upload_file.
type KidsUploadFileColumns struct {
	Id          string // 上传文件ID
	UserId      string // 上传用户ID
	MemberId    string // 绑定成员ID
	UsageType   string // 文件用途：avatar/task_photo/reward
	FileName    string // 文件名称
	FileUrl     string // 文件访问地址
	ContentType string // 文件类型
	FileSize    string // 文件大小字节数
	CreatedAt   string // 创建时间
}

// kidsUploadFileColumns holds the columns for the table kids_upload_file.
var kidsUploadFileColumns = KidsUploadFileColumns{
	Id:          "id",
	UserId:      "user_id",
	MemberId:    "member_id",
	UsageType:   "usage_type",
	FileName:    "file_name",
	FileUrl:     "file_url",
	ContentType: "content_type",
	FileSize:    "file_size",
	CreatedAt:   "created_at",
}

// NewKidsUploadFileDao creates and returns a new DAO object for table data access.
func NewKidsUploadFileDao(handlers ...gdb.ModelHandler) *KidsUploadFileDao {
	return &KidsUploadFileDao{
		group:    "kids",
		table:    "kids_upload_file",
		columns:  kidsUploadFileColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsUploadFileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsUploadFileDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsUploadFileDao) Columns() KidsUploadFileColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsUploadFileDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsUploadFileDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsUploadFileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
