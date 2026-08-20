// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsTaskDefinitionDao is the data access object for the table kids_task_definition.
type KidsTaskDefinitionDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsTaskDefinitionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsTaskDefinitionColumns defines and stores column names for the table kids_task_definition.
type KidsTaskDefinitionColumns struct {
	Id                      string // 主键
	TaskId                  string // 接口任务标识
	CircleId                string // 接口圈子标识
	Title                   string // 任务标题
	Notes                   string // 任务备注
	Emoji                   string // 任务图标
	Stars                   string // 星星数量
	StartDate               string // 系列开始日期
	ZoneId                  string // 时区
	RepeatRule              string // 重复规则
	EndRule                 string // 结束规则
	TimeLimitMinuteOfDay    string // 时限分钟
	ReminderConfig          string // 提醒配置
	PhotoRequired           string // 是否要求图片凭证
	TaskTagId               string // 接口任务标签标识
	SeriesRevision          string // 系列修订号
	FutureEffectiveFromDate string // 当前未来生效日
	Status                  string // 任务状态
	Version                 string // 版本号
	DeletedAt               string // 删除时间
	CreatedAt               string // 创建时间
	UpdatedAt               string // 更新时间
}

// kidsTaskDefinitionColumns holds the columns for the table kids_task_definition.
var kidsTaskDefinitionColumns = KidsTaskDefinitionColumns{
	Id:                      "id",
	TaskId:                  "task_id",
	CircleId:                "circle_id",
	Title:                   "title",
	Notes:                   "notes",
	Emoji:                   "emoji",
	Stars:                   "stars",
	StartDate:               "start_date",
	ZoneId:                  "zone_id",
	RepeatRule:              "repeat_rule",
	EndRule:                 "end_rule",
	TimeLimitMinuteOfDay:    "time_limit_minute_of_day",
	ReminderConfig:          "reminder_config",
	PhotoRequired:           "photo_required",
	TaskTagId:               "task_tag_id",
	SeriesRevision:          "series_revision",
	FutureEffectiveFromDate: "future_effective_from_date",
	Status:                  "status",
	Version:                 "version",
	DeletedAt:               "deleted_at",
	CreatedAt:               "created_at",
	UpdatedAt:               "updated_at",
}

// NewKidsTaskDefinitionDao creates and returns a new DAO object for table data access.
func NewKidsTaskDefinitionDao(handlers ...gdb.ModelHandler) *KidsTaskDefinitionDao {
	return &KidsTaskDefinitionDao{
		group:    "kids",
		table:    "kids_task_definition",
		columns:  kidsTaskDefinitionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsTaskDefinitionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsTaskDefinitionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsTaskDefinitionDao) Columns() KidsTaskDefinitionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsTaskDefinitionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsTaskDefinitionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsTaskDefinitionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
