package kids

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// v1UpsertReward 创建或更新奖励定义及其成员分配，整个写入使用同一提交序列。
func (s *sKids) v1UpsertReward(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, rewardID := in.PathParameters["circle_id"], in.PathParameters["reward_id"]
	assignedMemberIDs, err := v1StringSlice(in.Body["assigned_member_ids"])
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "assigned members are invalid")
	}
	now := time.Now()
	var reward, receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageRewards); err != nil {
			return err
		}
		row, err := tx.Model(consts.KidsV1RewardTable).Ctx(ctx).Where("reward_id", rewardID).Where("circle_id", circleID).LockUpdate().One()
		if err != nil {
			return err
		}
		expected, hasExpected := v1ExpectedVersion(in.Body["expected_version"])
		if row.IsEmpty() {
			if in.Body["expected_version"] != nil {
				return v1Error(409, "VERSION_CONFLICT", false, "reward does not exist")
			}
			reward = v1RewardProjection(rewardID, circleID, in.Body, 1, "active", now, now, time.Time{})
			if _, err = tx.Model(consts.KidsV1RewardTable).Ctx(ctx).Data(v1RewardValues(reward, now)).Insert(); err != nil {
				return err
			}
		} else {
			version := row["version"].Int64()
			if !hasExpected || expected != version {
				return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &version, Message: "reward version conflicts"}
			}
			if row["status"].String() != "active" {
				return v1Error(409, "VERSION_CONFLICT", false, "deleted reward cannot be updated")
			}
			reward = v1RewardProjection(rewardID, circleID, in.Body, version+1, "active", row["created_at"].Time(), now, time.Time{})
			if _, err = tx.Model(consts.KidsV1RewardTable).Ctx(ctx).Where("id", row["id"].Int64()).Data(v1RewardValues(reward, now)).Update(); err != nil {
				return err
			}
		}
		if _, err = tx.Model(consts.KidsV1RewardAssignmentTable).Ctx(ctx).Where("reward_id", rewardID).Delete(); err != nil {
			return err
		}
		for _, memberID := range assignedMemberIDs {
			member, queryErr := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").One()
			if queryErr != nil {
				return queryErr
			}
			if member.IsEmpty() {
				return v1Error(422, "VALIDATION_FAILED", false, "assigned member is missing")
			}
			if _, err = tx.Model(consts.KidsV1RewardAssignmentTable).Ctx(ctx).Data(gdb.Map{"reward_id": rewardID, "member_id": memberID, "created_at": now}).Insert(); err != nil {
				return err
			}
		}
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"reward_id": rewardID})
		return err
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "reward": reward}, cursor, nil
}

// v1DeleteReward 软删除奖励定义，并保留兑换和账本审计事实。
func (s *sKids) v1DeleteReward(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, rewardID := in.PathParameters["circle_id"], in.PathParameters["reward_id"]
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "reward version is invalid")
	}
	now := time.Now()
	var tombstone, receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageRewards)
		if err != nil {
			return err
		}
		row, err := tx.Model(consts.KidsV1RewardTable).Ctx(ctx).Where("reward_id", rewardID).Where("circle_id", circleID).LockUpdate().One()
		if err != nil {
			return err
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "reward is missing")
		}
		version := row["version"].Int64()
		if version != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &version, Message: "reward version conflicts"}
		}
		actor, err := v1ActorSnapshotTx(ctx, tx, membership)
		if err != nil {
			return err
		}
		version++
		if _, err = tx.Model(consts.KidsV1RewardTable).Ctx(ctx).Where("id", row["id"].Int64()).Data(gdb.Map{"status": "deleted", "version": version, "deleted_at": now, "updated_at": now}).Update(); err != nil {
			return err
		}
		tombstone = v1EntityTombstone("reward", rewardID, version, now, actor)
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"reward_tombstone": tombstone})
		return err
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "tombstone": tombstone}, cursor, nil
}

// v1RewardEligibility 返回奖励分配、冷却和余额共同决定的可兑换状态。
func (s *sKids) v1RewardEligibility(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, rewardID, memberID := in.PathParameters["circle_id"], in.PathParameters["reward_id"], v1QueryFirst(in.Query, "member_id")
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	reward, err := utils.KidsDB(ctx).Model(consts.KidsV1RewardTable).Ctx(ctx).Where("reward_id", rewardID).Where("circle_id", circleID).One()
	if err != nil {
		return nil, "", err
	}
	if reward.IsEmpty() {
		return nil, "", v1Error(404, "NOT_FOUND", false, "reward is missing")
	}
	balance, err := utils.KidsDB(ctx).Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).One()
	if err != nil {
		return nil, "", err
	}
	assignment, err := utils.KidsDB(ctx).Model(consts.KidsV1RewardAssignmentTable).Ctx(ctx).Where("reward_id", rewardID).Where("member_id", memberID).One()
	if err != nil {
		return nil, "", err
	}
	cooldown, err := utils.KidsDB(ctx).Model(consts.KidsV1RewardCooldownTable).Ctx(ctx).Where("reward_id", rewardID).Where("member_id", memberID).One()
	if err != nil {
		return nil, "", err
	}
	cursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	return v1RewardEligibilityProjection(reward, memberID, assignment, cooldown, balance, cursor), "", nil
}

// v1RedeemReward 在同一事务中扣减余额、写入流水、兑换审计、冷却和通知 outbox。
func (s *sKids) v1RedeemReward(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, rewardID, memberID := in.PathParameters["circle_id"], fmt.Sprint(in.Body["reward_id"]), fmt.Sprint(in.Body["member_id"])
	expected, ok := v1ExpectedVersion(in.Body["expected_reward_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "reward version is invalid")
	}
	now := time.Now()
	var exchange, ledger, balance, cooldownOut, notification, receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, "")
		if err != nil {
			return err
		}
		if membership["actor_type"].String() == "member" && membership["actor_id"].String() != memberID {
			return v1Error(403, "FORBIDDEN", false, "member can only redeem own reward")
		}
		reward, err := tx.Model(consts.KidsV1RewardTable).Ctx(ctx).Where("reward_id", rewardID).Where("circle_id", circleID).LockUpdate().One()
		if err != nil {
			return err
		}
		if reward.IsEmpty() || reward["status"].String() != "active" {
			return v1Error(404, "NOT_FOUND", false, "reward is unavailable")
		}
		if version := reward["version"].Int64(); version != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &version, Message: "reward version conflicts"}
		}
		assigned, err := tx.Model(consts.KidsV1RewardAssignmentTable).Ctx(ctx).Where("reward_id", rewardID).Where("member_id", memberID).One()
		if err != nil {
			return err
		}
		if assigned.IsEmpty() {
			return v1Error(422, "NOT_ASSIGNED", false, "reward is not assigned to member")
		}
		member, err := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").One()
		if err != nil {
			return err
		}
		if member.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "member is missing")
		}
		cooldown, err := tx.Model(consts.KidsV1RewardCooldownTable).Ctx(ctx).Where("reward_id", rewardID).Where("member_id", memberID).LockUpdate().One()
		if err != nil {
			return err
		}
		if !cooldown.IsEmpty() && cooldown["permanently_unavailable"].Bool() {
			return v1Error(422, "PERMANENTLY_UNAVAILABLE", false, "reward is permanently unavailable")
		}
		if !cooldown.IsEmpty() && !cooldown["cooldown_until_at"].Time().IsZero() && now.Before(cooldown["cooldown_until_at"].Time()) {
			return v1Error(422, "COOLDOWN_ACTIVE", false, "reward is cooling down")
		}
		balanceRow, err := tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).LockUpdate().One()
		if err != nil {
			return err
		}
		if balanceRow.IsEmpty() {
			return v1Error(409, "AUDIT_INCONSISTENT", false, "canonical member balance is missing")
		}
		if balanceRow["balance"].Int64() < reward["stars_required"].Int64() {
			return v1Error(422, "INSUFFICIENT_BALANCE", false, "member balance is insufficient")
		}
		actor, err := v1ActorSnapshotTx(ctx, tx, membership)
		if err != nil {
			return err
		}
		permanent, until := v1RewardCooldown(ruleString(reward["repeat_rule"].String()), reward["cooldown_days"].Int64(), now)
		ledgerID, exchangeID := v1ID("ledger", uuid.NewString()), fmt.Sprint(in.Body["exchange_id"])
		source := map[string]any{"source_type": "reward", "source_id": rewardID, "title_snapshot": reward["title"].String(), "stars_snapshot": reward["stars_required"].Int64(), "asset_id_snapshot": nil, "scheduled_date_snapshot": nil}
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Data(gdb.Map{"ledger_id": ledgerID, "circle_id": circleID, "member_id": memberID, "source": mustV1JSON(source), "delta": -reward["stars_required"].Int64(), "reason": nil, "actor": mustV1JSON(actor), "commit_sequence": 0, "created_at": now}).Insert(); err != nil {
			return err
		}
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"exchange_id": exchangeID, "ledger_id": ledgerID})
		if err != nil {
			return err
		}
		sequence, commitID := receipt["commit_sequence"].(int64), receipt["commit_id"].(string)
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("ledger_id", ledgerID).Data(gdb.Map{"commit_sequence": sequence}).Update(); err != nil {
			return err
		}
		cooldownVersion := int64(1)
		if !cooldown.IsEmpty() {
			cooldownVersion = cooldown["version"].Int64() + 1
			if _, err = tx.Model(consts.KidsV1RewardCooldownTable).Ctx(ctx).Where("id", cooldown["id"].Int64()).Data(gdb.Map{"cooldown_until_at": nullableTime(until), "last_redeemed_at": now, "permanently_unavailable": permanent, "version": cooldownVersion, "updated_at": now}).Update(); err != nil {
				return err
			}
		} else if _, err = tx.Model(consts.KidsV1RewardCooldownTable).Ctx(ctx).Data(gdb.Map{"reward_id": rewardID, "member_id": memberID, "cooldown_until_at": nullableTime(until), "last_redeemed_at": now, "permanently_unavailable": permanent, "version": cooldownVersion, "created_at": now, "updated_at": now}).Insert(); err != nil {
			return err
		}
		nextBalance, nextVersion := balanceRow["balance"].Int64()-reward["stars_required"].Int64(), balanceRow["version"].Int64()+1
		if _, err = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("id", balanceRow["id"].Int64()).Data(gdb.Map{"balance": nextBalance, "version": nextVersion, "source_commit_id": commitID, "source_commit_sequence": sequence, "updated_at": now}).Update(); err != nil {
			return err
		}
		avatar := v1JSONValue(member["avatar"].String())
		if _, err = tx.Model(consts.KidsV1ExchangeTable).Ctx(ctx).Data(gdb.Map{"exchange_id": exchangeID, "circle_id": circleID, "member_id": memberID, "member_name_snapshot": member["display_name"], "member_avatar_snapshot": mustV1JSON(avatar), "reward_id": rewardID, "reward_title_snapshot": reward["title"], "reward_visual_snapshot": reward["visual"], "stars_deducted_snapshot": reward["stars_required"], "reward_repeat_rule_snapshot": reward["repeat_rule"], "reward_cooldown_days_snapshot": nullableInt(reward["cooldown_days"].Int64()), "cooldown_until_at_snapshot": nullableTime(until), "permanently_unavailable_snapshot": permanent, "ledger_id": ledgerID, "commit_sequence": sequence, "exchanged_at": now}).Insert(); err != nil {
			return err
		}
		notificationID := v1ID("notification", uuid.NewString())
		if _, err = tx.Model(consts.KidsV1NotificationOutboxTable).Ctx(ctx).Data(gdb.Map{"notification_id": notificationID, "circle_id": circleID, "account_id": in.PrincipalID, "exchange_id": exchangeID, "event_type": "reward_redeemed", "payload": mustV1JSON(map[string]any{"exchange_id": exchangeID}), "commit_sequence": sequence, "status": "pending", "attempt_count": 0, "next_attempt_at": nil, "created_at": now, "updated_at": now, "version": 1}).Insert(); err != nil {
			return err
		}
		ledger = v1LedgerProjection(ledgerID, circleID, memberID, source, -reward["stars_required"].Int64(), nil, actor, nil, now, sequence)
		balance = v1BalanceProjection(circleID, memberID, nextBalance, nextVersion, commitID, sequence, now)
		cooldownOut = v1CooldownProjection(rewardID, memberID, until, permanent, now, cooldownVersion)
		exchange = v1ExchangeProjection(exchangeID, circleID, memberID, member["display_name"].String(), avatar, rewardID, reward, ledgerID, until, permanent, now, sequence)
		notification = v1NotificationProjection(notificationID, exchangeID, now)
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "exchange": exchange, "ledger_entry": ledger, "balance": balance, "cooldown": cooldownOut, "notification_outbox": notification, "change_cursor": cursor}, cursor, nil
}

// v1ExchangeHistory 返回指定成员按兑换时间倒序排列的稳定分页历史。
func (s *sKids) v1ExchangeHistory(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, memberID := in.PathParameters["circle_id"], v1QueryFirst(in.Query, "member_id")
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	limit, err := strconv.Atoi(v1QueryFirst(in.Query, "limit"))
	if err != nil || limit < 1 || limit > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "exchange limit is invalid")
	}
	offset, err := v1PageOffset(v1QueryFirst(in.Query, "cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "exchange cursor is invalid")
	}
	model := utils.KidsDB(ctx).Model(consts.KidsV1ExchangeTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID)
	if v1QueryFirst(in.Query, "filter_mode") == "calendar_range" {
		model = model.Where("exchanged_at >= ?", v1QueryFirst(in.Query, "start_date")).Where("exchanged_at < ?", v1QueryFirst(in.Query, "end_date_exclusive"))
	}
	total, err := model.Count()
	if err != nil {
		return nil, "", err
	}
	rows, err := model.OrderDesc("exchanged_at,id").Limit(offset, limit+1).All()
	if err != nil {
		return nil, "", err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, v1ExchangeRecordProjection(row))
	}
	cursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	next := any(nil)
	if hasMore {
		next = v1PageCursor(offset + limit)
	}
	return map[string]any{"items": items, "next_cursor": next, "has_more": hasMore, "snapshot_cursor": cursor, "total_exchanges": total}, "", nil
}

// v1RewardProjection 构造奖励定义的接口冻结投影。
func v1RewardProjection(rewardID, circleID string, input map[string]any, version int64, status string, createdAt, updatedAt, deletedAt time.Time) map[string]any {
	assignedMemberIDs, _ := v1StringSlice(input["assigned_member_ids"])
	return map[string]any{"reward_id": rewardID, "circle_id": circleID, "title": fmt.Sprint(input["title"]), "description": fmt.Sprint(input["description"]), "visual": input["visual"], "stars_required": mustV1Integer(input["stars_required"]), "repeat_rule": fmt.Sprint(input["repeat_rule"]), "cooldown_days": input["cooldown_days"], "zone_id": fmt.Sprint(input["zone_id"]), "assigned_member_ids": assignedMemberIDs, "status": status, "version": version, "created_at_ms": createdAt.UnixMilli(), "updated_at_ms": updatedAt.UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(deletedAt)}
}

// v1RewardValues 将奖励投影转换为数据库写入列。
func v1RewardValues(reward map[string]any, now time.Time) gdb.Map {
	return gdb.Map{"reward_id": reward["reward_id"], "circle_id": reward["circle_id"], "title": reward["title"], "description": reward["description"], "visual": mustV1JSON(reward["visual"]), "stars_required": reward["stars_required"], "repeat_rule": reward["repeat_rule"], "cooldown_days": reward["cooldown_days"], "zone_id": reward["zone_id"], "status": reward["status"], "version": reward["version"], "created_at": time.UnixMilli(reward["created_at_ms"].(int64)), "updated_at": now}
}

// v1RewardEligibilityProjection 构造不改变事实状态的奖励资格快照。
func v1RewardEligibilityProjection(reward gdb.Record, memberID string, assignment, cooldown, balance gdb.Record, cursor string) map[string]any {
	status, until, last := "eligible", any(nil), any(nil)
	if reward["status"].String() != "active" {
		status = "unavailable"
	}
	if assignment.IsEmpty() {
		status = "not_assigned"
	}
	if !cooldown.IsEmpty() {
		until, last = v1NullableTimeMillis(cooldown["cooldown_until_at"].Time()), cooldown["last_redeemed_at"].Time().UnixMilli()
		if cooldown["permanently_unavailable"].Bool() {
			status = "permanently_unavailable"
		} else if !cooldown["cooldown_until_at"].Time().IsZero() && time.Now().Before(cooldown["cooldown_until_at"].Time()) {
			status = "cooling_down"
		}
	}
	version := int64(1)
	if !balance.IsEmpty() {
		version = balance["version"].Int64()
	}
	return map[string]any{"reward_id": reward["reward_id"].String(), "member_id": memberID, "status": status, "cooldown_until_ms": until, "last_redeemed_at_ms": last, "reward_version": reward["version"].Int64(), "balance_version": version, "as_of_cursor": cursor}
}

// v1RewardCooldown 根据重复规则计算下一次允许兑换时间。
func v1RewardCooldown(repeatRule string, days int64, now time.Time) (bool, time.Time) {
	switch repeatRule {
	case consts.RewardRepeatNone:
		return true, time.Time{}
	case consts.RewardRepeatDaily:
		return false, now.AddDate(0, 0, 1)
	case consts.RewardRepeatWeekly:
		return false, now.AddDate(0, 0, 7)
	case consts.RewardRepeatMonthly:
		return false, now.AddDate(0, 1, 0)
	default:
		return false, now.AddDate(0, 0, int(days))
	}
}

// v1CooldownProjection 构造奖励冷却接口投影。
func v1CooldownProjection(rewardID, memberID string, until time.Time, permanent bool, last time.Time, version int64) map[string]any {
	return map[string]any{"reward_id": rewardID, "member_id": memberID, "cooldown_until_ms": v1NullableTimeMillis(until), "permanently_unavailable": permanent, "last_redeemed_at_ms": last.UnixMilli(), "version": version}
}

// v1ExchangeProjection 构造 append-only 兑换快照。
func v1ExchangeProjection(exchangeID, circleID, memberID, memberName string, avatar any, rewardID string, reward gdb.Record, ledgerID string, until time.Time, permanent bool, now time.Time, sequence int64) map[string]any {
	return map[string]any{"exchange_id": exchangeID, "circle_id": circleID, "member_id": memberID, "member_name_snapshot": memberName, "member_avatar_snapshot": avatar, "reward_id": rewardID, "reward_title_snapshot": reward["title"].String(), "reward_visual_snapshot": v1JSONValue(reward["visual"].String()), "stars_deducted_snapshot": reward["stars_required"].Int64(), "reward_repeat_rule_snapshot": reward["repeat_rule"].String(), "reward_cooldown_days_snapshot": nullableInt(reward["cooldown_days"].Int64()), "cooldown_until_ms_snapshot": v1NullableTimeMillis(until), "permanently_unavailable_snapshot": permanent, "ledger_id": ledgerID, "exchanged_at_ms": now.UnixMilli(), "commit_sequence": sequence}
}

// v1ExchangeRecordProjection 从持久化兑换记录恢复接口投影。
func v1ExchangeRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"exchange_id": row["exchange_id"].String(), "circle_id": row["circle_id"].String(), "member_id": row["member_id"].String(), "member_name_snapshot": row["member_name_snapshot"].String(), "member_avatar_snapshot": v1JSONValue(row["member_avatar_snapshot"].String()), "reward_id": row["reward_id"].String(), "reward_title_snapshot": row["reward_title_snapshot"].String(), "reward_visual_snapshot": v1JSONValue(row["reward_visual_snapshot"].String()), "stars_deducted_snapshot": row["stars_deducted_snapshot"].Int64(), "reward_repeat_rule_snapshot": row["reward_repeat_rule_snapshot"].String(), "reward_cooldown_days_snapshot": nullableInt(row["reward_cooldown_days_snapshot"].Int64()), "cooldown_until_ms_snapshot": v1NullableTimeMillis(row["cooldown_until_at_snapshot"].Time()), "permanently_unavailable_snapshot": row["permanently_unavailable_snapshot"].Bool(), "ledger_id": row["ledger_id"].String(), "exchanged_at_ms": row["exchanged_at"].Time().UnixMilli(), "commit_sequence": row["commit_sequence"].Int64()}
}

// v1NotificationProjection 构造独立于兑换提交重试状态的通知 outbox 投影。
func v1NotificationProjection(notificationID, exchangeID string, now time.Time) map[string]any {
	return map[string]any{"notification_id": notificationID, "event_type": "reward_redeemed", "exchange_id": exchangeID, "status": "pending", "attempt_count": int64(0), "next_attempt_at_ms": nil, "created_at_ms": now.UnixMilli(), "updated_at_ms": now.UnixMilli(), "version": int64(1)}
}

// mustV1Integer 返回经过 schema 校验的整数；异常值以零失败关闭。
func mustV1Integer(value any) int64 {
	integer, ok := v1Integer(value)
	if !ok {
		return 0
	}
	return integer
}

// nullableInt 将数据库零值转换为接口 nullable 整数。
func nullableInt(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

// nullableTime 将零时间转换为可写入数据库的 NULL。
func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

// ruleString 返回数据库重复规则字段。
func ruleString(value string) string { return value }

// v1PullCircleBootstrapDelta 按不可拆提交读取圈子同步页，不会拆分单个提交。
func (s *sKids) v1PullCircleBootstrapDelta(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	limit, err := strconv.Atoi(v1QueryFirst(in.Query, "limit"))
	if err != nil || limit < 1 || limit > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "sync limit is invalid")
	}
	if limit > 100 {
		limit = 100
	}
	sequence, err := v1SyncSequence(v1QueryFirst(in.Query, "change_cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "sync cursor is invalid")
	}
	rows, err := utils.KidsDB(ctx).Model(consts.KidsV1CommitTable).Ctx(ctx).Where("circle_id", circleID).Where("commit_sequence > ?", sequence).Order("commit_sequence ASC").Limit(0, limit+1).All()
	if err != nil {
		return nil, "", err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	commits := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var changes map[string]any
		if err = json.Unmarshal(row["change_payload"].Bytes(), &changes); err != nil {
			return nil, "", v1Error(409, "AUDIT_INCONSISTENT", false, "sync change payload is invalid")
		}
		commitSequence := row["commit_sequence"].Int64()
		commits = append(commits, map[string]any{"commit_id": row["commit_id"].String(), "commit_sequence": commitSequence, "committed_at_ms": row["created_at"].Time().UnixMilli(), "changes": v1SyncChanges(changes), "change_cursor": v1CommitCursor(commitSequence)})
	}
	snapshot, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	next := any(nil)
	if hasMore {
		next = v1CommitCursor(rows[len(rows)-1]["commit_sequence"].Int64())
	}
	return map[string]any{"commits": commits, "next_cursor": next, "has_more": hasMore, "snapshot_cursor": snapshot}, "", nil
}

// v1SyncSequence 解析由本服务签发的同步提交游标。
func v1SyncSequence(cursor string) (int64, error) {
	if cursor == "" {
		return 0, nil
	}
	const prefix = "cur:v1:"
	if !strings.HasPrefix(cursor, prefix) {
		return 0, fmt.Errorf("unsupported sync cursor")
	}
	sequence, err := strconv.ParseInt(strings.TrimPrefix(cursor, prefix), 16, 64)
	if err != nil || sequence < 0 {
		return 0, fmt.Errorf("invalid sync cursor")
	}
	return sequence, nil
}

// v1SyncChanges 将业务写入变化统一为 SyncChanges 的完整数组字段集合。
func v1SyncChanges(changes map[string]any) map[string]any {
	result := map[string]any{"circles": []any{}, "memberships": []any{}, "administrators": []any{}, "members": []any{}, "circle_selections": []any{}, "outgoing_invites": []any{}, "task_tags": []any{}, "tasks": []any{}, "task_occurrences": []any{}, "task_occurrence_tombstones": []any{}, "task_completions": []any{}, "task_cancellations": []any{}, "rewards": []any{}, "reward_cooldowns": []any{}, "ledger_entries": []any{}, "star_balances": []any{}, "exchanges": []any{}, "notification_outbox_entries": []any{}, "tombstones": []any{}}
	keyMap := map[string]string{"circle": "circles", "membership": "memberships", "administrator": "administrators", "member": "members", "selection": "circle_selections", "invite": "outgoing_invites", "tag": "task_tags", "task": "tasks", "occurrence": "task_occurrences", "completion": "task_completions", "cancellation": "task_cancellations", "reward": "rewards", "cooldown": "reward_cooldowns", "ledger": "ledger_entries", "ledger_entry": "ledger_entries", "reversal_ledger_entry": "ledger_entries", "balance": "star_balances", "exchange": "exchanges", "notification_outbox": "notification_outbox_entries"}
	for key, value := range changes {
		target, ok := keyMap[key]
		if !ok && strings.HasSuffix(key, "_tombstone") {
			target, ok = "tombstones", true
		}
		if !ok || value == nil {
			continue
		}
		if values, isArray := value.([]any); isArray {
			result[target] = append(result[target].([]any), values...)
		} else {
			result[target] = append(result[target].([]any), value)
		}
	}
	return result
}

// v1VerifyPlayPurchase 只允许服务端与 Google Play 的真实验证器完成 purchase 验证。
func (s *sKids) v1VerifyPlayPurchase(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	return nil, "", v1Error(503, "UNAVAILABLE", true, "Google Play verifier is not configured")
}
