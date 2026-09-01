package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// V1SessionTokenClaims 是 Clearwave v1 access token 的最小可验证声明。
type V1SessionTokenClaims struct {
	SessionID     string `json:"sid"`
	AccountID     string `json:"aid"`
	PrincipalKind string `json:"pk"`
	Environment   string `json:"env"`
	IssuedAtMs    int64  `json:"iat"`
	ExpiresAtMs   int64  `json:"exp"`
}

// GenerateV1SessionToken 使用当前环境专用密钥签发 v1 access token，所有时间均为 epoch milliseconds。
func GenerateV1SessionToken(claims V1SessionTokenClaims, secret string) (string, error) {
	if strings.TrimSpace(claims.SessionID) == "" || strings.TrimSpace(claims.PrincipalKind) == "" || strings.TrimSpace(claims.Environment) == "" || claims.IssuedAtMs <= 0 || claims.ExpiresAtMs <= claims.IssuedAtMs || strings.TrimSpace(secret) == "" {
		return "", errors.New("v1 session token claims are invalid")
	}
	headerPart := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadPart := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := headerPart + "." + payloadPart
	return signingInput + "." + signV1SessionToken(signingInput, secret), nil
}

// ParseV1SessionToken 校验 v1 access token 的签名、环境和毫秒级到期时间。
func ParseV1SessionToken(token, secret, environment string) (*V1SessionTokenClaims, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 || strings.TrimSpace(secret) == "" || strings.TrimSpace(environment) == "" {
		return nil, errors.New("v1 session token is invalid")
	}
	if !hmac.Equal([]byte(signV1SessionToken(parts[0]+"."+parts[1], secret)), []byte(parts[2])) {
		return nil, errors.New("v1 session token signature is invalid")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("v1 session token payload is invalid")
	}
	var claims V1SessionTokenClaims
	if err = json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("v1 session token claims are invalid")
	}
	if claims.SessionID == "" || claims.PrincipalKind == "" || claims.Environment != environment || claims.IssuedAtMs <= 0 || claims.ExpiresAtMs <= claims.IssuedAtMs || claims.ExpiresAtMs <= time.Now().UnixMilli() {
		return nil, errors.New("v1 session token is expired or invalid")
	}
	return &claims, nil
}

// signV1SessionToken 计算 v1 session token 的 HMAC 签名。
func signV1SessionToken(content, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(content))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
