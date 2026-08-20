// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsRewardDao is the data access object for the table kids_reward.
type KidsRewardDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KidsRewardColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KidsRewardColumns defines and stores column names for the table kids_reward.
type KidsRewardColumns struct {
	Id                 string // 奖励ID
	CircleId           string // 所属圈子ID
	Title              string // 奖励标题
	Icon               string // 图标标识
	ImageUrl           string // 奖励图片地址
	StarCost           string // 所需星星数量
	Stock              string // 可用库存，-1表示不限量
	Description        string // 奖励描述
	RepeatRule         string // 重复兑换规则：none/daily/weekly/monthly/custom
	RepeatIntervalDays string // 自定义重复兑换间隔天数
	CreatedAt          string // 创建时间
	UpdatedAt          string // 更新时间
	DeletedAt          string // 删除时间
}

// kidsRewardColumns holds the columns for the table kids_reward.
var kidsRewardColumns = KidsRewardColumns{
	Id:                 "id",
	CircleId:           "circle_id",
	Title:              "title",
	Icon:               "icon",
	ImageUrl:           "image_url",
	StarCost:           "star_cost",
	Stock:              "stock",
	Description:        "description",
	RepeatRule:         "repeat_rule",
	RepeatIntervalDays: "repeat_interval_days",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
}

// NewKidsRewardDao creates and returns a new DAO object for table data access.
func NewKidsRewardDao(handlers ...gdb.ModelHandler) *KidsRewardDao {
	return &KidsRewardDao{
		group:    "kids",
		table:    "kids_reward",
		columns:  kidsRewardColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsRewardDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsRewardDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsRewardDao) Columns() KidsRewardColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsRewardDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsRewardDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsRewardDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
