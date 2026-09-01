package kids

import (
	"context"
	"errors"
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1MemberBalanceDependencyError 验证真实读存储故障返回可重试的 UNAVAILABLE，而不是伪造成功响应。
func TestV1MemberBalanceDependencyError(t *testing.T) {
	err := v1MemberBalanceDependencyError(context.Background(), "balance_snapshot_store", errors.New("store unavailable"))
	protocolErr, ok := err.(*v1.V1Error)
	if !ok {
		t.Fatalf("读存储故障没有转换为 V1Error: %T", err)
	}
	if protocolErr.Status != 503 || protocolErr.Code != "UNAVAILABLE" || !protocolErr.Retryable {
		t.Fatalf("读存储故障错误语义不正确: %+v", protocolErr)
	}
	if protocolErr.RetryAfterMs == nil || *protocolErr.RetryAfterMs <= 0 {
		t.Fatalf("读存储故障缺少合理重试等待时间: %+v", protocolErr)
	}
}

// TestV1MemberBalanceDependencyErrorPreservesProtocolError 验证权限和资源错误不会被误映射为 503。
func TestV1MemberBalanceDependencyErrorPreservesProtocolError(t *testing.T) {
	original := &v1.V1Error{Status: 403, Code: "FORBIDDEN", Retryable: false}
	err := v1MemberBalanceDependencyError(context.Background(), "membership", original)
	if err != original {
		t.Fatal("已有协议错误不应被改写为依赖故障")
	}
}
