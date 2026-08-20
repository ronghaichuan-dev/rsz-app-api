// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsRewardAssigneeDao is the data access object for the table kids_reward_assignee.
type KidsRewardAssigneeDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsRewardAssigneeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsRewardAssigneeColumns defines and stores column names for the table kids_reward_assignee.
type KidsRewardAssigneeColumns struct {
	Id        string // 奖励指派ID
	RewardId  string // 奖励ID
	KidId     string // 儿童成员ID
	CreatedAt string // 创建时间
}

// kidsRewardAssigneeColumns holds the columns for the table kids_reward_assignee.
var kidsRewardAssigneeColumns = KidsRewardAssigneeColumns{
	Id:        "id",
	RewardId:  "reward_id",
	KidId:     "kid_id",
	CreatedAt: "created_at",
}

// NewKidsRewardAssigneeDao creates and returns a new DAO object for table data access.
func NewKidsRewardAssigneeDao(handlers ...gdb.ModelHandler) *KidsRewardAssigneeDao {
	return &KidsRewardAssigneeDao{
		group:    "kids",
		table:    "kids_reward_assignee",
		columns:  kidsRewardAssigneeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsRewardAssigneeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsRewardAssigneeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsRewardAssigneeDao) Columns() KidsRewardAssigneeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsRewardAssigneeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsRewardAssigneeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsRewardAssigneeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
