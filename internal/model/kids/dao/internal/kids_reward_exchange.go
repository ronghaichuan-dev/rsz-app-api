// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsRewardExchangeDao is the data access object for the table kids_reward_exchange.
type KidsRewardExchangeDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsRewardExchangeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsRewardExchangeColumns defines and stores column names for the table kids_reward_exchange.
type KidsRewardExchangeColumns struct {
	Id                             string // 主键
	ExchangeId                     string // 接口兑换标识
	CircleId                       string // 接口圈子标识
	MemberId                       string // 接口成员标识
	MemberNameSnapshot             string // 成员名称快照
	MemberAvatarSnapshot           string // 成员头像快照
	RewardId                       string // 接口奖励标识
	RewardTitleSnapshot            string // 奖励标题快照
	RewardVisualSnapshot           string // 奖励视觉快照
	StarsDeductedSnapshot          string // 扣减星星快照
	RewardRepeatRuleSnapshot       string // 重复规则快照
	RewardCooldownDaysSnapshot     string // 冷却天数快照
	CooldownUntilAtSnapshot        string // 冷却结束快照
	PermanentlyUnavailableSnapshot string // 永久不可兑换快照
	LedgerId                       string // 账本流水标识
	CommitSequence                 string // 提交序列
	ExchangedAt                    string // 兑换时间
	CreatedAt                      string // 创建时间
}

// kidsRewardExchangeColumns holds the columns for the table kids_reward_exchange.
var kidsRewardExchangeColumns = KidsRewardExchangeColumns{
	Id:                             "id",
	ExchangeId:                     "exchange_id",
	CircleId:                       "circle_id",
	MemberId:                       "member_id",
	MemberNameSnapshot:             "member_name_snapshot",
	MemberAvatarSnapshot:           "member_avatar_snapshot",
	RewardId:                       "reward_id",
	RewardTitleSnapshot:            "reward_title_snapshot",
	RewardVisualSnapshot:           "reward_visual_snapshot",
	StarsDeductedSnapshot:          "stars_deducted_snapshot",
	RewardRepeatRuleSnapshot:       "reward_repeat_rule_snapshot",
	RewardCooldownDaysSnapshot:     "reward_cooldown_days_snapshot",
	CooldownUntilAtSnapshot:        "cooldown_until_at_snapshot",
	PermanentlyUnavailableSnapshot: "permanently_unavailable_snapshot",
	LedgerId:                       "ledger_id",
	CommitSequence:                 "commit_sequence",
	ExchangedAt:                    "exchanged_at",
	CreatedAt:                      "created_at",
}

// NewKidsRewardExchangeDao creates and returns a new DAO object for table data access.
func NewKidsRewardExchangeDao(handlers ...gdb.ModelHandler) *KidsRewardExchangeDao {
	return &KidsRewardExchangeDao{
		group:    "kids",
		table:    "kids_reward_exchange",
		columns:  kidsRewardExchangeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsRewardExchangeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsRewardExchangeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsRewardExchangeDao) Columns() KidsRewardExchangeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsRewardExchangeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsRewardExchangeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsRewardExchangeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
