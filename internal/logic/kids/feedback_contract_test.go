package kids

import (
	"testing"
	"time"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1FeedbackReceiptUsesCanonicalCommitTime 验证首次提交与幂等重放均保留同一反馈接收时间和回执提交时间。
func TestV1FeedbackReceiptUsesCanonicalCommitTime(t *testing.T) {
	committedAt := time.UnixMilli(1_700_000_000_123)
	receipt := map[string]any{
		"receipt_id":      "receipt:v1:00000000-0000-4000-8000-000000000001",
		"commit_id":       "commit:v1:00000000-0000-4000-8000-000000000001",
		"commit_sequence": int64(42),
		"result_kind":     "first_committed",
		"committed_at_ms": committedAt.UnixMilli(),
	}
	data := v1FeedbackReceiptData("feedback:v1:00000000-0000-4000-8000-000000000001", committedAt, receipt)
	if err := v1.ValidateV1ResponseData("submitFeedback", data); err != nil {
		t.Fatalf("反馈成功回执不符合冻结合同: %v", err)
	}
	if data["received_at_ms"] != receipt["committed_at_ms"] {
		t.Fatalf("反馈接收时间与提交回执时间必须完全相等: data=%v receipt=%v", data, receipt)
	}
	replay := v1ReplayOutput(&v1.V1OperationOutput{Data: data, Status: 200})
	replayReceipt := replay.Data["receipt"].(map[string]any)
	if replay.Data["received_at_ms"] != replayReceipt["committed_at_ms"] || replayReceipt["result_kind"] != "idempotent_replay" {
		t.Fatalf("反馈幂等重放没有保留 canonical 回执: %v", replay.Data)
	}
	if err := v1.ValidateV1ResponseData("submitFeedback", replay.Data); err != nil {
		t.Fatalf("反馈幂等重放不符合冻结合同: %v", err)
	}
}
