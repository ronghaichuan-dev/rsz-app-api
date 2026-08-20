package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// JWTClaims 是本项目 JWT 载荷结构。
type JWTClaims struct {
	UserId uint64 `json:"userId"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// GenerateJWT 使用 HMAC-SHA256 生成 JWT 字符串。
func GenerateJWT(userId uint64, expiresAt int64, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := JWTClaims{UserId: userId, Exp: expiresAt, Iat: time.Now().Unix()}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	headerPart := base64.RawURLEncoding.EncodeToString(headerBytes)
	claimPart := base64.RawURLEncoding.EncodeToString(claimBytes)
	signature := signJWT(headerPart+"."+claimPart, secret)
	return headerPart + "." + claimPart + "." + signature, nil
}

// ParseJWT 校验 JWT 签名和过期时间，并返回载荷。
func ParseJWT(token string, secret string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid token")
	}
	expected := signJWT(parts[0]+"."+parts[1], secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid token signature")
	}
	claimBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid token payload")
	}
	var claims JWTClaims
	if err = json.Unmarshal(claimBytes, &claims); err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid token claims")
	}
	if claims.UserId == 0 || claims.Exp <= time.Now().Unix() {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "token expired")
	}
	return &claims, nil
}

// ExtractBearerToken 从 Authorization 请求头中提取 Bearer Token。
func ExtractBearerToken(authorization string) string {
	authorization = strings.TrimSpace(authorization)
	if authorization == "" {
		return ""
	}
	parts := strings.Fields(authorization)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return authorization
}

// Uint64String 将 uint64 转成字符串，便于日志和上下文使用。
func Uint64String(value uint64) string {
	return strconv.FormatUint(value, 10)
}

// signJWT 计算 JWT 签名部分。
func signJWT(content string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(content))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
