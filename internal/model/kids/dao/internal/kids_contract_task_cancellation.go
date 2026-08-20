// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractTaskCancellationDao is the data access object for the table kids_contract_task_cancellation.
type KidsContractTaskCancellationDao struct {
	table    string                              // table is the underlying table name of the DAO.
	group    string                              // group is the database configuration group name of the current DAO.
	columns  KidsContractTaskCancellationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                  // handlers for customized model modification.
}

// KidsContractTaskCancellationColumns defines and stores column names for the table kids_contract_task_cancellation.
type KidsContractTaskCancellationColumns struct {
	Id             string // 主键
	CancellationId string // 合同取消标识
	CompletionId   string // 合同完成标识
	ReasonCode     string // 取消原因
	CancelledBy    string // 取消操作者快照
	CancelledAt    string // 取消时间
	CommitSequence string // 提交序列
	CreatedAt      string // 创建时间
}

// kidsContractTaskCancellationColumns holds the columns for the table kids_contract_task_cancellation.
var kidsContractTaskCancellationColumns = KidsContractTaskCancellationColumns{
	Id:             "id",
	CancellationId: "cancellation_id",
	CompletionId:   "completion_id",
	ReasonCode:     "reason_code",
	CancelledBy:    "cancelled_by",
	CancelledAt:    "cancelled_at",
	CommitSequence: "commit_sequence",
	CreatedAt:      "created_at",
}

// NewKidsContractTaskCancellationDao creates and returns a new DAO object for table data access.
func NewKidsContractTaskCancellationDao(handlers ...gdb.ModelHandler) *KidsContractTaskCancellationDao {
	return &KidsContractTaskCancellationDao{
		group:    "kids",
		table:    "kids_contract_task_cancellation",
		columns:  kidsContractTaskCancellationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractTaskCancellationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractTaskCancellationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractTaskCancellationDao) Columns() KidsContractTaskCancellationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractTaskCancellationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractTaskCancellationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractTaskCancellationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
