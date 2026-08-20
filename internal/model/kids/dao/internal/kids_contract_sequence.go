// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractSequenceDao is the data access object for the table kids_contract_sequence.
type KidsContractSequenceDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  KidsContractSequenceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// KidsContractSequenceColumns defines and stores column names for the table kids_contract_sequence.
type KidsContractSequenceColumns struct {
	Id                 string // 固定序列标识
	NextCommitSequence string // 下一个提交序列
	UpdatedAt          string // 更新时间
}

// kidsContractSequenceColumns holds the columns for the table kids_contract_sequence.
var kidsContractSequenceColumns = KidsContractSequenceColumns{
	Id:                 "id",
	NextCommitSequence: "next_commit_sequence",
	UpdatedAt:          "updated_at",
}

// NewKidsContractSequenceDao creates and returns a new DAO object for table data access.
func NewKidsContractSequenceDao(handlers ...gdb.ModelHandler) *KidsContractSequenceDao {
	return &KidsContractSequenceDao{
		group:    "kids",
		table:    "kids_contract_sequence",
		columns:  kidsContractSequenceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractSequenceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractSequenceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractSequenceDao) Columns() KidsContractSequenceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractSequenceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractSequenceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractSequenceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
