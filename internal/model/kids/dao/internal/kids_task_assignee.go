// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsTaskAssigneeDao is the data access object for the table kids_task_assignee.
type KidsTaskAssigneeDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  KidsTaskAssigneeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// KidsTaskAssigneeColumns defines and stores column names for the table kids_task_assignee.
type KidsTaskAssigneeColumns struct {
	Id            string // 任务分配ID
	TaskId        string // 任务ID
	KidId         string // 儿童成员ID
	AssigneeOrder string // 分配顺序，用于轮流模式
	Completed     string // 该儿童是否已完成
	PhotoUrl      string // 该儿童照片凭证地址
	CompletedAt   string // 该儿童完成时间
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// kidsTaskAssigneeColumns holds the columns for the table kids_task_assignee.
var kidsTaskAssigneeColumns = KidsTaskAssigneeColumns{
	Id:            "id",
	TaskId:        "task_id",
	KidId:         "kid_id",
	AssigneeOrder: "assignee_order",
	Completed:     "completed",
	PhotoUrl:      "photo_url",
	CompletedAt:   "completed_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewKidsTaskAssigneeDao creates and returns a new DAO object for table data access.
func NewKidsTaskAssigneeDao(handlers ...gdb.ModelHandler) *KidsTaskAssigneeDao {
	return &KidsTaskAssigneeDao{
		group:    "kids",
		table:    "kids_task_assignee",
		columns:  kidsTaskAssigneeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsTaskAssigneeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsTaskAssigneeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsTaskAssigneeDao) Columns() KidsTaskAssigneeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsTaskAssigneeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsTaskAssigneeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsTaskAssigneeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
