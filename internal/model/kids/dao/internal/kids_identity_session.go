// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsIdentitySessionDao is the data access object for the table kids_identity_session.
type KidsIdentitySessionDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  KidsIdentitySessionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// KidsIdentitySessionColumns defines and stores column names for the table kids_identity_session.
type KidsIdentitySessionColumns struct {
	Id                    string // 主键
	SessionId             string // 接口会话标识
	AccountId             string // 接口账号标识
	PrincipalKind         string // 主体类型
	AccessTokenHash       string // 访问令牌摘要
	RefreshTokenHash      string // 刷新令牌摘要
	GuestUpgradeGrantHash string // 游客升级凭据摘要
	ExpiresAt             string // 访问令牌到期时间
	RefreshExpiresAt      string // 刷新令牌到期时间
	RevokedAt             string // 撤销时间
	CreatedAt             string // 创建时间
	UpdatedAt             string // 更新时间
}

// kidsIdentitySessionColumns holds the columns for the table kids_identity_session.
var kidsIdentitySessionColumns = KidsIdentitySessionColumns{
	Id:                    "id",
	SessionId:             "session_id",
	AccountId:             "account_id",
	PrincipalKind:         "principal_kind",
	AccessTokenHash:       "access_token_hash",
	RefreshTokenHash:      "refresh_token_hash",
	GuestUpgradeGrantHash: "guest_upgrade_grant_hash",
	ExpiresAt:             "expires_at",
	RefreshExpiresAt:      "refresh_expires_at",
	RevokedAt:             "revoked_at",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewKidsIdentitySessionDao creates and returns a new DAO object for table data access.
func NewKidsIdentitySessionDao(handlers ...gdb.ModelHandler) *KidsIdentitySessionDao {
	return &KidsIdentitySessionDao{
		group:    "kids",
		table:    "kids_identity_session",
		columns:  kidsIdentitySessionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsIdentitySessionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsIdentitySessionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsIdentitySessionDao) Columns() KidsIdentitySessionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsIdentitySessionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsIdentitySessionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsIdentitySessionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
