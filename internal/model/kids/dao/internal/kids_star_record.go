// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsStarRecordDao is the data access object for the table kids_star_record.
type KidsStarRecordDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  KidsStarRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// KidsStarRecordColumns defines and stores column names for the table kids_star_record.
type KidsStarRecordColumns struct {
	Id           string // 星星流水ID
	KidId        string // 儿童成员ID
	ChangeAmount string // 星星变动数量
	Balance      string // 变动后余额
	RecordType   string // 流水类型：task/reward/adjustment
	Title        string // 流水标题
	Remark       string // 流水备注
	CreatedAt    string // 创建时间
}

// kidsStarRecordColumns holds the columns for the table kids_star_record.
var kidsStarRecordColumns = KidsStarRecordColumns{
	Id:           "id",
	KidId:        "kid_id",
	ChangeAmount: "change_amount",
	Balance:      "balance",
	RecordType:   "record_type",
	Title:        "title",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewKidsStarRecordDao creates and returns a new DAO object for table data access.
func NewKidsStarRecordDao(handlers ...gdb.ModelHandler) *KidsStarRecordDao {
	return &KidsStarRecordDao{
		group:    "kids",
		table:    "kids_star_record",
		columns:  kidsStarRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsStarRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsStarRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsStarRecordDao) Columns() KidsStarRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsStarRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsStarRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsStarRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
