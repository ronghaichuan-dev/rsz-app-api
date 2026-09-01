package kids

import (
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1RefreshInviteUnavailableError 验证不可刷新的邀请码不会被错误标记为可重试的版本冲突。
func TestV1RefreshInviteUnavailableError(t *testing.T) {
	err := v1RefreshInviteUnavailableError()
	protocolErr, ok := err.(*v1.V1Error)
	if !ok {
		t.Fatalf("邀请不可刷新错误不是 V1Error: %T", err)
	}
	if protocolErr.Status != 409 || protocolErr.Code != "INVITE_USED" || protocolErr.Retryable {
		t.Fatalf("邀请不可刷新错误语义不正确: %+v", protocolErr)
	}
	if protocolErr.Version != nil {
		t.Fatalf("非版本冲突错误不应返回 current_version: %+v", protocolErr)
	}
}

// TestV1InviteVersionConflictError 验证真实版本不匹配会向客户端提供可用于受控重试的正数版本。
func TestV1InviteVersionConflictError(t *testing.T) {
	err := v1InviteVersionConflictError(2)
	protocolErr, ok := err.(*v1.V1Error)
	if !ok {
		t.Fatalf("邀请版本冲突错误不是 V1Error: %T", err)
	}
	if protocolErr.Status != 409 || protocolErr.Code != "VERSION_CONFLICT" || protocolErr.Version == nil || *protocolErr.Version != 2 {
		t.Fatalf("邀请版本冲突未返回当前正数版本: %+v", protocolErr)
	}
}

// TestV1InviteVersionConflictErrorRejectsInvalidVersion 验证损坏版本不会伪造出无法重试的 VERSION_CONFLICT。
func TestV1InviteVersionConflictErrorRejectsInvalidVersion(t *testing.T) {
	err := v1InviteVersionConflictError(0)
	protocolErr, ok := err.(*v1.V1Error)
	if !ok {
		t.Fatalf("无效邀请版本错误不是 V1Error: %T", err)
	}
	if protocolErr.Status != 502 || protocolErr.Code != "PROTOCOL_ERROR" || protocolErr.Version != nil {
		t.Fatalf("无效邀请版本错误语义不正确: %+v", protocolErr)
	}
}
