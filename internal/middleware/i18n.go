package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	commoni18n "rslytics-app-api/internal/common/i18n"
)

// I18n 从请求头读取语言并写入请求上下文。
func I18n(r *ghttp.Request) {
	commoni18n.WithRequestLanguage(r)
	r.Middleware.Next()
}
