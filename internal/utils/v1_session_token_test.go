package utils

import (
	"testing"
	"time"
)

// TestV1SessionToken 验证 v1 access token 绑定签发密钥、部署环境、session 和毫秒级有效期。
func TestV1SessionToken(t *testing.T) {
	now := time.Now().UnixMilli()
	claims := V1SessionTokenClaims{
		SessionID: "session:v1:00000000-0000-4000-8000-000000000001", AccountID: "account:v1:00000000-0000-4000-8000-000000000001", PrincipalKind: "account", Environment: "test", IssuedAtMs: now, ExpiresAtMs: now + 60_000,
	}
	token, err := GenerateV1SessionToken(claims, "test-signing-secret")
	if err != nil {
		t.Fatalf("签发 v1 session token 失败: %v", err)
	}
	parsed, err := ParseV1SessionToken(token, "test-signing-secret", "test")
	if err != nil || parsed.SessionID != claims.SessionID || parsed.ExpiresAtMs != claims.ExpiresAtMs {
		t.Fatalf("签发后的 v1 session token 无法立即验签: claims=%+v err=%v", parsed, err)
	}
	if _, err = ParseV1SessionToken(token, "test-signing-secret", "prod"); err == nil {
		t.Fatal("不同部署环境的 token 不应通过验签")
	}
	if _, err = ParseV1SessionToken(token, "wrong-signing-secret", "test"); err == nil {
		t.Fatal("不同签发密钥的 token 不应通过验签")
	}
}
