// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsFeedbackDao is the data access object for the table kids_feedback.
type KidsFeedbackDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  KidsFeedbackColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// KidsFeedbackColumns defines and stores column names for the table kids_feedback.
type KidsFeedbackColumns struct {
	Id                    string // 主键
	FeedbackId            string // 反馈标识
	AccountId             string // 账号标识
	Category              string // 反馈分类
	Content               string // 反馈内容
	ContactType           string // 联系类型
	Contact               string // 联系方式
	PrivacyConsentVersion string // 隐私同意版本
	AttachmentAssetIds    string // 附件资产标识
	Version               string // 版本号
	CreatedAt             string // 创建时间
}

// kidsFeedbackColumns holds the columns for the table kids_feedback.
var kidsFeedbackColumns = KidsFeedbackColumns{
	Id:                    "id",
	FeedbackId:            "feedback_id",
	AccountId:             "account_id",
	Category:              "category",
	Content:               "content",
	ContactType:           "contact_type",
	Contact:               "contact",
	PrivacyConsentVersion: "privacy_consent_version",
	AttachmentAssetIds:    "attachment_asset_ids",
	Version:               "version",
	CreatedAt:             "created_at",
}

// NewKidsFeedbackDao creates and returns a new DAO object for table data access.
func NewKidsFeedbackDao(handlers ...gdb.ModelHandler) *KidsFeedbackDao {
	return &KidsFeedbackDao{
		group:    "kids",
		table:    "kids_feedback",
		columns:  kidsFeedbackColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsFeedbackDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsFeedbackDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsFeedbackDao) Columns() KidsFeedbackColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsFeedbackDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsFeedbackDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsFeedbackDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
