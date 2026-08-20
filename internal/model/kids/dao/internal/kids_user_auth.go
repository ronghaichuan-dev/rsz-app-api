// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsUserAuthDao is the data access object for the table kids_user_auth.
type KidsUserAuthDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  KidsUserAuthColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// KidsUserAuthColumns defines and stores column names for the table kids_user_auth.
type KidsUserAuthColumns struct {
	Id        string // 授权记录ID
	UserId    string // kids用户ID
	Provider  string // 授权服务商：google/apple
	OpenId    string // 服务商开放ID或主体标识
	Email     string // 服务商邮箱
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// kidsUserAuthColumns holds the columns for the table kids_user_auth.
var kidsUserAuthColumns = KidsUserAuthColumns{
	Id:        "id",
	UserId:    "user_id",
	Provider:  "provider",
	OpenId:    "open_id",
	Email:     "email",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewKidsUserAuthDao creates and returns a new DAO object for table data access.
func NewKidsUserAuthDao(handlers ...gdb.ModelHandler) *KidsUserAuthDao {
	return &KidsUserAuthDao{
		group:    "kids",
		table:    "kids_user_auth",
		columns:  kidsUserAuthColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsUserAuthDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsUserAuthDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsUserAuthDao) Columns() KidsUserAuthColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsUserAuthDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsUserAuthDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsUserAuthDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
