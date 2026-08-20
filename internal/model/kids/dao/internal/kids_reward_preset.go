// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsRewardPresetDao is the data access object for the table kids_reward_preset.
type KidsRewardPresetDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  KidsRewardPresetColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// KidsRewardPresetColumns defines and stores column names for the table kids_reward_preset.
type KidsRewardPresetColumns struct {
	Id          string // 奖励预设ID
	Title       string // 预设奖励标题
	Icon        string // 图标标识
	ImageUrl    string // 奖励图片地址
	StarCost    string // 默认所需星星数量
	Description string // 预设描述
	Enabled     string // 是否启用预设
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// kidsRewardPresetColumns holds the columns for the table kids_reward_preset.
var kidsRewardPresetColumns = KidsRewardPresetColumns{
	Id:          "id",
	Title:       "title",
	Icon:        "icon",
	ImageUrl:    "image_url",
	StarCost:    "star_cost",
	Description: "description",
	Enabled:     "enabled",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewKidsRewardPresetDao creates and returns a new DAO object for table data access.
func NewKidsRewardPresetDao(handlers ...gdb.ModelHandler) *KidsRewardPresetDao {
	return &KidsRewardPresetDao{
		group:    "kids",
		table:    "kids_reward_preset",
		columns:  kidsRewardPresetColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsRewardPresetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsRewardPresetDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsRewardPresetDao) Columns() KidsRewardPresetColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsRewardPresetDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsRewardPresetDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsRewardPresetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
