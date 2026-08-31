package v1

import "testing"

// TestStarTransactionPageAcceptsEmptySnapshot 验证没有流水时仍返回可解码的规范空页。
func TestStarTransactionPageAcceptsEmptySnapshot(t *testing.T) {
	data := map[string]any{
		"items":           []any{},
		"next_cursor":     nil,
		"has_more":        false,
		"snapshot_cursor": "cur:v1:00000000",
	}
	if err := ValidateV1ResponseData("listStarTransactions", data); err != nil {
		t.Fatalf("合法空流水页被合同拒绝: %v", err)
	}
}

// TestStarTransactionPageRejectsLegacyLedgerID 防止历史 ledger:v1 标识再次在读取期导致协议失败。
func TestStarTransactionPageRejectsLegacyLedgerID(t *testing.T) {
	data := map[string]any{
		"items": []any{map[string]any{
			"ledger_id":             "ledger:v1:00000000-0000-4000-8000-000000000001",
			"circle_id":             "circle:v1:00000000-0000-4000-8000-000000000001",
			"member_id":             "member:v1:00000000-0000-4000-8000-000000000001",
			"source":                map[string]any{"source_type": "adjustment", "source_id": "adjustment:v1:00000000-0000-4000-8000-000000000001", "title_snapshot": nil, "stars_snapshot": nil, "asset_id_snapshot": nil, "scheduled_date_snapshot": nil},
			"delta":                 1,
			"reason":                nil,
			"actor":                 map[string]any{"actor_type": "owner", "actor_id": "admin:v1:00000000-0000-4000-8000-000000000001", "role": "owner", "display_name_snapshot": "管理员"},
			"reversal_of_ledger_id": nil,
			"created_at_ms":         int64(1_700_000_000_000),
			"commit_sequence":       int64(1),
		}},
		"next_cursor":     nil,
		"has_more":        false,
		"snapshot_cursor": "cur:v1:00000001",
	}
	if err := ValidateV1ResponseData("listStarTransactions", data); err == nil {
		t.Fatal("旧 ledger 命名空间不应通过流水页合同")
	}
}

// TestStarTransactionPageAcceptsCanonicalLedgerID 验证已有规范流水可与分页字段一起被客户端解码。
func TestStarTransactionPageAcceptsCanonicalLedgerID(t *testing.T) {
	data := map[string]any{
		"items": []any{map[string]any{
			"ledger_id":             "star-transaction:v1:00000000-0000-4000-8000-000000000001",
			"circle_id":             "circle:v1:00000000-0000-4000-8000-000000000001",
			"member_id":             "member:v1:00000000-0000-4000-8000-000000000001",
			"source":                map[string]any{"source_type": "adjustment", "source_id": "adjustment:v1:00000000-0000-4000-8000-000000000001", "title_snapshot": nil, "stars_snapshot": nil, "asset_id_snapshot": nil, "scheduled_date_snapshot": nil},
			"delta":                 1,
			"reason":                nil,
			"actor":                 map[string]any{"actor_type": "owner", "actor_id": "admin:v1:00000000-0000-4000-8000-000000000001", "role": "owner", "display_name_snapshot": "管理员"},
			"reversal_of_ledger_id": nil,
			"created_at_ms":         int64(1_700_000_000_000),
			"commit_sequence":       int64(1),
		}},
		"next_cursor":     nil,
		"has_more":        false,
		"snapshot_cursor": "cur:v1:00000001",
	}
	if err := ValidateV1ResponseData("listStarTransactions", data); err != nil {
		t.Fatalf("规范流水页被合同拒绝: %v", err)
	}
}
