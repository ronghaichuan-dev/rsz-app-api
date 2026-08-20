// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsStarBalanceDao is the data access object for the table kids_star_balance.
type KidsStarBalanceDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  KidsStarBalanceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// KidsStarBalanceColumns defines and stores column names for the table kids_star_balance.
type KidsStarBalanceColumns struct {
	Id                   string // 主键
	CircleId             string // 接口圈子标识
	MemberId             string // 接口成员标识
	Balance              string // 星星余额
	Version              string // 版本号
	SourceCommitId       string // 来源提交标识
	SourceCommitSequence string // 来源提交序列
	UpdatedAt            string // 更新时间
}

// kidsStarBalanceColumns holds the columns for the table kids_star_balance.
var kidsStarBalanceColumns = KidsStarBalanceColumns{
	Id:                   "id",
	CircleId:             "circle_id",
	MemberId:             "member_id",
	Balance:              "balance",
	Version:              "version",
	SourceCommitId:       "source_commit_id",
	SourceCommitSequence: "source_commit_sequence",
	UpdatedAt:            "updated_at",
}

// NewKidsStarBalanceDao creates and returns a new DAO object for table data access.
func NewKidsStarBalanceDao(handlers ...gdb.ModelHandler) *KidsStarBalanceDao {
	return &KidsStarBalanceDao{
		group:    "kids",
		table:    "kids_star_balance",
		columns:  kidsStarBalanceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsStarBalanceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsStarBalanceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsStarBalanceDao) Columns() KidsStarBalanceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsStarBalanceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsStarBalanceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsStarBalanceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
