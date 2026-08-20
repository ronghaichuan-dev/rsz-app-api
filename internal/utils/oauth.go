package utils

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	// googleJWKSURL 是 Google OpenID Connect 公钥地址。
	googleJWKSURL = "https://www.googleapis.com/oauth2/v3/certs"
	// appleJWKSURL 是 Apple Sign in 公钥地址。
	appleJWKSURL = "https://appleid.apple.com/auth/keys"
	// googleIssuerWithScheme 是 Google identityToken 带接口的签发方。
	googleIssuerWithScheme = "https://accounts.google.com"
	// googleIssuerWithoutScheme 是 Google identityToken 不带接口的签发方。
	googleIssuerWithoutScheme = "accounts.google.com"
	// appleIssuer 是 Apple identityToken 的签发方。
	appleIssuer = "https://appleid.apple.com"
)

// OAuthIdentity 是第三方 identityToken 校验后的授权身份。
type OAuthIdentity struct {
	OpenId   string
	Email    string
	Nickname string
	Avatar   string
}

// VerifyGoogleIdentityTokenWithNonce 校验 Google identityToken，并验证调用方绑定的 nonce。
func VerifyGoogleIdentityTokenWithNonce(ctx context.Context, identityToken, clientId, nonce string) (*OAuthIdentity, error) {
	return verifyOIDCIdentityTokenWithNonce(ctx, identityToken, clientId, []string{googleIssuerWithScheme, googleIssuerWithoutScheme}, googleJWKSURL, nonce)
}

// VerifyGoogleIdentityToken 校验 Google identityToken 并返回服务端可信身份。
func VerifyGoogleIdentityToken(ctx context.Context, identityToken, clientId string) (*OAuthIdentity, error) {
	return verifyOIDCIdentityTokenWithNonce(ctx, identityToken, clientId, []string{googleIssuerWithScheme, googleIssuerWithoutScheme}, googleJWKSURL, "")
}

// VerifyAppleIdentityToken 校验 Apple identityToken 并返回服务端可信身份。
func VerifyAppleIdentityToken(ctx context.Context, identityToken, clientId string) (*OAuthIdentity, error) {
	return verifyOIDCIdentityTokenWithNonce(ctx, identityToken, clientId, []string{appleIssuer}, appleJWKSURL, "")
}

// verifyOIDCIdentityToken 按 OpenID Connect 规则校验 identityToken。
func verifyOIDCIdentityToken(ctx context.Context, identityToken, audience string, issuers []string, jwksURL string) (*OAuthIdentity, error) {
	return verifyOIDCIdentityTokenWithNonce(ctx, identityToken, audience, issuers, jwksURL, "")
}

// verifyOIDCIdentityTokenWithNonce 按 OpenID Connect 规则校验签名、签发方、受众、到期和可选 nonce。
func verifyOIDCIdentityTokenWithNonce(ctx context.Context, identityToken, audience string, issuers []string, jwksURL, nonce string) (*OAuthIdentity, error) {
	identityToken = strings.TrimSpace(identityToken)
	audience = strings.TrimSpace(audience)
	if identityToken == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "identityToken is required")
	}
	if audience == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidConfiguration, "oauth clientId is required")
	}
	header, payload, signingInput, signature, err := parseOIDCToken(identityToken)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(header.Alg, "RS256") {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "unsupported identityToken algorithm")
	}
	publicKey, err := fetchOIDCPublicKey(ctx, jwksURL, header.Kid)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256([]byte(signingInput))
	if err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken signature invalid")
	}
	if !matchIssuer(payload.Issuer, issuers) {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken issuer invalid")
	}
	if !payload.MatchAudience(audience) {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken audience invalid")
	}
	if payload.ExpiresAt <= time.Now().Unix() {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken expired")
	}
	if strings.TrimSpace(payload.Subject) == "" {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken subject invalid")
	}
	if nonce != "" && payload.Nonce != nonce {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken nonce invalid")
	}
	return &OAuthIdentity{
		OpenId:   payload.Subject,
		Email:    payload.Email,
		Nickname: payload.Name,
		Avatar:   payload.Picture,
	}, nil
}

// matchIssuer 判断 token 签发方是否在允许列表内。
func matchIssuer(issuer string, allowed []string) bool {
	for _, item := range allowed {
		if issuer == item {
			return true
		}
	}
	return false
}

// oidcTokenHeader 承载 identityToken 头部中用于验签的字段。
type oidcTokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

// oidcTokenPayload 承载 identityToken 声明字段。
type oidcTokenPayload struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  any    `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Nonce     string `json:"nonce"`
}

// MatchAudience 判断 token 的 aud 是否包含指定客户端 ID。
func (p oidcTokenPayload) MatchAudience(audience string) bool {
	switch aud := p.Audience.(type) {
	case string:
		return aud == audience
	case []any:
		for _, item := range aud {
			if value, ok := item.(string); ok && value == audience {
				return true
			}
		}
	}
	return false
}

// parseOIDCToken 解析 identityToken 的头、载荷、签名和待签名内容。
func parseOIDCToken(identityToken string) (oidcTokenHeader, oidcTokenPayload, string, []byte, error) {
	parts := strings.Split(identityToken, ".")
	if len(parts) != 3 {
		return oidcTokenHeader{}, oidcTokenPayload{}, "", nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken invalid")
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return oidcTokenHeader{}, oidcTokenPayload{}, "", nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken header invalid")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return oidcTokenHeader{}, oidcTokenPayload{}, "", nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken payload invalid")
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return oidcTokenHeader{}, oidcTokenPayload{}, "", nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken signature invalid")
	}
	var header oidcTokenHeader
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return oidcTokenHeader{}, oidcTokenPayload{}, "", nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken header invalid")
	}
	var payload oidcTokenPayload
	if err = json.Unmarshal(payloadBytes, &payload); err != nil {
		return oidcTokenHeader{}, oidcTokenPayload{}, "", nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken payload invalid")
	}
	return header, payload, parts[0] + "." + parts[1], signature, nil
}

// oidcJWKS 承载 OpenID Connect 公钥列表响应。
type oidcJWKS struct {
	Keys []oidcJWK `json:"keys"`
}

// oidcJWK 承载 RSA 公钥参数。
type oidcJWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// fetchOIDCPublicKey 拉取并解析第三方平台公钥。
func fetchOIDCPublicKey(ctx context.Context, jwksURL, kid string) (*rsa.PublicKey, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, gerror.Wrap(err, "fetch identityToken public key failed")
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "fetch identityToken public key failed")
	}
	var jwks oidcJWKS
	if err = json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken public key invalid")
	}
	for _, key := range jwks.Keys {
		if key.Kid == kid && strings.EqualFold(key.Kty, "RSA") {
			return jwkToRSAPublicKey(key)
		}
	}
	return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken public key not found")
}

// jwkToRSAPublicKey 将 JWK 参数转换为 RSA 公钥。
func jwkToRSAPublicKey(key oidcJWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken public key invalid")
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken public key invalid")
	}
	exponent := 0
	for _, b := range eBytes {
		exponent = exponent<<8 + int(b)
	}
	if exponent == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "identityToken public key invalid")
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: exponent}, nil
}
