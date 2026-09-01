package kids

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// GetMemberBalancesV1 从接口余额投影读取指定成员的当前星星余额。
func (s *sKids) GetMemberBalancesV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1GetMemberBalances)
}

// ListStarTransactionsV1 查询接口账本中指定成员的追加式星星流水。
func (s *sKids) ListStarTransactionsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1ListStarTransactions)
}

// v1GetMemberBalances 校验当前账号的圈子成员身份后，返回指定成员的规范余额投影。
func (s *sKids) v1GetMemberBalances(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	memberIDs := in.Query["member_id"]
	items := make([]map[string]any, 0, len(memberIDs))
	cursor := ""
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, txErr := v1ReadMembershipTx(ctx, tx, in.PrincipalID, circleID); txErr != nil {
			return v1MemberBalanceDependencyError(ctx, "membership", txErr)
		}
		members, txErr := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).
			Where("circle_id", circleID).Where("member_id IN(?)", memberIDs).Where("status", "active").All()
		if txErr != nil {
			return v1MemberBalanceDependencyError(ctx, "member_read_store", txErr)
		}
		membersByID := make(map[string]gdb.Record, len(members))
		for _, member := range members {
			membersByID[member["member_id"].String()] = member
		}
		balances, txErr := tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).
			Where("circle_id", circleID).Where("member_id IN(?)", memberIDs).All()
		if txErr != nil {
			return v1MemberBalanceDependencyError(ctx, "balance_snapshot_store", txErr)
		}
		balancesByID := make(map[string]gdb.Record, len(balances))
		for _, balance := range balances {
			balancesByID[balance["member_id"].String()] = balance
		}
		for _, memberID := range memberIDs {
			if _, ok := membersByID[memberID]; !ok {
				return v1Error(404, "NOT_FOUND", false, "member is missing")
			}
			balance, ok := balancesByID[memberID]
			if !ok {
				return v1Error(409, "AUDIT_INCONSISTENT", false, "canonical member balance is missing")
			}
			items = append(items, v1BalanceProjectionFromRecord(balance))
		}
		var cursorErr error
		cursor, cursorErr = v1LatestCursorTx(ctx, tx)
		if cursorErr != nil {
			return v1MemberBalanceDependencyError(ctx, "commit_snapshot_store", cursorErr)
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"items": items, "snapshot_cursor": cursor}, cursor, nil
}

// v1ReadMembershipTx 在同一只读快照中确认当前账号仍拥有 active membership，不获取更新锁。
func v1ReadMembershipTx(ctx context.Context, tx gdb.TX, accountID, circleID string) (gdb.Record, error) {
	row, err := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).
		Where("account_id", accountID).Where("circle_id", circleID).Where("status", "active").One()
	if err != nil {
		return nil, err
	}
	if row.IsEmpty() {
		return nil, v1Error(403, "FORBIDDEN", false, "active circle membership is required")
	}
	return row, nil
}

// v1MemberBalanceDependencyError 保留稳定业务错误，并把真实余额读取依赖故障映射为可重试 UNAVAILABLE。
func v1MemberBalanceDependencyError(ctx context.Context, dependency string, err error) error {
	var protocolErr *v1.V1Error
	if errors.As(err, &protocolErr) {
		return err
	}
	requestID, traceID, operationID := "", "", "getMemberBalances"
	if request := ghttp.RequestFromCtx(ctx); request != nil {
		requestID = request.Header.Get(v1.V1RequestIDHeader)
		traceID = request.GetCtxVar(consts.CtxTraceIDKey).String()
		operationID = request.GetCtxVar(consts.CtxV1OperationIDKey).String()
	}
	g.Log().Errorf(ctx, "event=kids_member_balance_unavailable operation_id=%s request_id=%s trace_id=%s dependency=%s error=%+v", operationID, requestID, traceID, dependency, err)
	retryAfterMs := int64(1000)
	return &v1.V1Error{Status: 503, Code: "UNAVAILABLE", Retryable: true, RetryAfterMs: &retryAfterMs, Message: "member balance snapshot dependency is unavailable"}
}

// v1ListStarTransactions 返回经过成员权限、时间范围和来源类型筛选后的接口账本分页。
func (s *sKids) v1ListStarTransactions(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	limit := v1QueryFirst(in.Query, "limit")
	offset, err := v1PageOffset(v1QueryFirst(in.Query, "cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "ledger cursor is invalid")
	}
	var pageSize int
	if _, err = fmt.Sscan(limit, &pageSize); err != nil || pageSize < 1 || pageSize > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "ledger limit is invalid")
	}
	sourceType := v1QueryFirst(in.Query, "source_type")
	items := make([]map[string]any, 0, pageSize)
	cursor := ""
	hasMore := false
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, txErr := v1ReadMembershipTx(ctx, tx, in.PrincipalID, circleID); txErr != nil {
			return v1StarTransactionDependencyError(ctx, "membership", txErr)
		}
		model := tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).
			Fields("ledger_id,circle_id,member_id,source,delta,reason,actor,reversal_of_ledger_id,commit_sequence,created_at").
			Where("circle_id", circleID).
			Where("member_id", v1QueryFirst(in.Query, "member_id"))
		if startText := v1QueryFirst(in.Query, "start_at_ms"); startText != "" {
			var startAt int64
			if _, txErr := fmt.Sscan(startText, &startAt); txErr != nil {
				return v1Error(422, "VALIDATION_FAILED", false, "ledger start time is invalid")
			}
			model = model.Where("created_at >= ?", time.UnixMilli(startAt))
		}
		if endText := v1QueryFirst(in.Query, "end_at_ms"); endText != "" {
			var endAt int64
			if _, txErr := fmt.Sscan(endText, &endAt); txErr != nil {
				return v1Error(422, "VALIDATION_FAILED", false, "ledger end time is invalid")
			}
			model = model.Where("created_at < ?", time.UnixMilli(endAt))
		}
		if sourceType != "" {
			model = model.Where("JSON_UNQUOTE(JSON_EXTRACT(source, '$.source_type')) = ?", sourceType)
		}
		rows, txErr := model.Order("created_at DESC,id DESC").Limit(offset, pageSize+1).All()
		if txErr != nil {
			return v1StarTransactionDependencyError(ctx, "ledger_read_store", txErr)
		}
		hasMore = len(rows) > pageSize
		if hasMore {
			rows = rows[:pageSize]
		}
		for _, row := range rows {
			items = append(items, v1LedgerRecordProjection(row))
		}
		var cursorErr error
		cursor, cursorErr = v1LatestCursorTx(ctx, tx)
		if cursorErr != nil {
			return v1StarTransactionDependencyError(ctx, "commit_snapshot_store", cursorErr)
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	next := any(nil)
	if hasMore {
		next = v1PageCursor(offset + len(items))
	}
	page := map[string]any{"items": items, "next_cursor": next, "has_more": hasMore, "snapshot_cursor": cursor}
	if validationErr := v1.ValidateV1ResponseData(in.OperationID, page); validationErr != nil {
		v1LogStarTransactionProjectionFailure(ctx, in, len(items), validationErr)
		return nil, "", v1Error(409, "AUDIT_INCONSISTENT", false, "star transaction projection is inconsistent")
	}
	return page, cursor, nil
}

// v1StarTransactionDependencyError 保留稳定业务错误，并将真实流水读取依赖故障映射为可重试 UNAVAILABLE。
func v1StarTransactionDependencyError(ctx context.Context, dependency string, err error) error {
	var protocolErr *v1.V1Error
	if errors.As(err, &protocolErr) {
		return err
	}
	requestID, traceID, operationID := "", "", "listStarTransactions"
	if request := ghttp.RequestFromCtx(ctx); request != nil {
		requestID = request.Header.Get(v1.V1RequestIDHeader)
		traceID = request.GetCtxVar(consts.CtxTraceIDKey).String()
		operationID = request.GetCtxVar(consts.CtxV1OperationIDKey).String()
	}
	g.Log().Errorf(ctx, "event=kids_star_transaction_unavailable operation_id=%s request_id=%s trace_id=%s dependency=%s error=%+v", operationID, requestID, traceID, dependency, err)
	retryAfterMs := int64(1000)
	return &v1.V1Error{Status: 503, Code: "UNAVAILABLE", Retryable: true, RetryAfterMs: &retryAfterMs, Message: "star transaction ledger dependency is unavailable"}
}

// v1LogStarTransactionProjectionFailure 记录不含请求正文的流水投影合同故障，供 requestId 与 traceId 关联排查。
func v1LogStarTransactionProjectionFailure(ctx context.Context, in v1.V1OperationInput, itemCount int, err error) {
	traceID := ""
	if request := ghttp.RequestFromCtx(ctx); request != nil {
		traceID = request.GetCtxVar(consts.CtxTraceIDKey).String()
	}
	g.Log().Errorf(ctx, "event=kids_star_transaction_projection_invalid operation_id=%s request_id=%s trace_id=%s circle_id=%s member_id=%s item_count=%d error=%v", in.OperationID, in.RequestID, traceID, in.PathParameters["circle_id"], v1QueryFirst(in.Query, "member_id"), itemCount, err)
}

// GetStarBalance 从星星流水表读取指定儿童的当前星星余额。
func (s *sKids) GetStarBalance(ctx context.Context, in v1.StarBalanceInput) (*v1.StarBalanceOutput, error) {
	balance, err := currentStarBalance(ctx, nil, in.KidId)
	if err != nil {
		return nil, err
	}
	return &v1.StarBalanceOutput{KidId: in.KidId, Balance: balance}, nil
}

// AdjustStars 在事务中写入手动调整流水，并返回调整后的星星余额。
func (s *sKids) AdjustStars(ctx context.Context, in v1.StarAdjustInput) (*v1.StarAdjustOutput, error) {
	if in.Amount == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "amount cannot be 0")
	}
	if strings.TrimSpace(in.Reason) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reason is required")
	}

	var record v1.StarRecord
	var balance int
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var err error
		balance, err = addStarRecord(ctx, tx, in.KidId, in.Amount, v1.StarRecordTypeAdjustment, commoni18n.T(ctx, "Balance Adjustment"), strings.TrimSpace(in.Reason))
		if err != nil {
			return err
		}
		record, err = latestStarRecord(ctx, tx, in.KidId)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.StarAdjustOutput{KidId: in.KidId, Balance: balance, Record: record}, nil
}

// ListStarRecords 从数据库按儿童、类型和时间范围筛选星星流水。
func (s *sKids) ListStarRecords(ctx context.Context, in v1.StarRecordListInput) (*v1.StarRecordListOutput, error) {
	if !validStarRecordType(in.Type) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "unsupported star record type")
	}
	model := utils.KidsDB(ctx).Model(consts.KidsStarRecordTable).Ctx(ctx).Where("kid_id", in.KidId)
	if strings.TrimSpace(in.Type) != "" {
		model = model.Where("record_type", strings.TrimSpace(in.Type))
	}
	if strings.TrimSpace(in.From) != "" {
		model = model.WhereGTE("created_at", strings.TrimSpace(in.From)+" 00:00:00")
	}
	if strings.TrimSpace(in.To) != "" {
		model = model.WhereLTE("created_at", strings.TrimSpace(in.To)+" 23:59:59")
	}
	records, err := model.OrderDesc("id").All()
	if err != nil {
		return nil, err
	}
	out := &v1.StarRecordListOutput{}
	for _, record := range records {
		out.List = append(out.List, starRecordFromDB(record))
	}
	return out, nil
}

// validStarRecordType 校验星星流水类型筛选，空值表示不过滤。
func validStarRecordType(recordType string) bool {
	switch strings.TrimSpace(recordType) {
	case "", v1.StarRecordTypeTask, v1.StarRecordTypeReward, v1.StarRecordTypeAdjustment:
		return true
	default:
		return false
	}
}

// addStarRecord 根据上一条流水计算余额，并持久化新的星星流水。
func addStarRecord(ctx context.Context, tx gdb.TX, kidId uint64, change int, recordType, title, remark string) (int, error) {
	balance, err := currentStarBalance(ctx, tx, kidId)
	if err != nil {
		return 0, err
	}
	balance += change
	if balance < 0 {
		balance = 0
	}
	_, err = tx.Model(consts.KidsStarRecordTable).Ctx(ctx).Data(map[string]any{
		"kid_id":        kidId,
		"change_amount": change,
		"balance":       balance,
		"record_type":   recordType,
		"title":         title,
		"remark":        strings.TrimSpace(remark),
	}).Insert()
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// currentStarBalance 查询某个儿童最新一条星星流水中的余额。
func currentStarBalance(ctx context.Context, tx gdb.TX, kidId uint64) (int, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsStarRecordTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsStarRecordTable).Ctx(ctx)
	}
	value, err := model.Where("kid_id", kidId).OrderDesc("id").Value("balance")
	if err != nil || value == nil {
		return 0, err
	}
	return value.Int(), nil
}

// latestStarRecord 查询某个儿童最新一条星星流水。
func latestStarRecord(ctx context.Context, tx gdb.TX, kidId uint64) (v1.StarRecord, error) {
	record, err := tx.Model(consts.KidsStarRecordTable).Ctx(ctx).Where("kid_id", kidId).OrderDesc("id").One()
	if err != nil || record.IsEmpty() {
		return v1.StarRecord{}, err
	}
	return starRecordFromDB(record), nil
}

// starRecordFromDB 将数据库星星流水记录转换为接口响应结构。
func starRecordFromDB(record gdb.Record) v1.StarRecord {
	return v1.StarRecord{
		Id:        record["id"].Uint64(),
		KidId:     record["kid_id"].Uint64(),
		Change:    record["change_amount"].Int(),
		Balance:   record["balance"].Int(),
		Type:      record["record_type"].String(),
		Title:     record["title"].String(),
		Remark:    record["remark"].String(),
		CreatedAt: utils.ParseDBTime(record["created_at"].Val()),
	}
}
