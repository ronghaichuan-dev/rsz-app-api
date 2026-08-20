// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractMemberDao is the data access object for the table kids_contract_member.
type KidsContractMemberDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  KidsContractMemberColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// KidsContractMemberColumns defines and stores column names for the table kids_contract_member.
type KidsContractMemberColumns struct {
	Id             string // 主键
	MemberId       string // 合同成员标识
	CircleId       string // 合同圈子标识
	BoundAccountId string // 绑定账号标识
	DisplayName    string // 显示名称
	Gender         string // 性别
	Avatar         string // 头像视觉引用
	Status         string // 成员状态
	Version        string // 版本号
	DeletedAt      string // 删除时间
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// kidsContractMemberColumns holds the columns for the table kids_contract_member.
var kidsContractMemberColumns = KidsContractMemberColumns{
	Id:             "id",
	MemberId:       "member_id",
	CircleId:       "circle_id",
	BoundAccountId: "bound_account_id",
	DisplayName:    "display_name",
	Gender:         "gender",
	Avatar:         "avatar",
	Status:         "status",
	Version:        "version",
	DeletedAt:      "deleted_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewKidsContractMemberDao creates and returns a new DAO object for table data access.
func NewKidsContractMemberDao(handlers ...gdb.ModelHandler) *KidsContractMemberDao {
	return &KidsContractMemberDao{
		group:    "kids",
		table:    "kids_contract_member",
		columns:  kidsContractMemberColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractMemberDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractMemberDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractMemberDao) Columns() KidsContractMemberColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractMemberDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractMemberDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractMemberDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
