// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractTaskAssignmentDao is the data access object for the table kids_contract_task_assignment.
type KidsContractTaskAssignmentDao struct {
	table    string                            // table is the underlying table name of the DAO.
	group    string                            // group is the database configuration group name of the current DAO.
	columns  KidsContractTaskAssignmentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                // handlers for customized model modification.
}

// KidsContractTaskAssignmentColumns defines and stores column names for the table kids_contract_task_assignment.
type KidsContractTaskAssignmentColumns struct {
	Id        string // 主键
	TaskId    string // 合同任务标识
	MemberId  string // 合同成员标识
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// kidsContractTaskAssignmentColumns holds the columns for the table kids_contract_task_assignment.
var kidsContractTaskAssignmentColumns = KidsContractTaskAssignmentColumns{
	Id:        "id",
	TaskId:    "task_id",
	MemberId:  "member_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewKidsContractTaskAssignmentDao creates and returns a new DAO object for table data access.
func NewKidsContractTaskAssignmentDao(handlers ...gdb.ModelHandler) *KidsContractTaskAssignmentDao {
	return &KidsContractTaskAssignmentDao{
		group:    "kids",
		table:    "kids_contract_task_assignment",
		columns:  kidsContractTaskAssignmentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractTaskAssignmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractTaskAssignmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractTaskAssignmentDao) Columns() KidsContractTaskAssignmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractTaskAssignmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractTaskAssignmentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractTaskAssignmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
