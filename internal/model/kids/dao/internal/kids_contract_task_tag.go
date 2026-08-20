// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractTaskTagDao is the data access object for the table kids_contract_task_tag.
type KidsContractTaskTagDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  KidsContractTaskTagColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// KidsContractTaskTagColumns defines and stores column names for the table kids_contract_task_tag.
type KidsContractTaskTagColumns struct {
	Id        string // 主键
	TaskTagId string // 合同任务标签标识
	CircleId  string // 合同圈子标识
	Name      string // 标签名称
	Status    string // 标签状态
	Version   string // 版本号
	DeletedAt string // 删除时间
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// kidsContractTaskTagColumns holds the columns for the table kids_contract_task_tag.
var kidsContractTaskTagColumns = KidsContractTaskTagColumns{
	Id:        "id",
	TaskTagId: "task_tag_id",
	CircleId:  "circle_id",
	Name:      "name",
	Status:    "status",
	Version:   "version",
	DeletedAt: "deleted_at",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewKidsContractTaskTagDao creates and returns a new DAO object for table data access.
func NewKidsContractTaskTagDao(handlers ...gdb.ModelHandler) *KidsContractTaskTagDao {
	return &KidsContractTaskTagDao{
		group:    "kids",
		table:    "kids_contract_task_tag",
		columns:  kidsContractTaskTagColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractTaskTagDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractTaskTagDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractTaskTagDao) Columns() KidsContractTaskTagColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractTaskTagDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractTaskTagDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractTaskTagDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
