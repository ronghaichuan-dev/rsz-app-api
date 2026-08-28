package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"

	"rslytics-app-api/internal/consts"
)

// Ctx 初始化请求级 trace ID，供网关、应用和数据库错误关联，且不读取或记录请求正文。
func Ctx(r *ghttp.Request) {
	traceID := strings.TrimSpace(r.Header.Get("X-Trace-Id"))
	if len(traceID) < 8 || len(traceID) > 128 {
		traceID = "trace:" + uuid.NewString()
	}
	r.SetCtxVar(consts.CtxTraceIDKey, traceID)
	r.Middleware.Next()
}
