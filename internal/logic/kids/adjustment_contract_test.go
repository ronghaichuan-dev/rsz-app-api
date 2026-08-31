package kids

import (
	"strings"
	"testing"
	"time"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1AdjustmentBundleUsesCanonicalLedgerID 验证正常人工调整的完整 bundle 可通过冻结合同，并保持 receipt、流水、余额和 cursor 的跨字段关系。
func TestV1AdjustmentBundleUsesCanonicalLedgerID(t *testing.T) {
	bundle := v1AdjustmentBundleFixture()
	if err := v1.ValidateV1ResponseData("adjustMemberStars", bundle); err != nil {
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
