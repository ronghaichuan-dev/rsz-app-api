package login

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// DBProvider 根据当前上下文返回调用方微服务自己的数据库连接。
type DBProvider func(ctx context.Context) gdb.DB

// StringProvider 根据当前上下文和登录提供方返回配置字符串。
type StringProvider func(ctx context.Context, provider string) string

// UserHook 允许调用方在登录事务内追加当前微服务自己的业务写入。
type UserHook func(ctx context.Context, tx gdb.TX, user *UserRecord) error

// Config 定义公共登录模块需要的数据库表、令牌和 OAuth 配置。
type Config struct {
	DB              DBProvider
	UserTable       string
	UserAuthTable   string
	UserTokenTable  string
	AccessTokenTTL  time.Duration
	JWTSecret       func(ctx context.Context) string
	OAuthClientId   StringProvider
	DefaultNickname StringProvider
}

// Input 是公共登录模块的统一入参。
type Input struct {
	Provider      string
	DeviceId      string
	AuthCode      string
	IdentityToken string
	OpenId        string
	Email         string
	Nickname      string
	Avatar        string
}

// Output 是公共登录模块的统一出参。
type Output struct {
	UserId       uint64
	Token        string
	Provider     string
	IsGuest      bool
	BoundGuest   bool
	IsNewUser    bool
	DeviceId     string
	Nickname     string
	Avatar       string
	AccessExpire int64
}

// UserRecord 承载用户表中登录流程需要的通用字段。
type UserRecord struct {
	Id       uint64
	DeviceId string
	Provider string
	Email    string
	Nickname string
	Avatar   string
	IsGuest  uint
}

// OAuthIdentity 承载服务端校验后的授权登录资料。
type OAuthIdentity struct {
	OpenId   string
	Email    string
	Nickname string
	Avatar   string
}

// Options 承载单次登录调用的可选行为。
type Options struct {
	UserHook UserHook
}

// Option 用于定制单次登录事务中的扩展行为。
type Option func(options *Options)

// WithUserHook 设置登录事务内的用户扩展钩子。
func WithUserHook(hook UserHook) Option {
	return func(options *Options) {
		options.UserHook = hook
	}
}

// MergeOptions 合并单次登录选项。
func MergeOptions(opts ...Option) Options {
	var options Options
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

// NormalizeConfig 补齐公共登录模块的默认配置。
func NormalizeConfig(cfg Config) Config {
	if cfg.DB == nil {
		cfg.DB = func(ctx context.Context) gdb.DB {
			return g.DB(consts.DefaultDBGroup)
		}
	}
	if cfg.AccessTokenTTL <= 0 {
		cfg.AccessTokenTTL = consts.DefaultAccessTokenTTL
	}
	if cfg.JWTSecret == nil {
		cfg.JWTSecret = SecretFromConfig
	}
	if cfg.OAuthClientId == nil {
		cfg.OAuthClientId = OAuthClientIdFromConfig
	}
	if cfg.DefaultNickname == nil {
		cfg.DefaultNickname = func(ctx context.Context, provider string) string {
			if provider == consts.LoginProviderGuest {
				return "Guest"
			}
			return utils.DefaultProviderName(provider)
		}
	}
	return cfg
}

// TokenExists 校验 JWT 是否仍存在于访问令牌表且未过期。
func TokenExists(ctx context.Context, db gdb.DB, tokenTable string, userId uint64, token string) bool {
	count, err := db.Model(tokenTable).Ctx(ctx).
		Where("user_id", userId).
		Where("token", token).
		WhereGT("expired_at", time.Now().Format(consts.MySQLTimeLayout)).
		Count()
	return err == nil && count > 0
}

// SecretFromConfig 从配置读取 JWT 密钥，未配置时使用开发默认值。
func SecretFromConfig(ctx context.Context) string {
	secret := strings.TrimSpace(g.Cfg().MustGet(ctx, "auth.jwt.secret", consts.DefaultJWTSecret).String())
	if secret == "" {
		return consts.DefaultJWTSecret
	}
	return secret
}

// OAuthClientIdFromConfig 按登录服务商读取 OAuth 客户端 ID。
func OAuthClientIdFromConfig(ctx context.Context, provider string) string {
	key := ""
	switch provider {
	case consts.LoginProviderGoogle:
		key = "auth.google.clientId"
	case consts.LoginProviderApple:
		key = "auth.apple.clientId"
	}
	if key == "" {
		return ""
	}
	return strings.TrimSpace(g.Cfg().MustGet(ctx, key, "").String())
}
