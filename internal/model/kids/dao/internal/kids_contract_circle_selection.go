// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractCircleSelectionDao is the data access object for the table kids_contract_circle_selection.
type KidsContractCircleSelectionDao struct {
	table    string                             // table is the underlying table name of the DAO.
	group    string                             // group is the database configuration group name of the current DAO.
	columns  KidsContractCircleSelectionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                 // handlers for customized model modification.
}

// KidsContractCircleSelectionColumns defines and stores column names for the table kids_contract_circle_selection.
type KidsContractCircleSelectionColumns struct {
	Id              string // 主键
	SelectionId     string // 合同选择标识
	AccountId       string // 合同账号标识
	CurrentCircleId string // 当前圈子标识
	Version         string // 版本号
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// kidsContractCircleSelectionColumns holds the columns for the table kids_contract_circle_selection.
var kidsContractCircleSelectionColumns = KidsContractCircleSelectionColumns{
	Id:              "id",
	SelectionId:     "selection_id",
	AccountId:       "account_id",
	CurrentCircleId: "current_circle_id",
	Version:         "version",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewKidsContractCircleSelectionDao creates and returns a new DAO object for table data access.
func NewKidsContractCircleSelectionDao(handlers ...gdb.ModelHandler) *KidsContractCircleSelectionDao {
	return &KidsContractCircleSelectionDao{
		group:    "kids",
		table:    "kids_contract_circle_selection",
		columns:  kidsContractCircleSelectionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractCircleSelectionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractCircleSelectionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractCircleSelectionDao) Columns() KidsContractCircleSelectionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractCircleSelectionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractCircleSelectionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractCircleSelectionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
