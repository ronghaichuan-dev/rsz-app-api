// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsRewardRedeemRecordDao is the data access object for the table kids_reward_redeem_record.
type KidsRewardRedeemRecordDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  KidsRewardRedeemRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// KidsRewardRedeemRecordColumns defines and stores column names for the table kids_reward_redeem_record.
type KidsRewardRedeemRecordColumns struct {
	Id        string // 奖励兑换记录ID
	CircleId  string // 圈子ID
	RewardId  string // 奖励ID
	KidId     string // 儿童成员ID
	UserId    string // 兑换用户ID
	Title     string // 奖励标题
	Icon      string // 图标标识
	ImageUrl  string // 奖励图片地址
	StarCost  string // 消耗星星数量
	Remark    string // 兑换备注
	CreatedAt string // 创建时间
}

// kidsRewardRedeemRecordColumns holds the columns for the table kids_reward_redeem_record.
var kidsRewardRedeemRecordColumns = KidsRewardRedeemRecordColumns{
	Id:        "id",
	CircleId:  "circle_id",
	RewardId:  "reward_id",
	KidId:     "kid_id",
	UserId:    "user_id",
	Title:     "title",
	Icon:      "icon",
	ImageUrl:  "image_url",
	StarCost:  "star_cost",
	Remark:    "remark",
	CreatedAt: "created_at",
}

// NewKidsRewardRedeemRecordDao creates and returns a new DAO object for table data access.
func NewKidsRewardRedeemRecordDao(handlers ...gdb.ModelHandler) *KidsRewardRedeemRecordDao {
	return &KidsRewardRedeemRecordDao{
		group:    "kids",
		table:    "kids_reward_redeem_record",
		columns:  kidsRewardRedeemRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsRewardRedeemRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsRewardRedeemRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsRewardRedeemRecordDao) Columns() KidsRewardRedeemRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsRewardRedeemRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsRewardRedeemRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsRewardRedeemRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
