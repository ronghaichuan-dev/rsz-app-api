package kids

import (
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1ReplayOutputMarksInviteRedemption 验证邀请码兑换的成功快照在重放时返回 idempotent_replay。
func TestV1ReplayOutputMarksInviteRedemption(t *testing.T) {
	first := &v1.V1OperationOutput{
		Status:       200,
		ChangeCursor: "cur:v1:2a",
		Data: map[string]any{
			"redemption_receipt": map[string]any{
				"receipt_id":  "receipt:v1:00000000-0000-4000-8000-000000000001",
				"result_kind": "first_committed",
			},
		},
	}
	replay := v1ReplayOutput(first)
	receipt, ok := replay.Data["redemption_receipt"].(map[string]any)
	if !ok {
		t.Fatal("幂等重放缺少兑换回执")
	}
	if receipt["result_kind"] != "idempotent_replay" {
		t.Fatalf("兑换回执重放标识错误: %v", receipt["result_kind"])
	}
	if first.Data["redemption_receipt"].(map[string]any)["result_kind"] != "first_committed" {
		t.Fatal("幂等重放改写了首次成功快照")
	}
}
