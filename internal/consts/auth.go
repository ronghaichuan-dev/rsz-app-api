package consts

import "time"

const (
	// CtxUserIdKey 是请求上下文中的用户ID键。
	CtxUserIdKey = "userId"
	// CtxTokenKey 是请求上下文中的令牌键。
	CtxTokenKey = "token"
	// DefaultJWTSecret 是未配置 JWT 密钥时的开发环境默认密钥。
	DefaultJWTSecret = "rslytics-kids-dev-secret"
	// DefaultAccessTokenTTL 是公共登录访问令牌默认有效期。
	DefaultAccessTokenTTL = 7 * 24 * time.Hour
)
