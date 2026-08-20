// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KidsContractIdempotencyDao is the data access object for the table kids_contract_idempotency.
type KidsContractIdempotencyDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  KidsContractIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// KidsContractIdempotencyColumns defines and stores column names for the table kids_contract_idempotency.
type KidsContractIdempotencyColumns struct {
	Id                   string // 主键
	PrincipalScope       string // 主体和成员作用域
	IdempotencyKey       string // 幂等键
	OperationId          string // 合同操作标识
	RouteFingerprint     string // 路由摘要
	BodyFingerprint      string // 请求体摘要
	ResponseStatus       string // 首次响应状态码
	ResponseBody         string // 首次规范响应
	ResponseChangeCursor string // 首次响应变更游标
	ResponseEtag         string // 首次响应实体标签
	CreatedAt            string // 创建时间
	UpdatedAt            string // 更新时间
}

// kidsContractIdempotencyColumns holds the columns for the table kids_contract_idempotency.
var kidsContractIdempotencyColumns = KidsContractIdempotencyColumns{
	Id:                   "id",
	PrincipalScope:       "principal_scope",
	IdempotencyKey:       "idempotency_key",
	OperationId:          "operation_id",
	RouteFingerprint:     "route_fingerprint",
	BodyFingerprint:      "body_fingerprint",
	ResponseStatus:       "response_status",
	ResponseBody:         "response_body",
	ResponseChangeCursor: "response_change_cursor",
	ResponseEtag:         "response_etag",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewKidsContractIdempotencyDao creates and returns a new DAO object for table data access.
func NewKidsContractIdempotencyDao(handlers ...gdb.ModelHandler) *KidsContractIdempotencyDao {
	return &KidsContractIdempotencyDao{
		group:    "kids",
		table:    "kids_contract_idempotency",
		columns:  kidsContractIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KidsContractIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KidsContractIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KidsContractIdempotencyDao) Columns() KidsContractIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KidsContractIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KidsContractIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KidsContractIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
