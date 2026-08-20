// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractCircleDao is the data access object for the table kids_contract_circle.
type KidsContractCircleDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsContractCircleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsContractCircleColumns defines and stores column names for the table kids_contract_circle.
type KidsContractCircleColumns struct {
	Id                   string // 主键
	CircleId             string // 合同圈子标识
	Name                 string // 圈子名称
	Icon                 string // 圈子视觉引用
	OwnerAdministratorId string // 所有者管理员标识
	Status               string // 圈子状态
	Version              string // 版本号
	DeletedAt            string // 删除时间
	CreatedAt            string // 创建时间
	UpdatedAt            string // 更新时间
}

// kidsContractCircleColumns holds the columns for the table kids_contract_circle.
var kidsContractCircleColumns = KidsContractCircleColumns{
	Id:                   "id",
	CircleId:             "circle_id",
	Name:                 "name",
	Icon:                 "icon",
	OwnerAdministratorId: "owner_administrator_id",
	Status:               "status",
	Version:              "version",
	DeletedAt:            "deleted_at",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewKidsContractCircleDao creates and returns a new DAO object for table data access.
func NewKidsContractCircleDao(handlers ...gdb.ModelHandler) *KidsContractCircleDao {
	return &KidsContractCircleDao{
		group:    "kids",
		table:    "kids_contract_circle",
		columns:  kidsContractCircleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractCircleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractCircleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractCircleDao) Columns() KidsContractCircleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractCircleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractCircleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractCircleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
