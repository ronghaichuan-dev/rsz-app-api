// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractMembershipDao is the data access object for the table kids_contract_membership.
type KidsContractMembershipDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  KidsContractMembershipColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// KidsContractMembershipColumns defines and stores column names for the table kids_contract_membership.
type KidsContractMembershipColumns struct {
	Id           string // 主键
	MembershipId string // 合同成员身份标识
	CircleId     string // 合同圈子标识
	AccountId    string // 合同账号标识
	ActorType    string // 授权主体类型
	ActorId      string // 授权主体标识
	Role         string // 成员角色
	Permissions  string // 权限集合
	Status       string // 成员身份状态
	Version      string // 版本号
	DeletedAt    string // 删除时间
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// kidsContractMembershipColumns holds the columns for the table kids_contract_membership.
var kidsContractMembershipColumns = KidsContractMembershipColumns{
	Id:           "id",
	MembershipId: "membership_id",
	CircleId:     "circle_id",
	AccountId:    "account_id",
	ActorType:    "actor_type",
	ActorId:      "actor_id",
	Role:         "role",
	Permissions:  "permissions",
	Status:       "status",
	Version:      "version",
	DeletedAt:    "deleted_at",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewKidsContractMembershipDao creates and returns a new DAO object for table data access.
func NewKidsContractMembershipDao(handlers ...gdb.ModelHandler) *KidsContractMembershipDao {
	return &KidsContractMembershipDao{
		group:    "kids",
		table:    "kids_contract_membership",
		columns:  kidsContractMembershipColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractMembershipDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractMembershipDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractMembershipDao) Columns() KidsContractMembershipColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractMembershipDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractMembershipDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractMembershipDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
