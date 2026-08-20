package kids

import (
	"testing"
	"time"
)

// TestV1SyncChanges 验证任意业务提交都会被规整为完整同步变更集合。
func TestV1SyncChanges(t *testing.T) {
	changes := v1SyncChanges(map[string]any{
		"circle":           map[string]any{"circle_id": "circle:v1:00000000-0000-4000-8000-000000000001"},
		"ledger_entry":     map[string]any{"ledger_id": "ledger:v1:00000000-0000-4000-8000-000000000001"},
		"reward_tombstone": map[string]any{"entity_id": "reward:v1:00000000-0000-4000-8000-000000000001"},
	})
	if len(changes) != 19 {
		t.Fatalf("同步变更字段数量不正确: %d", len(changes))
	}
	if got := len(changes["circles"].([]any)); got != 1 {
		t.Fatalf("圈子变更数量不正确: %d", got)
	}
	if got := len(changes["ledger_entries"].([]any)); got != 1 {
		t.Fatalf("流水变更数量不正确: %d", got)
	}
	if got := len(changes["tombstones"].([]any)); got != 1 {
		t.Fatalf("删除标记数量不正确: %d", got)
	}
}

// TestV1SyncSequence 验证同步游标只能解析本服务签发的十六进制提交序列。
func TestV1SyncSequence(t *testing.T) {
	sequence, err := v1SyncSequence(v1CommitCursor(42))
	if err != nil || sequence != 42 {
		t.Fatalf("同步游标解析失败: sequence=%d err=%v", sequence, err)
	}
	if _, err = v1SyncSequence("cur:v1:invalid"); err == nil {
		t.Fatal("无效同步游标不应被接受")
	}
}

// TestV1RewardCooldown 验证奖励重复规则不会把非一次性奖励误判为永久不可用。
func TestV1RewardCooldown(t *testing.T) {
	now := mustTestTime(t)
	permanent, until := v1RewardCooldown("none", 0, now)
	if !permanent || !until.IsZero() {
		t.Fatal("一次性奖励冷却计算不正确")
	}
	permanent, until = v1RewardCooldown("daily", 0, now)
	if permanent || !until.Equal(now.AddDate(0, 0, 1)) {
		t.Fatal("每日奖励冷却计算不正确")
	}
}

// mustTestTime 返回测试中固定可比较的时间值。
func mustTestTime(t *testing.T) time.Time {
	t.Helper()
	value, err := time.Parse(time.RFC3339, "2026-08-20T00:00:00Z")
	if err != nil {
		t.Fatalf("测试时间解析失败: %v", err)
	}
	return value
}
