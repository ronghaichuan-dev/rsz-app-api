// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsTaskDao is the data access object for the table kids_task.
type KidsTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KidsTaskColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KidsTaskColumns defines and stores column names for the table kids_task.
type KidsTaskColumns struct {
	Id                    string // 任务ID
	Title                 string // 任务标题
	Icon                  string // 图标标识
	Note                  string // 任务说明
	Star                  string // 星星奖励数量
	TaskDate              string // 计划日期
	CompletionMode        string // 完成模式：single/rotation/anyone/everyone
	RepeatRule            string // 重复规则
	RepeatEndType         string // 重复结束类型：never/date/count
	RepeatEndDate         string // 重复结束日期
	RepeatEndCount        string // 重复结束次数
	TimeLimitType         string // 时间限制类型：all_day/range
	TimeLimitStart        string // 开始时间，格式HH:mm
	TimeLimitEnd          string // 结束时间，格式HH:mm
	ReminderType          string // 提醒类型：none/at_time/before_start
	ReminderAt            string // 提醒时间，格式HH:mm
	ReminderOffsetMinutes string // 提前提醒分钟数
	NeedPhotoProof        string // 是否需要照片凭证
	TagId                 string // 标签ID
	Completed             string // 是否已完成
	CompletedBy           string // 完成任务的儿童成员ID
	PhotoUrl              string // 照片凭证地址
	CompletedAt           string // 完成时间
	DeletedAt             string // 删除时间
	CreatedAt             string // 创建时间
	UpdatedAt             string // 更新时间
}

// kidsTaskColumns holds the columns for the table kids_task.
var kidsTaskColumns = KidsTaskColumns{
	Id:                    "id",
	Title:                 "title",
	Icon:                  "icon",
	Note:                  "note",
	Star:                  "star",
	TaskDate:              "task_date",
	CompletionMode:        "completion_mode",
	RepeatRule:            "repeat_rule",
	RepeatEndType:         "repeat_end_type",
	RepeatEndDate:         "repeat_end_date",
	RepeatEndCount:        "repeat_end_count",
	TimeLimitType:         "time_limit_type",
	TimeLimitStart:        "time_limit_start",
	TimeLimitEnd:          "time_limit_end",
	ReminderType:          "reminder_type",
	ReminderAt:            "reminder_at",
	ReminderOffsetMinutes: "reminder_offset_minutes",
	NeedPhotoProof:        "need_photo_proof",
	TagId:                 "tag_id",
	Completed:             "completed",
	CompletedBy:           "completed_by",
	PhotoUrl:              "photo_url",
	CompletedAt:           "completed_at",
	DeletedAt:             "deleted_at",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewKidsTaskDao creates and returns a new DAO object for table data access.
func NewKidsTaskDao(handlers ...gdb.ModelHandler) *KidsTaskDao {
	return &KidsTaskDao{
		group:    "kids",
		table:    "kids_task",
		columns:  kidsTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsTaskDao) Columns() KidsTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
