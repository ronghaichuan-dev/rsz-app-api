package kids

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

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
