// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractAccountBindingDao is the data access object for the table kids_contract_account_binding.
type KidsContractAccountBindingDao struct {
	table    string                            // table is the underlying table name of the DAO.
	group    string                            // group is the database configuration group name of the current DAO.
	columns  KidsContractAccountBindingColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                // handlers for customized model modification.
}

// KidsContractAccountBindingColumns defines and stores column names for the table kids_contract_account_binding.
type KidsContractAccountBindingColumns struct {
	Id              string // 主键
	BindingId       string // 合同绑定标识
	AccountId       string // 合同账号标识
	Environment     string // 绑定环境
	MigrationPolicy string // 迁移策略
	Version         string // 版本号
	IssuedAt        string // 签发时间
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// kidsContractAccountBindingColumns holds the columns for the table kids_contract_account_binding.
var kidsContractAccountBindingColumns = KidsContractAccountBindingColumns{
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

// NewKidsContractAccountBindingDao creates and returns a new DAO object for table data access.
func NewKidsContractAccountBindingDao(handlers ...gdb.ModelHandler) *KidsContractAccountBindingDao {
	return &KidsContractAccountBindingDao{
		group:    "kids",
		table:    "kids_contract_account_binding",
		columns:  kidsContractAccountBindingColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractAccountBindingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractAccountBindingDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractAccountBindingDao) Columns() KidsContractAccountBindingColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractAccountBindingDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractAccountBindingDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractAccountBindingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
