// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsNotificationDao is the data access object for the table kids_notification.
type KidsNotificationDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  KidsNotificationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// KidsNotificationColumns defines and stores column names for the table kids_notification.
type KidsNotificationColumns struct {
	Id               string // 通知ID
	MemberId         string // 目标成员ID，0表示全家庭
	NotificationType string // 通知类型
	Title            string // 通知标题
	Content          string // 通知内容
	IsRead           string // 是否已读
	CreatedAt        string // 创建时间
	UpdatedAt        string // 更新时间
}

// kidsNotificationColumns holds the columns for the table kids_notification.
var kidsNotificationColumns = KidsNotificationColumns{
	Id:               "id",
	MemberId:         "member_id",
	NotificationType: "notification_type",
	Title:            "title",
	Content:          "content",
	IsRead:           "is_read",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

// NewKidsNotificationDao creates and returns a new DAO object for table data access.
func NewKidsNotificationDao(handlers ...gdb.ModelHandler) *KidsNotificationDao {
	return &KidsNotificationDao{
		group:    "kids",
		table:    "kids_notification",
		columns:  kidsNotificationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsNotificationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsNotificationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsNotificationDao) Columns() KidsNotificationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsNotificationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsNotificationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsNotificationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
