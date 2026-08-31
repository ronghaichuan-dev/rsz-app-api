package middleware

import (
	"net/url"
	"reflect"
	"testing"
)

// TestSafeQueryParameters 验证普通查询参数完整保留，可能含凭据的参数统一脱敏。
func TestSafeQueryParameters(t *testing.T) {
	query := url.Values{
		"member_id":      {"member:v1:00000000-0000-4000-8000-000000000001"},
		"identity_token": {"credential-value"},
		"page":           {"1", "2"},
	}
	want := url.Values{
		"member_id":      {"member:v1:00000000-0000-4000-8000-000000000001"},
		"identity_token": {"[REDACTED]"},
		"page":           {"1", "2"},
	}
	if got := safeQueryParameters(query); !reflect.DeepEqual(got, want) {
		t.Fatalf("查询参数脱敏结果不符合预期: got=%v want=%v", got, want)
	}
	if query.Get("identity_token") != "credential-value" {
		t.Fatalf("查询参数脱敏不应修改原始请求值: %v", query)
	}
}
