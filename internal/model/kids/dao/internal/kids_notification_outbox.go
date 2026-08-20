// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsNotificationOutboxDao is the data access object for the table kids_notification_outbox.
type KidsNotificationOutboxDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  KidsNotificationOutboxColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// KidsNotificationOutboxColumns defines and stores column names for the table kids_notification_outbox.
type KidsNotificationOutboxColumns struct {
	Id             string // 主键
	NotificationId string // 接口通知标识
	CircleId       string // 接口圈子标识
	AccountId      string // 接收账号标识
	ExchangeId     string // 接口兑换标识
	EventType      string // 通知事件类型
	Payload        string // 通知载荷
	CommitSequence string // 提交序列
	Status         string // 投递状态
	AttemptCount   string // 已尝试投递次数
	NextAttemptAt  string // 下次允许投递时间
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	Version        string // 版本号
}

// kidsNotificationOutboxColumns holds the columns for the table kids_notification_outbox.
var kidsNotificationOutboxColumns = KidsNotificationOutboxColumns{
	Id:             "id",
	NotificationId: "notification_id",
	CircleId:       "circle_id",
	AccountId:      "account_id",
	ExchangeId:     "exchange_id",
	EventType:      "event_type",
	Payload:        "payload",
	CommitSequence: "commit_sequence",
	Status:         "status",
	AttemptCount:   "attempt_count",
	NextAttemptAt:  "next_attempt_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	Version:        "version",
}

// NewKidsNotificationOutboxDao creates and returns a new DAO object for table data access.
func NewKidsNotificationOutboxDao(handlers ...gdb.ModelHandler) *KidsNotificationOutboxDao {
	return &KidsNotificationOutboxDao{
		group:    "kids",
		table:    "kids_notification_outbox",
		columns:  kidsNotificationOutboxColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsNotificationOutboxDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsNotificationOutboxDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsNotificationOutboxDao) Columns() KidsNotificationOutboxColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsNotificationOutboxDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsNotificationOutboxDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsNotificationOutboxDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
