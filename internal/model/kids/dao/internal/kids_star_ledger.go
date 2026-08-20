// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsStarLedgerDao is the data access object for the table kids_star_ledger.
type KidsStarLedgerDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  KidsStarLedgerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// KidsStarLedgerColumns defines and stores column names for the table kids_star_ledger.
type KidsStarLedgerColumns struct {
	Id                 string // 主键
	LedgerId           string // 接口流水标识
	CircleId           string // 接口圈子标识
	MemberId           string // 接口成员标识
	Source             string // 来源快照
	Delta              string // 余额变化
	Reason             string // 原因
	Actor              string // 操作者快照
	ReversalOfLedgerId string // 原流水标识
	CommitSequence     string // 提交序列
	CreatedAt          string // 创建时间
}

// kidsStarLedgerColumns holds the columns for the table kids_star_ledger.
var kidsStarLedgerColumns = KidsStarLedgerColumns{
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

// NewKidsStarLedgerDao creates and returns a new DAO object for table data access.
func NewKidsStarLedgerDao(handlers ...gdb.ModelHandler) *KidsStarLedgerDao {
	return &KidsStarLedgerDao{
		group:    "kids",
		table:    "kids_star_ledger",
		columns:  kidsStarLedgerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsStarLedgerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsStarLedgerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsStarLedgerDao) Columns() KidsStarLedgerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsStarLedgerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsStarLedgerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsStarLedgerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
