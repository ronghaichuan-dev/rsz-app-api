package v1

import "testing"

// TestMemberBalanceSnapshotAcceptsCanonicalZeroBalance 验证零余额也必须带可用于后续乐观锁调整的规范版本和提交来源。
func TestMemberBalanceSnapshotAcceptsCanonicalZeroBalance(t *testing.T) {
	data := memberBalanceSnapshotFixture("commit:v1:00000000-0000-4000-8000-000000000001", 1)
	if err := ValidateV1ResponseData("getMemberBalances", data); err != nil {
		t.Fatalf("规范零余额快照被拒绝: %v", err)
	}
}

// TestMemberBalanceSnapshotRejectsLegacyBackfillSource 防止非规范回填来源在运行时被折叠为 UNAVAILABLE。
func TestMemberBalanceSnapshotRejectsLegacyBackfillSource(t *testing.T) {
	data := memberBalanceSnapshotFixture("migration_star_balance_backfill", 0)
	if err := ValidateV1ResponseData("getMemberBalances", data); err == nil {
		t.Fatal("非规范的历史余额回填来源不应通过响应合同")
	}
}

// memberBalanceSnapshotFixture 构造符合读取接口最小字段集的单成员余额快照。
func memberBalanceSnapshotFixture(sourceCommitID string, sourceCommitSequence int64) map[string]any {
	return map[string]any{
		"items": []map[string]any{{
			"circle_id":              "circle:v1:00000000-0000-4000-8000-000000000001",
			"member_id":              "member:v1:00000000-0000-4000-8000-000000000001",
			"balance":                int64(0),
			"version":                int64(1),
			"source_commit_id":       sourceCommitID,
			"source_commit_sequence": sourceCommitSequence,
			"updated_at_ms":          int64(1),
		}},
		"snapshot_cursor": "cur:v1:00000001",
	}
}
