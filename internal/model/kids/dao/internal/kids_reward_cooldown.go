// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsRewardCooldownDao is the data access object for the table kids_reward_cooldown.
type KidsRewardCooldownDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsRewardCooldownColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsRewardCooldownColumns defines and stores column names for the table kids_reward_cooldown.
type KidsRewardCooldownColumns struct {
	Id                     string // 主键
	RewardId               string // 接口奖励标识
	MemberId               string // 接口成员标识
	CooldownUntilAt        string // 冷却结束时间
	LastRedeemedAt         string // 最近兑换时间
	PermanentlyUnavailable string // 是否永久不可兑换
	Version                string // 版本号
	CreatedAt              string // 创建时间
	UpdatedAt              string // 更新时间
}

// kidsRewardCooldownColumns holds the columns for the table kids_reward_cooldown.
var kidsRewardCooldownColumns = KidsRewardCooldownColumns{
	Id:                     "id",
	RewardId:               "reward_id",
	MemberId:               "member_id",
	CooldownUntilAt:        "cooldown_until_at",
	LastRedeemedAt:         "last_redeemed_at",
	PermanentlyUnavailable: "permanently_unavailable",
	Version:                "version",
	CreatedAt:              "created_at",
	UpdatedAt:              "updated_at",
}

// NewKidsRewardCooldownDao creates and returns a new DAO object for table data access.
func NewKidsRewardCooldownDao(handlers ...gdb.ModelHandler) *KidsRewardCooldownDao {
	return &KidsRewardCooldownDao{
		group:    "kids",
		table:    "kids_reward_cooldown",
		columns:  kidsRewardCooldownColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsRewardCooldownDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsRewardCooldownDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsRewardCooldownDao) Columns() KidsRewardCooldownColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsRewardCooldownDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsRewardCooldownDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsRewardCooldownDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
