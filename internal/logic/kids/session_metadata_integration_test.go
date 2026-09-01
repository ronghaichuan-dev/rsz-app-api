package kids

import (
	"reflect"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
)

// TestV1SessionMetadataIntegrationExchangeRefreshAndBootstrap 验证 exchange、refresh 和 bearer token bootstrap 都投影对应持久化 session 的完整毫秒 metadata。
func TestV1SessionMetadataIntegrationExchangeRefreshAndBootstrap(t *testing.T) {
	issuedAt := time.UnixMilli(1787900222276)
	persistedSession := v1TestPersistedAccountSession("session:v1:00000000-0000-4000-8000-000000000001", issuedAt)

	// exchange 在写入 session 后从持久化记录构造响应；access token 仅用于后续 bearer 鉴权。
	exchangeResponse := v1AuthSessionProjection(persistedSession, "access-token", "refresh-token")
	exchangeMetadata, ok := exchangeResponse["metadata"].(map[string]any)
	if !ok {
		t.Fatal("exchange 响应缺少 session metadata")
	}

	// bootstrap 通过 access token 找到同一条 session 后，直接使用相同投影函数返回 metadata。
	bootstrapMetadata := v1SessionMetadataProjection(persistedSession)
	expected := v1ExpectedSessionMetadata("session:v1:00000000-0000-4000-8000-000000000001", issuedAt)
	if !reflect.DeepEqual(exchangeMetadata, expected) {
		t.Fatalf("exchange session metadata 不符合持久化记录: got=%v want=%v", exchangeMetadata, expected)
	}
	if !reflect.DeepEqual(bootstrapMetadata, expected) {
		t.Fatalf("bootstrap session metadata 不符合持久化记录: got=%v want=%v", bootstrapMetadata, expected)
	}
	if issuedAt.UnixMilli()%1000 == 0 {
		t.Fatal("测试时间必须包含非零毫秒")
	}

	// refresh 在同一 session 上原地轮换 token pair；后续 bootstrap 必须读取同一条记录的最新 metadata。
	rotatedIssuedAt := time.UnixMilli(1787900222377)
	rotatedSession := v1TestPersistedAccountSession("session:v1:00000000-0000-4000-8000-000000000001", rotatedIssuedAt)
	refreshResponse := v1AuthSessionProjection(rotatedSession, "rotated-access-token", "rotated-refresh-token")
	refreshMetadata, ok := refreshResponse["metadata"].(map[string]any)
	if !ok {
		t.Fatal("refresh 响应缺少 session metadata")
	}
	rotatedExpected := v1ExpectedSessionMetadata("session:v1:00000000-0000-4000-8000-000000000001", rotatedIssuedAt)
	if !reflect.DeepEqual(refreshMetadata, rotatedExpected) {
		t.Fatalf("refresh session metadata 不符合新持久化记录: got=%v want=%v", refreshMetadata, rotatedExpected)
	}
	if bootstrapAfterRefresh := v1SessionMetadataProjection(rotatedSession); !reflect.DeepEqual(bootstrapAfterRefresh, rotatedExpected) {
		t.Fatalf("refresh 后 bootstrap session metadata 不符合新持久化记录: got=%v want=%v", bootstrapAfterRefresh, rotatedExpected)
	}
}

// v1TestPersistedAccountSession 构造数据库读取后的账号 session 记录，时间值固定为带非零毫秒的值。
func v1TestPersistedAccountSession(sessionID string, issuedAt time.Time) gdb.Record {
	return gdb.Record{
		"session_id":            gvar.New(sessionID),
		"principal_kind":        gvar.New("account"),
		"status":                gvar.New("active"),
		"issued_at_ms":          gvar.New(issuedAt.UnixMilli()),
		"access_expires_at_ms":  gvar.New(issuedAt.Add(time.Hour).UnixMilli()),
		"refresh_expires_at_ms": gvar.New(issuedAt.Add(30 * 24 * time.Hour).UnixMilli()),
	}
}

// v1ExpectedSessionMetadata 返回固定持久化 session 应对外返回的完整 metadata。
func v1ExpectedSessionMetadata(sessionID string, issuedAt time.Time) map[string]any {
	return map[string]any{
		"session_id":            sessionID,
		"principal_type":        "account",
		"status":                "active",
		"issued_at_ms":          issuedAt.UnixMilli(),
		"access_expires_at_ms":  issuedAt.Add(time.Hour).UnixMilli(),
		"refresh_expires_at_ms": issuedAt.Add(30 * 24 * time.Hour).UnixMilli(),
	}
}
