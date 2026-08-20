// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsTaskTagDao is the data access object for the table kids_task_tag.
type KidsTaskTagDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KidsTaskTagColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KidsTaskTagColumns defines and stores column names for the table kids_task_tag.
type KidsTaskTagColumns struct {
	Id        string // 任务标签ID
	Name      string // 标签名称
	Color     string // 标签颜色
	SortOrder string // 排序值
	DeletedAt string // 删除时间
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// kidsTaskTagColumns holds the columns for the table kids_task_tag.
var kidsTaskTagColumns = KidsTaskTagColumns{
	Id:        "id",
	Name:      "name",
	Color:     "color",
	SortOrder: "sort_order",
	DeletedAt: "deleted_at",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewKidsTaskTagDao creates and returns a new DAO object for table data access.
func NewKidsTaskTagDao(handlers ...gdb.ModelHandler) *KidsTaskTagDao {
	return &KidsTaskTagDao{
		group:    "kids",
		table:    "kids_task_tag",
		columns:  kidsTaskTagColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsTaskTagDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsTaskTagDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsTaskTagDao) Columns() KidsTaskTagColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsTaskTagDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsTaskTagDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsTaskTagDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
