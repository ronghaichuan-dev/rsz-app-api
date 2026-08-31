package kids

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1AdjustmentBundleUsesCanonicalLedgerID 验证正常人工调整的完整 bundle 可通过冻结合同，并保持 receipt、流水、余额和 cursor 的跨字段关系。
func TestV1AdjustmentBundleUsesCanonicalLedgerID(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	if err := v1ValidateAdjustmentCommitBundle(bundle, "circle:v1:00000000-0000-4000-8000-000000000001", "member:v1:00000000-0000-4000-8000-000000000001", "adjustment:v1:00000000-0000-4000-8000-000000000001", 1, 5, "测试调整", time.UnixMilli(1_700_000_000_000)); err != nil {
		t.Fatalf("规范 adjustment bundle 被拒绝: %v", err)
	}
	ledger := bundle["ledger_entry"].(map[string]any)
	if !strings.HasPrefix(ledger["ledger_id"].(string), "star-transaction:v1:") {
		t.Fatalf("流水标识未使用冻结命名空间: %s", ledger["ledger_id"])
	}
	receipt := bundle["receipt"].(map[string]any)
	balance := bundle["balance"].(map[string]any)
	if ledger["commit_sequence"] != receipt["commit_sequence"] || balance["source_commit_id"] != receipt["commit_id"] || balance["source_commit_sequence"] != receipt["commit_sequence"] || bundle["change_cursor"] != v1CommitCursor(receipt["commit_sequence"].(int64)) {
		t.Fatal("adjustment bundle 的 receipt、流水、余额和 cursor 不一致")
	}
}

// TestV1AdjustmentBundleSupportsLargestLegalIncrease 验证合同允许的最大人工加星仍可组成规范的原子提交响应。
func TestV1AdjustmentBundleSupportsLargestLegalIncrease(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	bundle["ledger_entry"].(map[string]any)["delta"] = int64(1_000_000)
	bundle["balance"].(map[string]any)["balance"] = int64(1_000_000)
	if err := v1ValidateAdjustmentCommitBundle(bundle, "circle:v1:00000000-0000-4000-8000-000000000001", "member:v1:00000000-0000-4000-8000-000000000001", "adjustment:v1:00000000-0000-4000-8000-000000000001", 1, 1_000_000, "测试调整", time.UnixMilli(1_700_000_000_000)); err != nil {
		t.Fatalf("最大合法人工加星不应被响应合同拒绝: %v", err)
	}
}

// TestV1AdjustmentBundleRejectsInconsistentCommittedTime 防止 receipt、ledger 与 balance 使用不同时间却返回 200。
func TestV1AdjustmentBundleRejectsInconsistentCommittedTime(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	bundle["balance"].(map[string]any)["updated_at_ms"] = int64(1_700_000_000_001)
	if err := v1ValidateAdjustmentCommitBundle(bundle, "circle:v1:00000000-0000-4000-8000-000000000001", "member:v1:00000000-0000-4000-8000-000000000001", "adjustment:v1:00000000-0000-4000-8000-000000000001", 1, 5, "测试调整", time.UnixMilli(1_700_000_000_000)); err == nil {
		t.Fatal("不一致的调整提交时间不应通过原子 bundle 校验")
	}
}

// TestV1AdjustmentEnvelopeAndReplay 验证成功 envelope 与幂等重放保留同一原子 bundle 和 cursor。
func TestV1AdjustmentEnvelopeAndReplay(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	cursor := bundle["change_cursor"].(string)
	response := &v1.V1Response{V1Version: v1.V1Version, RequestID: "request:adjustment-envelope", ServerTimeMs: time.Now().UnixMilli(), Data: bundle, ChangeCursor: cursor}
	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("序列化 adjustment envelope 失败: %v", err)
	}
	var envelope map[string]any
	if err = json.Unmarshal(encoded, &envelope); err != nil {
		t.Fatalf("解析 adjustment envelope 失败: %v", err)
	}
	data := envelope["data"].(map[string]any)
	if len(data) != 4 || envelope["change_cursor"] != data["change_cursor"] {
		t.Fatalf("adjustment envelope 的 data 字段或 cursor 不符合冻结合同: %v", envelope)
	}

	replay := v1ReplayOutput(&v1.V1OperationOutput{Data: bundle, Status: 200, ChangeCursor: cursor})
	if err := v1ValidateAdjustmentOutputCursor("adjustMemberStars", replay); err != nil {
		t.Fatalf("幂等重放的 envelope cursor 不应被拒绝: %v", err)
	}
	replayReceipt := replay.Data["receipt"].(map[string]any)
	if replay.ChangeCursor != cursor || replay.Data["change_cursor"] != cursor || replayReceipt["result_kind"] != "idempotent_replay" {
		t.Fatalf("幂等重放没有保留原子 cursor 或 replay receipt: %v", replay)
	}
	firstLedger, _ := json.Marshal(bundle["ledger_entry"])
	replayLedger, _ := json.Marshal(replay.Data["ledger_entry"])
	firstBalance, _ := json.Marshal(bundle["balance"])
	replayBalance, _ := json.Marshal(replay.Data["balance"])
	if string(firstLedger) != string(replayLedger) || string(firstBalance) != string(replayBalance) {
		t.Fatal("幂等重放不应生成第二条流水或再次改变余额")
	}
}

// TestV1AdjustmentOutputRejectsMismatchedEnvelopeCursor 防止外层 success envelope 与 data 使用不同 cursor。
func TestV1AdjustmentOutputRejectsMismatchedEnvelopeCursor(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	if err := v1ValidateAdjustmentOutputCursor("adjustMemberStars", &v1.V1OperationOutput{Data: bundle, Status: 200, ChangeCursor: "cur:v1:00000043"}); err == nil {
		t.Fatal("不一致的 adjustment envelope cursor 不应通过")
	}
}

// TestV1AdjustmentBundleRejectsLegacyLedgerNamespace 防止 ledger:v1 命名空间再次导致成功写入在提交后才被拒绝。
func TestV1AdjustmentBundleRejectsLegacyLedgerNamespace(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	bundle["ledger_entry"].(map[string]any)["ledger_id"] = "ledger:v1:00000000-0000-4000-8000-000000000001"
	if err := v1.ValidateV1ResponseData("adjustMemberStars", bundle); err == nil {
		t.Fatal("旧 ledger 命名空间不应通过 adjustment 合同")
	}
}

// TestV1AdjustmentSyncChangesIncludesLedgerAndBalance 验证 pull 场景可获得与 mutation bundle 相同的流水和余额投影。
func TestV1AdjustmentSyncChangesIncludesLedgerAndBalance(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	changes := v1SyncChanges(map[string]any{"ledger_entry": bundle["ledger_entry"], "balance": bundle["balance"]})
	if len(changes["ledger_entries"].([]any)) != 1 || len(changes["star_balances"].([]any)) != 1 {
		t.Fatal("adjustment commit 未包含完整的 ledger 和 balance 同步投影")
	}
}

// v1AdjustmentBundleFixture 构造无需数据库的合法 adjustment canonical bundle，用于固定 wire contract。
func v1AdjustmentBundleFixture() map[string]any {
	now := time.UnixMilli(1_700_000_000_000)
	circleID := "circle:v1:00000000-0000-4000-8000-000000000001"
	memberID := "member:v1:00000000-0000-4000-8000-000000000001"
	adjustmentID := "adjustment:v1:00000000-0000-4000-8000-000000000001"
	commitID := "commit:v1:00000000-0000-4000-8000-000000000001"
	sequence := int64(42)
	actor := map[string]any{"actor_type": "owner", "actor_id": "admin:v1:00000000-0000-4000-8000-000000000001", "role": "owner", "display_name_snapshot": "管理员"}
	source := map[string]any{"source_type": "adjustment", "source_id": adjustmentID, "title_snapshot": nil, "stars_snapshot": nil, "asset_id_snapshot": nil, "scheduled_date_snapshot": nil}
	ledger := v1LedgerProjection(v1ID("ledger", "adjustment-contract"), circleID, memberID, source, 5, "测试调整", actor, nil, now, sequence)
	balance := v1BalanceProjection(circleID, memberID, 5, 2, commitID, sequence, now)
	receipt := map[string]any{"receipt_id": "receipt:v1:00000000-0000-4000-8000-000000000001", "commit_id": commitID, "commit_sequence": sequence, "result_kind": "first_committed", "committed_at_ms": now.UnixMilli()}
	return map[string]any{"receipt": receipt, "ledger_entry": ledger, "balance": balance, "change_cursor": v1CommitCursor(sequence)}
}
