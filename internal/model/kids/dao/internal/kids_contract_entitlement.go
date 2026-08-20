// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractEntitlementDao is the data access object for the table kids_contract_entitlement.
type KidsContractEntitlementDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  KidsContractEntitlementColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// KidsContractEntitlementColumns defines and stores column names for the table kids_contract_entitlement.
type KidsContractEntitlementColumns struct {
	Id            string // 主键
	EntitlementId string // 权益标识
	AccountId     string // 账号标识
	ProductId     string // 商品标识
	Status        string // 权益状态
	ValidUntilAt  string // 有效截止时间
	VerifiedAt    string // 验证时间
	RevokedAt     string // 撤销时间
	Version       string // 版本号
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// kidsContractEntitlementColumns holds the columns for the table kids_contract_entitlement.
var kidsContractEntitlementColumns = KidsContractEntitlementColumns{
	Id:            "id",
	EntitlementId: "entitlement_id",
	AccountId:     "account_id",
	ProductId:     "product_id",
	Status:        "status",
	ValidUntilAt:  "valid_until_at",
	VerifiedAt:    "verified_at",
	RevokedAt:     "revoked_at",
	Version:       "version",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewKidsContractEntitlementDao creates and returns a new DAO object for table data access.
func NewKidsContractEntitlementDao(handlers ...gdb.ModelHandler) *KidsContractEntitlementDao {
	return &KidsContractEntitlementDao{
		group:    "kids",
		table:    "kids_contract_entitlement",
		columns:  kidsContractEntitlementColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractEntitlementDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractEntitlementDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractEntitlementDao) Columns() KidsContractEntitlementColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractEntitlementDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractEntitlementDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractEntitlementDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
