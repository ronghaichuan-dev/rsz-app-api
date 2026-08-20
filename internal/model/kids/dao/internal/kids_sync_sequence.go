// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsSyncSequenceDao is the data access object for the table kids_sync_sequence.
type KidsSyncSequenceDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  KidsSyncSequenceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// KidsSyncSequenceColumns defines and stores column names for the table kids_sync_sequence.
type KidsSyncSequenceColumns struct {
	Id                 string // 固定序列标识
	NextCommitSequence string // 下一个提交序列
	UpdatedAt          string // 更新时间
}

// kidsSyncSequenceColumns holds the columns for the table kids_sync_sequence.
var kidsSyncSequenceColumns = KidsSyncSequenceColumns{
	Id:                 "id",
	NextCommitSequence: "next_commit_sequence",
	UpdatedAt:          "updated_at",
}

// NewKidsSyncSequenceDao creates and returns a new DAO object for table data access.
func NewKidsSyncSequenceDao(handlers ...gdb.ModelHandler) *KidsSyncSequenceDao {
	return &KidsSyncSequenceDao{
		group:    "kids",
		table:    "kids_sync_sequence",
		columns:  kidsSyncSequenceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsSyncSequenceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsSyncSequenceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsSyncSequenceDao) Columns() KidsSyncSequenceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsSyncSequenceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsSyncSequenceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsSyncSequenceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
