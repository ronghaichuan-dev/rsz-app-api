package kids

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// ListRewardPresets 从数据库读取启用的奖励预设，并支持按关键字搜索。
func (s *sKids) ListRewardPresets(ctx context.Context, in v1.RewardPresetListInput) (*v1.RewardPresetListOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsRewardPresetTable).Ctx(ctx).Where("enabled", 1)
	if strings.TrimSpace(in.Keyword) != "" {
		model = model.WhereLike("title", "%"+strings.TrimSpace(in.Keyword)+"%")
	}
	records, err := model.OrderAsc("id").All()
	if err != nil {
		return nil, err
	}
	out := &v1.RewardPresetListOutput{}
	for _, record := range records {
		out.List = append(out.List, rewardPresetFromDB(record))
	}
	return out, nil
}

// ListRewards 从数据库读取奖励列表，支持按群组、儿童和兑换进度筛选。
func (s *sKids) ListRewards(ctx context.Context, in v1.RewardListInput) (*v1.RewardListOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsRewardTable).Ctx(ctx).Where("deleted_at IS NULL")
	if in.CircleId > 0 {
		model = model.Where("circle_id", in.CircleId)
	}
	records, err := model.OrderAsc("id").All()
	if err != nil {
		return nil, err
	}
	out := &v1.RewardListOutput{}
	for _, record := range records {
		item, err := rewardFromDB(ctx, nil, record, in.KidId)
		if err != nil {
			return nil, err
		}
		if in.KidId > 0 && !rewardVisibleForKid(item, in.KidId) {
			continue
		}
		out.List = append(out.List, item)
	}
	return out, nil
}

// GetReward 查询单个未删除奖励详情。
func (s *sKids) GetReward(ctx context.Context, in v1.RewardDetailInput) (*v1.RewardDetailOutput, error) {
	record, err := utils.KidsDB(ctx).Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "reward not found")
	}
	reward, err := rewardFromDB(ctx, nil, record, in.KidId)
	if err != nil {
		return nil, err
	}
	return &v1.RewardDetailOutput{Reward: reward}, nil
}

// CreateReward 校验奖励入参并持久化创建奖励。
func (s *sKids) CreateReward(ctx context.Context, in v1.RewardCreateInput) (*v1.RewardCreateOutput, error) {
	if err := validateRewardPayload(in.Title, in.StarCost, in.RepeatRule, in.RepeatIntervalDays); err != nil {
		return nil, err
	}
	if in.CircleId > 0 {
		if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, in.CircleId); err != nil || !ok {
			if err != nil {
				return nil, err
			}
			return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can manage reward")
		}
	}
	stock := in.Stock
	if stock == 0 {
		stock = -1
	}
	var reward v1.RewardItem
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model(consts.KidsRewardTable).Ctx(ctx).Data(map[string]any{
			"circle_id":            in.CircleId,
			"title":                strings.TrimSpace(in.Title),
			"icon":                 strings.TrimSpace(in.Icon),
			"image_url":            strings.TrimSpace(in.ImageUrl),
			"star_cost":            in.StarCost,
			"stock":                stock,
			"description":          strings.TrimSpace(in.Description),
			"repeat_rule":          normalizedRewardRepeatRule(in.RepeatRule),
			"repeat_interval_days": in.RepeatIntervalDays,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		if err = replaceRewardAssignees(ctx, tx, uint64(id), in.KidIds); err != nil {
			return err
		}
		record, err := tx.Model(consts.KidsRewardTable).Ctx(ctx).Where("id", id).One()
		if err != nil {
			return err
		}
		reward, err = rewardFromDB(ctx, tx, record, 0)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.RewardCreateOutput{Reward: reward}, nil
}

// UpdateReward 校验奖励入参并更新奖励配置和指派儿童。
func (s *sKids) UpdateReward(ctx context.Context, in v1.RewardUpdateInput) (*v1.RewardUpdateOutput, error) {
	if err := validateRewardPayload(in.Title, in.StarCost, in.RepeatRule, in.RepeatIntervalDays); err != nil {
		return nil, err
	}
	if in.CircleId > 0 {
		if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, in.CircleId); err != nil || !ok {
			if err != nil {
				return nil, err
			}
			return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can manage reward")
		}
	}
	stock := in.Stock
	if stock == 0 {
		stock = -1
	}
	var reward v1.RewardItem
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").Data(map[string]any{
			"circle_id":            in.CircleId,
			"title":                strings.TrimSpace(in.Title),
			"icon":                 strings.TrimSpace(in.Icon),
			"image_url":            strings.TrimSpace(in.ImageUrl),
			"star_cost":            in.StarCost,
			"stock":                stock,
			"description":          strings.TrimSpace(in.Description),
			"repeat_rule":          normalizedRewardRepeatRule(in.RepeatRule),
			"repeat_interval_days": in.RepeatIntervalDays,
		}).Update()
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "reward not found")
		}
		if err = replaceRewardAssignees(ctx, tx, in.Id, in.KidIds); err != nil {
			return err
		}
		record, err := tx.Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).One()
		if err != nil {
			return err
		}
		reward, err = rewardFromDB(ctx, tx, record, 0)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.RewardUpdateOutput{Reward: reward}, nil
}

// DeleteReward 软删除奖励并保留兑换历史。
func (s *sKids) DeleteReward(ctx context.Context, in v1.RewardDeleteInput) (*v1.RewardDeleteOutput, error) {
	record, err := utils.KidsDB(ctx).Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "reward not found")
	}
	if record["circle_id"].Uint64() > 0 {
		if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, record["circle_id"].Uint64()); err != nil || !ok {
			if err != nil {
				return nil, err
			}
			return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can manage reward")
		}
	}
	result, err := utils.KidsDB(ctx).Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").Data(map[string]any{"deleted_at": time.Now().Format(consts.MySQLTimeLayout)}).Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "reward not found")
	}
	return &v1.RewardDeleteOutput{Id: in.Id}, nil
}

// RedeemReward 在事务中完成奖励兑换、扣减库存、校验重复周期并写入星星流水。
func (s *sKids) RedeemReward(ctx context.Context, in v1.RewardRedeemInput) (*v1.RewardRedeemOutput, error) {
	var reward v1.RewardItem
	var balance int
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		record, err := tx.Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").One()
		if err != nil {
			return err
		}
		if record.IsEmpty() {
			return gerror.NewCode(gcode.CodeNotFound, "reward not found")
		}
		currentReward, err := rewardFromDB(ctx, tx, record, in.KidId)
		if err != nil {
			return err
		}
		if !rewardVisibleForKid(currentReward, in.KidId) {
			return gerror.NewCode(gcode.CodeNotAuthorized, "reward is not assigned to this kid")
		}
		if record["stock"].Int() == 0 {
			return gerror.NewCode(gcode.CodeInvalidOperation, "reward is out of stock")
		}
		nextAt, err := nextRewardRedeemAt(ctx, tx, in.Id, in.KidId, record["repeat_rule"].String(), record["repeat_interval_days"].Int())
		if err != nil {
			return err
		}
		if nextAt > time.Now().Unix() {
			return gerror.NewCode(gcode.CodeInvalidOperation, "reward repeat interval has not expired")
		}
		currentBalance, err := currentStarBalance(ctx, tx, in.KidId)
		if err != nil {
			return err
		}
		if currentBalance < record["star_cost"].Int() {
			return gerror.NewCode(gcode.CodeInvalidOperation, "not enough stars")
		}
		if record["stock"].Int() > 0 {
			if _, err = tx.Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).Data(map[string]any{"stock": record["stock"].Int() - 1}).Update(); err != nil {
				return err
			}
		}
		balance, err = addStarRecord(ctx, tx, in.KidId, -record["star_cost"].Int(), v1.StarRecordTypeReward, commoni18n.Tf(ctx, "star.reward_redeemed", record["title"].String()), strings.TrimSpace(in.Remark))
		if err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsRewardRedeemRecordTable).Ctx(ctx).Data(map[string]any{
			"circle_id": record["circle_id"].Uint64(),
			"reward_id": in.Id,
			"kid_id":    in.KidId,
			"title":     record["title"].String(),
			"icon":      record["icon"].String(),
			"image_url": record["image_url"].String(),
			"star_cost": record["star_cost"].Int(),
			"remark":    strings.TrimSpace(in.Remark),
		}).Insert(); err != nil {
			return err
		}
		if _, err = createNotification(ctx, tx, in.KidId, "reward", commoni18n.T(ctx, "Reward Redeemed!"), commoni18n.T(ctx, "You have received a reward! Tap here to see it!")); err != nil {
			return err
		}
		record, err = tx.Model(consts.KidsRewardTable).Ctx(ctx).Where("id", in.Id).One()
		if err != nil {
			return err
		}
		reward, err = rewardFromDB(ctx, tx, record, in.KidId)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.RewardRedeemOutput{Reward: reward, StarBalance: balance}, nil
}

// ListRewardRedeemRecords 从数据库查询奖励兑换历史，支持月份和日期范围筛选。
func (s *sKids) ListRewardRedeemRecords(ctx context.Context, in v1.RewardRecordListInput) (*v1.RewardRecordListOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsRewardRedeemRecordTable+" rr").Ctx(ctx).
		LeftJoin(consts.KidsFamilyMemberTable+" fm", "fm.id = rr.kid_id").
		Fields("rr.*, fm.name AS kid_name")
	if in.CircleId > 0 {
		model = model.Where("rr.circle_id", in.CircleId)
	}
	if in.KidId > 0 {
		model = model.Where("rr.kid_id", in.KidId)
	}
	from, to := rewardRecordRange(in)
	if from != "" {
		model = model.WhereGTE("rr.created_at", from)
	}
	if to != "" {
		model = model.WhereLTE("rr.created_at", to)
	}
	records, err := model.OrderDesc("rr.id").All()
	if err != nil {
		return nil, err
	}
	out := &v1.RewardRecordListOutput{Total: len(records)}
	for _, record := range records {
		out.List = append(out.List, rewardRedeemRecordFromDB(record))
	}
	return out, nil
}

// rewardPresetFromDB 将数据库奖励预设转换为接口结构。
func rewardPresetFromDB(record gdb.Record) v1.RewardPreset {
	return v1.RewardPreset{Id: record["id"].Uint64(), Title: record["title"].String(), Icon: record["icon"].String(), ImageUrl: record["image_url"].String(), StarCost: record["star_cost"].Int(), Description: record["description"].String()}
}

// rewardFromDB 将数据库奖励记录转换为接口响应结构。
func rewardFromDB(ctx context.Context, tx gdb.TX, record gdb.Record, kidId uint64) (v1.RewardItem, error) {
	kidIds, err := rewardAssigneeIds(ctx, tx, record["id"].Uint64())
	if err != nil {
		return v1.RewardItem{}, err
	}
	item := v1.RewardItem{Id: record["id"].Uint64(), CircleId: record["circle_id"].Uint64(), Title: record["title"].String(), Icon: record["icon"].String(), ImageUrl: record["image_url"].String(), StarCost: record["star_cost"].Int(), Stock: record["stock"].Int(), Description: record["description"].String(), RepeatRule: record["repeat_rule"].String(), RepeatIntervalDays: record["repeat_interval_days"].Int(), KidIds: kidIds, CanRedeem: true}
	if kidId > 0 {
		balance, err := currentStarBalance(ctx, tx, kidId)
		if err != nil {
			return v1.RewardItem{}, err
		}
		if item.StarCost > 0 {
			item.Progress = math.Min(float64(balance)/float64(item.StarCost), 1)
		}
		item.NextRedeemAt, err = nextRewardRedeemAt(ctx, tx, item.Id, kidId, item.RepeatRule, item.RepeatIntervalDays)
		if err != nil {
			return v1.RewardItem{}, err
		}
		item.CanRedeem = balance >= item.StarCost && (item.Stock != 0) && rewardVisibleForKid(item, kidId) && item.NextRedeemAt <= time.Now().Unix()
	}
	return item, nil
}

// rewardRedeemRecordFromDB 将数据库奖励兑换记录转换为接口结构。
func rewardRedeemRecordFromDB(record gdb.Record) v1.RewardRedeemRecord {
	return v1.RewardRedeemRecord{Id: record["id"].Uint64(), CircleId: record["circle_id"].Uint64(), RewardId: record["reward_id"].Uint64(), KidId: record["kid_id"].Uint64(), KidName: record["kid_name"].String(), UserId: record["user_id"].Uint64(), Title: record["title"].String(), Icon: record["icon"].String(), ImageUrl: record["image_url"].String(), StarCost: record["star_cost"].Int(), Remark: record["remark"].String(), CreatedAt: utils.ParseDBTime(record["created_at"].Val())}
}

// validateRewardPayload 校验奖励核心字段和重复兑换规则。
func validateRewardPayload(title string, starCost int, repeatRule string, repeatIntervalDays int) error {
	if strings.TrimSpace(title) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "title is required")
	}
	if starCost <= 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "starCost must be greater than 0")
	}
	if !validRewardRepeatRule(repeatRule) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported repeatRule")
	}
	if normalizedRewardRepeatRule(repeatRule) == consts.RewardRepeatCustom && repeatIntervalDays <= 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "repeat interval days is required")
	}
	if repeatIntervalDays < 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "repeat interval days cannot be negative")
	}
	return nil
}

// replaceRewardAssignees 重建奖励指派儿童关系。
func replaceRewardAssignees(ctx context.Context, tx gdb.TX, rewardId uint64, kidIds []uint64) error {
	if _, err := tx.Model(consts.KidsRewardAssigneeTable).Ctx(ctx).Where("reward_id", rewardId).Delete(); err != nil {
		return err
	}
	for _, kidId := range uniqueUint64s(kidIds) {
		if _, err := tx.Model(consts.KidsRewardAssigneeTable).Ctx(ctx).Data(map[string]any{"reward_id": rewardId, "kid_id": kidId}).Insert(); err != nil {
			return err
		}
	}
	return nil
}

// rewardAssigneeIds 查询奖励指派儿童列表。
func rewardAssigneeIds(ctx context.Context, tx gdb.TX, rewardId uint64) ([]uint64, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsRewardAssigneeTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsRewardAssigneeTable).Ctx(ctx)
	}
	records, err := model.Where("reward_id", rewardId).OrderAsc("id").All()
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record["kid_id"].Uint64())
	}
	return ids, nil
}

// rewardVisibleForKid 判断奖励是否对指定儿童可见。
func rewardVisibleForKid(reward v1.RewardItem, kidId uint64) bool {
	if kidId == 0 || len(reward.KidIds) == 0 {
		return true
	}
	for _, id := range reward.KidIds {
		if id == kidId {
			return true
		}
	}
	return false
}

// nextRewardRedeemAt 根据重复规则计算下一次可兑换时间。
func nextRewardRedeemAt(ctx context.Context, tx gdb.TX, rewardId, kidId uint64, repeatRule string, intervalDays int) (int64, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsRewardRedeemRecordTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsRewardRedeemRecordTable).Ctx(ctx)
	}
	record, err := model.Where("reward_id", rewardId).Where("kid_id", kidId).OrderDesc("id").One()
	if err != nil || record.IsEmpty() {
		return 0, err
	}
	last := utils.ParseDBTime(record["created_at"].Val())
	if last == 0 {
		return 0, nil
	}
	lastTime := time.Unix(last, 0)
	switch normalizedRewardRepeatRule(repeatRule) {
	case consts.RewardRepeatDaily:
		return lastTime.AddDate(0, 0, 1).Unix(), nil
	case consts.RewardRepeatWeekly:
		return lastTime.AddDate(0, 0, 7).Unix(), nil
	case consts.RewardRepeatMonthly:
		return lastTime.AddDate(0, 1, 0).Unix(), nil
	case consts.RewardRepeatCustom:
		if intervalDays <= 0 {
			intervalDays = 1
		}
		return lastTime.AddDate(0, 0, intervalDays).Unix(), nil
	default:
		return lastTime.AddDate(100, 0, 0).Unix(), nil
	}
}

// validRewardRepeatRule 校验奖励重复兑换规则，空值允许走默认不重复。
func validRewardRepeatRule(rule string) bool {
	switch strings.TrimSpace(rule) {
	case "", consts.RewardRepeatNone, consts.RewardRepeatDaily, consts.RewardRepeatWeekly, consts.RewardRepeatMonthly, consts.RewardRepeatCustom:
		return true
	default:
		return false
	}
}

// normalizedRewardRepeatRule 标准化奖励重复兑换规则，空值默认为不重复。
func normalizedRewardRepeatRule(rule string) string {
	switch strings.TrimSpace(rule) {
	case consts.RewardRepeatNone, consts.RewardRepeatDaily, consts.RewardRepeatWeekly, consts.RewardRepeatMonthly, consts.RewardRepeatCustom:
		return strings.TrimSpace(rule)
	default:
		return consts.RewardRepeatNone
	}
}

// rewardRecordRange 计算兑换历史筛选时间范围。
func rewardRecordRange(in v1.RewardRecordListInput) (string, string) {
	if strings.TrimSpace(in.Month) != "" {
		start, err := time.ParseInLocation("2006-01", strings.TrimSpace(in.Month), time.Local)
		if err == nil {
			end := start.AddDate(0, 1, -1)
			return start.Format(consts.DateLayout) + " 00:00:00", end.Format(consts.DateLayout) + " 23:59:59"
		}
	}
	from := ""
	to := ""
	if strings.TrimSpace(in.From) != "" {
		from = strings.TrimSpace(in.From) + " 00:00:00"
	}
	if strings.TrimSpace(in.To) != "" {
		to = strings.TrimSpace(in.To) + " 23:59:59"
	}
	return from, to
}
