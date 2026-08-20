// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractTaskOccurrenceDao is the data access object for the table kids_contract_task_occurrence.
type KidsContractTaskOccurrenceDao struct {
	table    string                            // table is the underlying table name of the DAO.
	group    string                            // group is the database configuration group name of the current DAO.
	columns  KidsContractTaskOccurrenceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                // handlers for customized model modification.
}

// KidsContractTaskOccurrenceColumns defines and stores column names for the table kids_contract_task_occurrence.
type KidsContractTaskOccurrenceColumns struct {
	Id                    string // 内部主键
	CircleId              string // 合同圈子标识
	TaskId                string // 合同任务标识
	MemberId              string // 合同成员标识
	ScheduledDate         string // 预定日
	ZoneId                string // 时区
	DefinitionRevision    string // 定义修订号
	TitleSnapshot         string // 标题快照
	NotesSnapshot         string // 备注快照
	EmojiSnapshot         string // 图标快照
	StarsSnapshot         string // 星星快照
	PhotoRequiredSnapshot string // 图片要求快照
	TaskTagIdSnapshot     string // 标签快照
	State                 string // occurrence 状态
	CompletionId          string // 完成事实标识
	Version               string // 版本号
	CreatedAt             string // 创建时间
	UpdatedAt             string // 更新时间
}

// kidsContractTaskOccurrenceColumns holds the columns for the table kids_contract_task_occurrence.
var kidsContractTaskOccurrenceColumns = KidsContractTaskOccurrenceColumns{
	Id:                    "id",
	CircleId:              "circle_id",
	TaskId:                "task_id",
	MemberId:              "member_id",
	ScheduledDate:         "scheduled_date",
	ZoneId:                "zone_id",
	DefinitionRevision:    "definition_revision",
	TitleSnapshot:         "title_snapshot",
	NotesSnapshot:         "notes_snapshot",
	EmojiSnapshot:         "emoji_snapshot",
	StarsSnapshot:         "stars_snapshot",
	PhotoRequiredSnapshot: "photo_required_snapshot",
	TaskTagIdSnapshot:     "task_tag_id_snapshot",
	State:                 "state",
	CompletionId:          "completion_id",
	Version:               "version",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewKidsContractTaskOccurrenceDao creates and returns a new DAO object for table data access.
func NewKidsContractTaskOccurrenceDao(handlers ...gdb.ModelHandler) *KidsContractTaskOccurrenceDao {
	return &KidsContractTaskOccurrenceDao{
		group:    "kids",
		table:    "kids_contract_task_occurrence",
		columns:  kidsContractTaskOccurrenceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractTaskOccurrenceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractTaskOccurrenceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractTaskOccurrenceDao) Columns() KidsContractTaskOccurrenceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractTaskOccurrenceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractTaskOccurrenceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractTaskOccurrenceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
