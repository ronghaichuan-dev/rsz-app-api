package kids

import (
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1CompleteTaskCanonicalBundle 验证首次提交与幂等重放都保留同一个任务完成规范时间。
func TestV1CompleteTaskCanonicalBundle(t *testing.T) {
	const committedAt = int64(1788229645577)
	first := &v1.V1OperationOutput{
		Status:       200,
		ChangeCursor: "cur:v1:2a",
		Data: map[string]any{
			"receipt":      map[string]any{"result_kind": "first_committed", "committed_at_ms": committedAt},
			"occurrence":   map[string]any{"updated_at_ms": committedAt},
			"completion":   map[string]any{"completed_at_ms": committedAt},
			"ledger_entry": map[string]any{"created_at_ms": committedAt},
			"balance":      map[string]any{"updated_at_ms": committedAt},
		},
	}
	if err := v1ValidateCompleteTaskCanonicalBundle(first.Data); err != nil {
		t.Fatalf("首次提交 canonical bundle 校验失败: %v", err)
	}
	replay := v1ReplayOutput(first)
	if err := v1ValidateCompleteTaskCanonicalBundle(replay.Data); err != nil {
		t.Fatalf("幂等重放 canonical bundle 校验失败: %v", err)
	}
	if replay.Data["receipt"].(map[string]any)["result_kind"] != "idempotent_replay" {
		t.Fatal("幂等重放未标记为 idempotent_replay")
	}
}

// TestV1CompleteTaskCanonicalBundleRejectsMismatch 验证任一投影时间不一致都会被拒绝，避免写入错误幂等快照。
func TestV1CompleteTaskCanonicalBundleRejectsMismatch(t *testing.T) {
	bundle := map[string]any{
		"receipt":      map[string]any{"committed_at_ms": int64(1788229645577)},
		"occurrence":   map[string]any{"updated_at_ms": int64(1788229645577)},
		"completion":   map[string]any{"completed_at_ms": int64(1788229645577)},
		"ledger_entry": map[string]any{"created_at_ms": int64(1788229645577)},
		"balance":      map[string]any{"updated_at_ms": int64(1788229646000)},
	}
	if err := v1ValidateCompleteTaskCanonicalBundle(bundle); err == nil {
		t.Fatal("不一致的 canonical bundle 未被拒绝")
	}
}
