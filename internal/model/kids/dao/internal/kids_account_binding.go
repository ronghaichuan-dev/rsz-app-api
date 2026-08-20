// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsAccountBindingDao is the data access object for the table kids_account_binding.
type KidsAccountBindingDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsAccountBindingColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsAccountBindingColumns defines and stores column names for the table kids_account_binding.
type KidsAccountBindingColumns struct {
	Id              string // 主键
	BindingId       string // 接口绑定标识
	AccountId       string // 接口账号标识
	Environment     string // 绑定环境
	MigrationPolicy string // 迁移策略
	Version         string // 版本号
	IssuedAt        string // 签发时间
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// kidsAccountBindingColumns holds the columns for the table kids_account_binding.
var kidsAccountBindingColumns = KidsAccountBindingColumns{
	Id:              "id",
	BindingId:       "binding_id",
	AccountId:       "account_id",
	Environment:     "environment",
	MigrationPolicy: "migration_policy",
	Version:         "version",
	IssuedAt:        "issued_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewKidsAccountBindingDao creates and returns a new DAO object for table data access.
func NewKidsAccountBindingDao(handlers ...gdb.ModelHandler) *KidsAccountBindingDao {
	return &KidsAccountBindingDao{
		group:    "kids",
		table:    "kids_account_binding",
		columns:  kidsAccountBindingColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsAccountBindingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsAccountBindingDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsAccountBindingDao) Columns() KidsAccountBindingColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsAccountBindingDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsAccountBindingDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsAccountBindingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
