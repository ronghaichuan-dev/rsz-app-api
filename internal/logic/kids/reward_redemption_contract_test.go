package kids

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1RewardRedemptionBundleUsesOneCanonicalCommitTime 验证兑换成功、幂等重放和同步变更都复用同一事务时间。
func TestV1RewardRedemptionBundleUsesOneCanonicalCommitTime(t *testing.T) {
	committedAt := time.UnixMilli(1_700_000_000_123)
	bundle := v1RewardRedemptionBundleFixture(committedAt)
	if err := v1ValidateRewardRedemptionCommitBundle(bundle, committedAt); err != nil {
		t.Fatalf("规范兑换 bundle 被拒绝: %v", err)
	}
	if err := v1.ValidateV1ResponseData("redeemReward", bundle); err != nil {
		t.Fatalf("兑换 bundle 不符合冻结合同: %v", err)
	}
	changes := v1SyncChanges(map[string]any{
		"exchange":            bundle["exchange"],
		"ledger_entry":        bundle["ledger_entry"],
		"balance":             bundle["balance"],
		"cooldown":            bundle["cooldown"],
		"notification_outbox": bundle["notification_outbox"],
	})
	if changes["exchanges"].([]any)[0].(map[string]any)["exchanged_at_ms"] != committedAt.UnixMilli() || changes["ledger_entries"].([]any)[0].(map[string]any)["created_at_ms"] != committedAt.UnixMilli() || changes["star_balances"].([]any)[0].(map[string]any)["updated_at_ms"] != committedAt.UnixMilli() {
		t.Fatalf("同步投影没有保留 canonical commit 时间: %v", changes)
	}
	replay := v1ReplayOutput(&v1.V1OperationOutput{Data: bundle, Status: 200, ChangeCursor: bundle["change_cursor"].(string)})
	if err := v1ValidateRewardRedemptionCommitBundle(replay.Data, committedAt); err != nil {
		t.Fatalf("兑换幂等重放没有保留规范 bundle: %v", err)
	}
	if replay.Data["receipt"].(map[string]any)["result_kind"] != "idempotent_replay" {
		t.Fatalf("兑换幂等重放未标记为 replay: %v", replay.Data)
	}
}

// TestV1ExchangeHistoryProjectionRetainsLedgerAuditTime 验证兑换历史从持久化审计记录读取同一 ledger 标识、金额对应时间和提交序列。
func TestV1ExchangeHistoryProjectionRetainsLedgerAuditTime(t *testing.T) {
	committedAt := time.UnixMilli(1_700_000_000_123)
	bundle := v1RewardRedemptionBundleFixture(committedAt)
	exchange := bundle["exchange"].(map[string]any)
	avatarJSON, err := json.Marshal(exchange["member_avatar_snapshot"])
	if err != nil {
		t.Fatalf("序列化成员头像快照失败: %v", err)
	}
	visualJSON, err := json.Marshal(exchange["reward_visual_snapshot"])
	if err != nil {
		t.Fatalf("序列化奖励视觉快照失败: %v", err)
	}
	row := gdb.Record{
		"exchange_id":                      gvar.New(exchange["exchange_id"]),
		"circle_id":                        gvar.New(exchange["circle_id"]),
		"member_id":                        gvar.New(exchange["member_id"]),
		"member_name_snapshot":             gvar.New(exchange["member_name_snapshot"]),
		"member_avatar_snapshot":           gvar.New(string(avatarJSON)),
		"reward_id":                        gvar.New(exchange["reward_id"]),
		"reward_title_snapshot":            gvar.New(exchange["reward_title_snapshot"]),
		"reward_visual_snapshot":           gvar.New(string(visualJSON)),
		"stars_deducted_snapshot":          gvar.New(exchange["stars_deducted_snapshot"]),
		"reward_repeat_rule_snapshot":      gvar.New(exchange["reward_repeat_rule_snapshot"]),
		"reward_cooldown_days_snapshot":    gvar.New(exchange["reward_cooldown_days_snapshot"]),
		"cooldown_until_at_snapshot":       gvar.New(time.UnixMilli(exchange["cooldown_until_ms_snapshot"].(int64))),
		"permanently_unavailable_snapshot": gvar.New(exchange["permanently_unavailable_snapshot"]),
		"ledger_id":                        gvar.New(exchange["ledger_id"]),
		"exchanged_at":                     gvar.New(committedAt),
		"commit_sequence":                  gvar.New(exchange["commit_sequence"]),
	}
	history := v1ExchangeRecordProjection(row)
	if !reflect.DeepEqual(history, exchange) {
		t.Fatalf("兑换历史改变了不可变审计投影: history=%v exchange=%v", history, exchange)
	}
	ledger := bundle["ledger_entry"].(map[string]any)
	if history["ledger_id"] != ledger["ledger_id"] || history["exchanged_at_ms"] != ledger["created_at_ms"] || ledger["delta"] != -history["stars_deducted_snapshot"].(int64) {
		t.Fatalf("兑换历史与流水没有闭合: history=%v ledger=%v", history, ledger)
	}
}

// v1RewardRedemptionBundleFixture 构造与数据库无关的完整临时奖励兑换 canonical bundle。
func v1RewardRedemptionBundleFixture(committedAt time.Time) map[string]any {
	circleID := "circle:v1:00000000-0000-4000-8000-000000000001"
	memberID := "member:v1:00000000-0000-4000-8000-000000000001"
	rewardID := "reward:v1:00000000-0000-4000-8000-000000000001"
	exchangeID := "exchange:v1:00000000-0000-4000-8000-000000000001"
	commitID := "commit:v1:00000000-0000-4000-8000-000000000001"
	ledgerID := "star-transaction:v1:00000000-0000-4000-8000-000000000001"
	sequence := int64(42)
	visual := map[string]any{"kind": "emoji", "emoji": "🎁"}
	actor := map[string]any{"actor_type": "owner", "actor_id": "admin:v1:00000000-0000-4000-8000-000000000001", "role": "owner", "display_name_snapshot": "管理员"}
	source := map[string]any{"source_type": "reward", "source_id": rewardID, "title_snapshot": "临时奖励", "stars_snapshot": int64(10), "asset_id_snapshot": nil, "scheduled_date_snapshot": nil}
	ledger := v1LedgerProjection(ledgerID, circleID, memberID, source, -10, nil, actor, nil, committedAt, sequence)
	balance := v1BalanceProjection(circleID, memberID, 90, 2, commitID, sequence, committedAt)
	receipt := map[string]any{"receipt_id": "receipt:v1:00000000-0000-4000-8000-000000000001", "commit_id": commitID, "commit_sequence": sequence, "result_kind": "first_committed", "committed_at_ms": committedAt.UnixMilli()}
	exchange := map[string]any{"exchange_id": exchangeID, "circle_id": circleID, "member_id": memberID, "member_name_snapshot": "成员", "member_avatar_snapshot": visual, "reward_id": rewardID, "reward_title_snapshot": "临时奖励", "reward_visual_snapshot": visual, "stars_deducted_snapshot": int64(10), "reward_repeat_rule_snapshot": "daily", "reward_cooldown_days_snapshot": nil, "cooldown_until_ms_snapshot": committedAt.AddDate(0, 0, 1).UnixMilli(), "permanently_unavailable_snapshot": false, "ledger_id": ledgerID, "exchanged_at_ms": committedAt.UnixMilli(), "commit_sequence": sequence}
	cooldown := v1CooldownProjection(rewardID, memberID, committedAt.AddDate(0, 0, 1), false, committedAt, 1)
	notification := v1NotificationProjection("notification:v1:00000000-0000-4000-8000-000000000001", exchangeID, committedAt)
	return map[string]any{"receipt": receipt, "exchange": exchange, "ledger_entry": ledger, "balance": balance, "cooldown": cooldown, "notification_outbox": notification, "change_cursor": v1CommitCursor(sequence)}
}
