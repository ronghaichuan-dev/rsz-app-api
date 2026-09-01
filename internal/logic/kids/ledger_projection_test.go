package kids

import (
	"reflect"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
)

// TestV1LedgerProjectionFromRecord 验证变更响应、流水列表和同步提交都只使用同一不可变流水行的规范字段。
func TestV1LedgerProjectionFromRecord(t *testing.T) {
	createdAt := time.UnixMilli(1_700_000_000_123)
	row := gdb.Record{
		"ledger_id":             gvar.New("star-transaction:v1:00000000-0000-4000-8000-000000000001"),
		"circle_id":             gvar.New("circle:v1:00000000-0000-4000-8000-000000000001"),
		"member_id":             gvar.New("member:v1:00000000-0000-4000-8000-000000000001"),
		"source":                gvar.New(`{"source_type":"adjustment","source_id":"adjustment:v1:00000000-0000-4000-8000-000000000001","title_snapshot":null,"stars_snapshot":null,"asset_id_snapshot":null,"scheduled_date_snapshot":null}`),
		"delta":                 gvar.New(-3),
		"reason":                gvar.New("测试反向调整"),
		"actor":                 gvar.New(`{"actor_type":"owner","actor_id":"admin:v1:00000000-0000-4000-8000-000000000001","role":"owner","display_name_snapshot":"管理员"}`),
		"reversal_of_ledger_id": gvar.New(""),
		"commit_sequence":       gvar.New(44),
		"created_at":            gvar.New(createdAt),
	}

	responseLedger := v1LedgerRecordProjection(row)
	listLedger := v1LedgerRecordProjection(row)
	syncLedger := v1SyncChanges(map[string]any{"ledger_entry": v1LedgerRecordProjection(row)})["ledger_entries"].([]any)[0].(map[string]any)

	if !reflect.DeepEqual(responseLedger, listLedger) || !reflect.DeepEqual(responseLedger, syncLedger) {
		t.Fatalf("同一不可变流水行生成了不同规范投影: response=%v list=%v sync=%v", responseLedger, listLedger, syncLedger)
	}
	if responseLedger["created_at_ms"] != createdAt.UnixMilli() || responseLedger["commit_sequence"] != int64(44) {
		t.Fatalf("持久化流水字段未被完整保留: %v", responseLedger)
	}
}
