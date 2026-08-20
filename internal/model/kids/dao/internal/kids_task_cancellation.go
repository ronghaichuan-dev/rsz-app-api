// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsTaskCancellationDao is the data access object for the table kids_task_cancellation.
type KidsTaskCancellationDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  KidsTaskCancellationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// KidsTaskCancellationColumns defines and stores column names for the table kids_task_cancellation.
type KidsTaskCancellationColumns struct {
	Id             string // 主键
	CancellationId string // 接口取消标识
	CompletionId   string // 接口完成标识
	ReasonCode     string // 取消原因
	CancelledBy    string // 取消操作者快照
	CancelledAt    string // 取消时间
	CommitSequence string // 提交序列
	CreatedAt      string // 创建时间
}

// kidsTaskCancellationColumns holds the columns for the table kids_task_cancellation.
var kidsTaskCancellationColumns = KidsTaskCancellationColumns{
	Id:             "id",
	CancellationId: "cancellation_id",
	CompletionId:   "completion_id",
	ReasonCode:     "reason_code",
	CancelledBy:    "cancelled_by",
	CancelledAt:    "cancelled_at",
	CommitSequence: "commit_sequence",
	CreatedAt:      "created_at",
}

// NewKidsTaskCancellationDao creates and returns a new DAO object for table data access.
func NewKidsTaskCancellationDao(handlers ...gdb.ModelHandler) *KidsTaskCancellationDao {
	return &KidsTaskCancellationDao{
		group:    "kids",
		table:    "kids_task_cancellation",
		columns:  kidsTaskCancellationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsTaskCancellationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsTaskCancellationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsTaskCancellationDao) Columns() KidsTaskCancellationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsTaskCancellationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsTaskCancellationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsTaskCancellationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
