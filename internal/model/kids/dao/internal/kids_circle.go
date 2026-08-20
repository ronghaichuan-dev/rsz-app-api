// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsCircleDao is the data access object for the table kids_circle.
type KidsCircleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KidsCircleColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KidsCircleColumns defines and stores column names for the table kids_circle.
type KidsCircleColumns struct {
	Id          string // 圈子ID
	Name        string // 圈子名称
	Icon        string // 圈子图标标识
	OwnerUserId string // 创建者用户ID
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// kidsCircleColumns holds the columns for the table kids_circle.
var kidsCircleColumns = KidsCircleColumns{
	Id:          "id",
	Name:        "name",
	Icon:        "icon",
	OwnerUserId: "owner_user_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewKidsCircleDao creates and returns a new DAO object for table data access.
func NewKidsCircleDao(handlers ...gdb.ModelHandler) *KidsCircleDao {
	return &KidsCircleDao{
		group:    "kids",
		table:    "kids_circle",
		columns:  kidsCircleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsCircleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsCircleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsCircleDao) Columns() KidsCircleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsCircleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsCircleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsCircleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
