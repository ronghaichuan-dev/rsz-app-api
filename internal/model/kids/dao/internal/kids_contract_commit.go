// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractCommitDao is the data access object for the table kids_contract_commit.
type KidsContractCommitDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsContractCommitColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsContractCommitColumns defines and stores column names for the table kids_contract_commit.
type KidsContractCommitColumns struct {
	Id             string // 主键
	CommitId       string // 合同提交标识
	CircleId       string // 圈子标识
	CommitSequence string // 单调提交序列
	ChangePayload  string // 完整变更集合
	CreatedAt      string // 创建时间
}

// kidsContractCommitColumns holds the columns for the table kids_contract_commit.
var kidsContractCommitColumns = KidsContractCommitColumns{
	Id:             "id",
	CommitId:       "commit_id",
	CircleId:       "circle_id",
	CommitSequence: "commit_sequence",
	ChangePayload:  "change_payload",
	CreatedAt:      "created_at",
}

// NewKidsContractCommitDao creates and returns a new DAO object for table data access.
func NewKidsContractCommitDao(handlers ...gdb.ModelHandler) *KidsContractCommitDao {
	return &KidsContractCommitDao{
		group:    "kids",
		table:    "kids_contract_commit",
		columns:  kidsContractCommitColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractCommitDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractCommitDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractCommitDao) Columns() KidsContractCommitColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractCommitDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractCommitDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractCommitDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
