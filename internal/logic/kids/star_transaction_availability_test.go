package kids

import (
	"context"
	"errors"
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1StarTransactionDependencyError 验证流水读取依赖故障会返回可重试错误，而不是错误映射成协议故障。
func TestV1StarTransactionDependencyError(t *testing.T) {
	err := v1StarTransactionDependencyError(context.Background(), "ledger_read_store", errors.New("store unavailable"))
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

// TestV1StarTransactionDependencyErrorPreservesProtocolError 验证权限与参数错误不会被误写为依赖故障。
func TestV1StarTransactionDependencyErrorPreservesProtocolError(t *testing.T) {
	original := &v1.V1Error{Status: 403, Code: "FORBIDDEN", Retryable: false}
	if err := v1StarTransactionDependencyError(context.Background(), "membership", original); err != original {
		t.Fatal("已有协议错误不应被改写为依赖故障")
	}
}
