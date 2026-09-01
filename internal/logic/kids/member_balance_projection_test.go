package kids

import (
	"reflect"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
)

// TestV1BalanceProjectionFromRecord 验证变更响应、余额读取和同步提交都只使用同一持久化余额行的规范字段。
func TestV1BalanceProjectionFromRecord(t *testing.T) {
	updatedAt := time.UnixMilli(1_700_000_000_123)
	row := gdb.Record{
		"circle_id":              gvar.New("circle:v1:00000000-0000-4000-8000-000000000001"),
		"member_id":              gvar.New("member:v1:00000000-0000-4000-8000-000000000001"),
		"balance":                gvar.New(9),
		"version":                gvar.New(4),
		"source_commit_id":       gvar.New("commit:v1:00000000-0000-4000-8000-000000000004"),
		"source_commit_sequence": gvar.New(44),
		"updated_at":             gvar.New(updatedAt),
	}

	mutationBalance := v1BalanceProjectionFromRecord(row)
	readBalance := v1BalanceProjectionFromRecord(row)
	syncBalance := v1SyncChanges(map[string]any{"balance": v1BalanceProjectionFromRecord(row)})["star_balances"].([]any)[0].(map[string]any)

	if !reflect.DeepEqual(mutationBalance, readBalance) || !reflect.DeepEqual(mutationBalance, syncBalance) {
		t.Fatalf("同一持久化余额行生成了不同规范投影: response=%v read=%v sync=%v", mutationBalance, readBalance, syncBalance)
	}
	if mutationBalance["updated_at_ms"] != updatedAt.UnixMilli() || mutationBalance["source_commit_sequence"] != int64(44) {
		t.Fatalf("持久化余额字段未被完整保留: %v", mutationBalance)
	}
}
