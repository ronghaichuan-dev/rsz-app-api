package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"rslytics-app-api/internal/consts"
)

// KidsDB 返回当前 kids 微服务配置的默认数据库连接。
func KidsDB(ctx context.Context) gdb.DB {
	return g.DB(consts.DefaultDBGroup)
}

// NormalizeDate 统一日期入参，空值时返回今天日期。
func NormalizeDate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().Format(consts.DateLayout)
	}
	return value
}

// ParseDBTime 将数据库时间字段转换为 Unix 时间戳。
func ParseDBTime(value any) int64 {
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case *gtime.Time:
		if v == nil {
			return 0
		}
		return v.Timestamp()
	case gtime.Time:
		return v.Timestamp()
	case time.Time:
		if v.IsZero() {
			return 0
		}
		return v.Unix()
	case string:
		return parseTimeString(v)
	case []byte:
		return parseTimeString(string(v))
	default:
		return 0
	}
}

// BoolToInt 将布尔值转换为数据库中的 0/1 值。
func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// DefaultString 在字符串为空时返回默认值。
func DefaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

// DefaultProviderName 根据授权提供方生成默认展示名称。
func DefaultProviderName(provider string) string {
	if provider == "" {
		return "OAuth"
	}
	return strings.ToUpper(provider[:1]) + provider[1:]
}

// GenerateAccessToken 生成带用户 ID 和过期时间的随机访问令牌。
func GenerateAccessToken(userId uint64, expiresAt int64) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", gerror.Wrap(err, "generate access token failed")
	}
	return fmt.Sprintf("kids_%d_%d_%s", userId, expiresAt, hex.EncodeToString(buf)), nil
}

// parseTimeString 按常见数据库时间格式解析字符串时间。
func parseTimeString(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "0000-00-00") {
		return 0
	}
	for _, layout := range []string{consts.MySQLTimeLayout, time.RFC3339, consts.DateLayout} {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t.Unix()
		}
	}
	return 0
}
