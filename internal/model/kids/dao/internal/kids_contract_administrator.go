// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractAdministratorDao is the data access object for the table kids_contract_administrator.
type KidsContractAdministratorDao struct {
	table    string                           // table is the underlying table name of the DAO.
	group    string                           // group is the database configuration group name of the current DAO.
	columns  KidsContractAdministratorColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler               // handlers for customized model modification.
}

// KidsContractAdministratorColumns defines and stores column names for the table kids_contract_administrator.
type KidsContractAdministratorColumns struct {
	Id              string // 主键
	AdministratorId string // 合同管理员标识
	CircleId        string // 合同圈子标识
	BoundAccountId  string // 绑定账号标识
	DisplayName     string // 显示名称
	Avatar          string // 头像视觉引用
	Role            string // 管理员角色
	Permissions     string // 权限集合
	Status          string // 管理员状态
	Version         string // 版本号
	DeletedAt       string // 删除时间
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// kidsContractAdministratorColumns holds the columns for the table kids_contract_administrator.
var kidsContractAdministratorColumns = KidsContractAdministratorColumns{
	Id:              "id",
	AdministratorId: "administrator_id",
	CircleId:        "circle_id",
	BoundAccountId:  "bound_account_id",
	DisplayName:     "display_name",
	Avatar:          "avatar",
	Role:            "role",
	Permissions:     "permissions",
	Status:          "status",
	Version:         "version",
	DeletedAt:       "deleted_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewKidsContractAdministratorDao creates and returns a new DAO object for table data access.
func NewKidsContractAdministratorDao(handlers ...gdb.ModelHandler) *KidsContractAdministratorDao {
	return &KidsContractAdministratorDao{
		group:    "kids",
		table:    "kids_contract_administrator",
		columns:  kidsContractAdministratorColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractAdministratorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractAdministratorDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractAdministratorDao) Columns() KidsContractAdministratorColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractAdministratorDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractAdministratorDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractAdministratorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
