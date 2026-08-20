// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsInvitationDao is the data access object for the table kids_invitation.
type KidsInvitationDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  KidsInvitationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// KidsInvitationColumns defines and stores column names for the table kids_invitation.
type KidsInvitationColumns struct {
	Id                    string // 主键
	InviteId              string // 接口邀请标识
	CircleId              string // 接口圈子标识
	TargetRole            string // 目标角色
	TargetAdministratorId string // 目标管理员标识
	TargetMemberId        string // 目标成员标识
	PermissionScope       string // 邀请权限范围
	CodeHash              string // 邀请码摘要
	SingleUse             string // 是否单次使用
	Generation            string // 邀请码代次
	Status                string // 邀请状态
	ExpiresAt             string // 过期时间
	UsedAt                string // 使用时间
	RevokedAt             string // 撤销时间
	Version               string // 版本号
	CreatedByAccountId    string // 创建账号标识
	CreatedAt             string // 创建时间
	UpdatedAt             string // 更新时间
}

// kidsInvitationColumns holds the columns for the table kids_invitation.
var kidsInvitationColumns = KidsInvitationColumns{
	Id:                    "id",
	InviteId:              "invite_id",
	CircleId:              "circle_id",
	TargetRole:            "target_role",
	TargetAdministratorId: "target_administrator_id",
	TargetMemberId:        "target_member_id",
	PermissionScope:       "permission_scope",
	CodeHash:              "code_hash",
	SingleUse:             "single_use",
	Generation:            "generation",
	Status:                "status",
	ExpiresAt:             "expires_at",
	UsedAt:                "used_at",
	RevokedAt:             "revoked_at",
	Version:               "version",
	CreatedByAccountId:    "created_by_account_id",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewKidsInvitationDao creates and returns a new DAO object for table data access.
func NewKidsInvitationDao(handlers ...gdb.ModelHandler) *KidsInvitationDao {
	return &KidsInvitationDao{
		group:    "kids",
		table:    "kids_invitation",
		columns:  kidsInvitationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsInvitationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsInvitationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsInvitationDao) Columns() KidsInvitationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsInvitationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsInvitationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsInvitationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
