// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsTaskCompletionDao is the data access object for the table kids_task_completion.
type KidsTaskCompletionDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsTaskCompletionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsTaskCompletionColumns defines and stores column names for the table kids_task_completion.
type KidsTaskCompletionColumns struct {
	Id             string // 主键
	CompletionId   string // 接口完成标识
	CircleId       string // 接口圈子标识
	TaskId         string // 接口任务标识
	MemberId       string // 接口成员标识
	ScheduledDate  string // 预定日
	ZoneId         string // 时区
	ProofAssetId   string // 凭证资产标识
	TitleSnapshot  string // 标题快照
	StarsSnapshot  string // 星星快照
	CompletedBy    string // 完成操作者快照
	CompletedAt    string // 完成时间
	CommitSequence string // 提交序列
	Version        string // 版本号
	CreatedAt      string // 创建时间
}

// kidsTaskCompletionColumns holds the columns for the table kids_task_completion.
var kidsTaskCompletionColumns = KidsTaskCompletionColumns{
	Id:             "id",
	CompletionId:   "completion_id",
	CircleId:       "circle_id",
	TaskId:         "task_id",
	MemberId:       "member_id",
	ScheduledDate:  "scheduled_date",
	ZoneId:         "zone_id",
	ProofAssetId:   "proof_asset_id",
	TitleSnapshot:  "title_snapshot",
	StarsSnapshot:  "stars_snapshot",
	CompletedBy:    "completed_by",
	CompletedAt:    "completed_at",
	CommitSequence: "commit_sequence",
	Version:        "version",
	CreatedAt:      "created_at",
}

// NewKidsTaskCompletionDao creates and returns a new DAO object for table data access.
func NewKidsTaskCompletionDao(handlers ...gdb.ModelHandler) *KidsTaskCompletionDao {
	return &KidsTaskCompletionDao{
		group:    "kids",
		table:    "kids_task_completion",
		columns:  kidsTaskCompletionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsTaskCompletionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsTaskCompletionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsTaskCompletionDao) Columns() KidsTaskCompletionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsTaskCompletionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsTaskCompletionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsTaskCompletionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
