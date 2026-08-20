// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractLedgerDao is the data access object for the table kids_contract_ledger.
type KidsContractLedgerDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsContractLedgerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsContractLedgerColumns defines and stores column names for the table kids_contract_ledger.
type KidsContractLedgerColumns struct {
	Id                 string // 主键
	LedgerId           string // 合同流水标识
	CircleId           string // 合同圈子标识
	MemberId           string // 合同成员标识
	Source             string // 来源快照
	Delta              string // 余额变化
	Reason             string // 原因
	Actor              string // 操作者快照
	ReversalOfLedgerId string // 原流水标识
	CommitSequence     string // 提交序列
	CreatedAt          string // 创建时间
}

// kidsContractLedgerColumns holds the columns for the table kids_contract_ledger.
var kidsContractLedgerColumns = KidsContractLedgerColumns{
	Id:                 "id",
	LedgerId:           "ledger_id",
	CircleId:           "circle_id",
	MemberId:           "member_id",
	Source:             "source",
	Delta:              "delta",
	Reason:             "reason",
	Actor:              "actor",
	ReversalOfLedgerId: "reversal_of_ledger_id",
	CommitSequence:     "commit_sequence",
	CreatedAt:          "created_at",
}

// NewKidsContractLedgerDao creates and returns a new DAO object for table data access.
func NewKidsContractLedgerDao(handlers ...gdb.ModelHandler) *KidsContractLedgerDao {
	return &KidsContractLedgerDao{
		group:    "kids",
		table:    "kids_contract_ledger",
		columns:  kidsContractLedgerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractLedgerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractLedgerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractLedgerDao) Columns() KidsContractLedgerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractLedgerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractLedgerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractLedgerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
