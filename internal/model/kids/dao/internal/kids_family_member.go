// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsFamilyMemberDao is the data access object for the table kids_family_member.
type KidsFamilyMemberDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  KidsFamilyMemberColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// KidsFamilyMemberColumns defines and stores column names for the table kids_family_member.
type KidsFamilyMemberColumns struct {
	Id          string // 家庭成员ID
	CircleId    string // 所属圈子ID
	Name        string // 显示名称
	Gender      string // 性别：male/female
	Avatar      string // 头像地址或预设标识
	AvatarStyle string // 虚拟形象风格标识
	Relation    string // 家庭关系
	Owner       string // 是否家庭拥有者
	BindUserId  string // 绑定用户ID
	BoundAt     string // 绑定时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// kidsFamilyMemberColumns holds the columns for the table kids_family_member.
var kidsFamilyMemberColumns = KidsFamilyMemberColumns{
	Id:          "id",
	CircleId:    "circle_id",
	Name:        "name",
	Gender:      "gender",
	Avatar:      "avatar",
	AvatarStyle: "avatar_style",
	Relation:    "relation",
	Owner:       "owner",
	BindUserId:  "bind_user_id",
	BoundAt:     "bound_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewKidsFamilyMemberDao creates and returns a new DAO object for table data access.
func NewKidsFamilyMemberDao(handlers ...gdb.ModelHandler) *KidsFamilyMemberDao {
	return &KidsFamilyMemberDao{
		group:    "kids",
		table:    "kids_family_member",
		columns:  kidsFamilyMemberColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsFamilyMemberDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsFamilyMemberDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsFamilyMemberDao) Columns() KidsFamilyMemberColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsFamilyMemberDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsFamilyMemberDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsFamilyMemberDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
