// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsDeviceNotificationDao is the data access object for the table kids_device_notification.
type KidsDeviceNotificationDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  KidsDeviceNotificationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// KidsDeviceNotificationColumns defines and stores column names for the table kids_device_notification.
type KidsDeviceNotificationColumns struct {
	Id            string // 设备通知ID
	UserId        string // 用户ID
	DeviceId      string // 设备ID
	Platform      string // 平台：ios/android/web
	DeviceToken   string // 推送设备令牌
	Authorized    string // 是否已授权通知
	TaskEnabled   string // 是否开启任务提醒
	RewardEnabled string // 是否开启奖励提醒
	MemberEnabled string // 是否开启成员提醒
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// kidsDeviceNotificationColumns holds the columns for the table kids_device_notification.
var kidsDeviceNotificationColumns = KidsDeviceNotificationColumns{
	Id:            "id",
	UserId:        "user_id",
	DeviceId:      "device_id",
	Platform:      "platform",
	DeviceToken:   "device_token",
	Authorized:    "authorized",
	TaskEnabled:   "task_enabled",
	RewardEnabled: "reward_enabled",
	MemberEnabled: "member_enabled",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewKidsDeviceNotificationDao creates and returns a new DAO object for table data access.
func NewKidsDeviceNotificationDao(handlers ...gdb.ModelHandler) *KidsDeviceNotificationDao {
	return &KidsDeviceNotificationDao{
		group:    "kids",
		table:    "kids_device_notification",
		columns:  kidsDeviceNotificationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsDeviceNotificationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsDeviceNotificationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsDeviceNotificationDao) Columns() KidsDeviceNotificationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsDeviceNotificationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsDeviceNotificationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsDeviceNotificationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
