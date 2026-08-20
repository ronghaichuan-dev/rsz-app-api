// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsUserDao is the data access object for the table kids_user.
type KidsUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KidsUserColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KidsUserColumns defines and stores column names for the table kids_user.
type KidsUserColumns struct {
	Id        string // 用户ID
	DeviceId  string // 最近登录设备ID
	Provider  string // 当前登录方式：guest/google/apple
	Email     string // 授权服务商返回的邮箱
	Nickname  string // 昵称
	Avatar    string // 头像地址
	IsGuest   string // 是否游客账号
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// kidsUserColumns holds the columns for the table kids_user.
var kidsUserColumns = KidsUserColumns{
	Id:        "id",
	DeviceId:  "device_id",
	Provider:  "provider",
	Email:     "email",
	Nickname:  "nickname",
	Avatar:    "avatar",
	IsGuest:   "is_guest",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewKidsUserDao creates and returns a new DAO object for table data access.
func NewKidsUserDao(handlers ...gdb.ModelHandler) *KidsUserDao {
	return &KidsUserDao{
		group:    "kids",
		table:    "kids_user",
		columns:  kidsUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsUserDao) Columns() KidsUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
