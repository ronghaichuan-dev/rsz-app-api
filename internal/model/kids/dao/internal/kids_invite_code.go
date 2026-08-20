// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsInviteCodeDao is the data access object for the table kids_invite_code.
type KidsInviteCodeDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  KidsInviteCodeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// KidsInviteCodeColumns defines and stores column names for the table kids_invite_code.
type KidsInviteCodeColumns struct {
	Id             string // 邀请码ID
	Code           string // 六位邀请码
	CircleId       string // 圈子ID
	InviterUserId  string // 邀请人用户ID
	InviteRole     string // 邀请加入角色：admin/member
	TargetMemberId string // 目标家庭成员ID，0表示不指定
	ExpiredAt      string // 过期时间
	UsedAt         string // 使用时间
	UsedByUserId   string // 使用者用户ID
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// kidsInviteCodeColumns holds the columns for the table kids_invite_code.
var kidsInviteCodeColumns = KidsInviteCodeColumns{
	Id:             "id",
	Code:           "code",
	CircleId:       "circle_id",
	InviterUserId:  "inviter_user_id",
	InviteRole:     "invite_role",
	TargetMemberId: "target_member_id",
	ExpiredAt:      "expired_at",
	UsedAt:         "used_at",
	UsedByUserId:   "used_by_user_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewKidsInviteCodeDao creates and returns a new DAO object for table data access.
func NewKidsInviteCodeDao(handlers ...gdb.ModelHandler) *KidsInviteCodeDao {
	return &KidsInviteCodeDao{
		group:    "kids",
		table:    "kids_invite_code",
		columns:  kidsInviteCodeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsInviteCodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsInviteCodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsInviteCodeDao) Columns() KidsInviteCodeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsInviteCodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsInviteCodeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsInviteCodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
