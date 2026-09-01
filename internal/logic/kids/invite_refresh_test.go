package kids

import (
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

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

// TestV1CircleMembershipExistsError 验证已加入目标圈子的账号不会被误报为邀请码已使用。
func TestV1CircleMembershipExistsError(t *testing.T) {
	err := v1CircleMembershipExistsError()
	protocolErr, ok := err.(*v1.V1Error)
	if !ok {
		t.Fatalf("圈子成员已存在错误不是 V1Error: %T", err)
	}
	if protocolErr.Status != 409 || protocolErr.Code != "ALREADY_CIRCLE_MEMBER" || protocolErr.Retryable || protocolErr.Version != nil {
		t.Fatalf("圈子成员已存在错误语义不正确: %+v", protocolErr)
	}
}
