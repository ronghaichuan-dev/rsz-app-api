package middleware

import (
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"

	"rslytics-app-api/internal/consts"
)

// Ctx 初始化请求级 trace ID，并在请求结束后记录脱敏的路径和查询参数，不读取或记录请求正文。
func Ctx(r *ghttp.Request) {
	traceID := strings.TrimSpace(r.Header.Get("X-Trace-Id"))
	if len(traceID) < 8 || len(traceID) > 128 {
		traceID = "trace:" + uuid.NewString()
	}
	r.SetCtxVar(consts.CtxTraceIDKey, traceID)
	r.Middleware.Next()
	logRequestParameters(r)
}

// logRequestParameters 统一记录所有 HTTP 接口的可安全诊断参数，敏感查询值会脱敏且请求正文永不写入日志。
func logRequestParameters(r *ghttp.Request) {
	g.Log().Infof(r.Context(), "event=http_request_parameters operation_id=%s request_id=%s trace_id=%s method=%s path=%s path_parameters=%v query_parameters=%v response_status=%d", r.GetCtxVar(consts.CtxV1OperationIDKey).String(), r.Header.Get("X-Request-Id"), r.GetCtxVar(consts.CtxTraceIDKey).String(), r.Method, r.URL.Path, r.GetRouterMap(), safeQueryParameters(r.URL.Query()), r.Response.Status)
}

// safeQueryParameters 返回可记录的查询参数副本，避免未来接口把 credential、token 或密码写入普通日志。
func safeQueryParameters(query url.Values) url.Values {
	result := make(url.Values, len(query))
	for key, values := range query {
		if isSensitiveParameter(key) {
			result[key] = []string{"[REDACTED]"}
			continue
		}
		result[key] = append([]string(nil), values...)
	}
	return result
}

// isSensitiveParameter 判断参数名是否可能承载 credential 或其他敏感凭据。
func isSensitiveParameter(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	return strings.Contains(normalized, "token") || strings.Contains(normalized, "secret") || strings.Contains(normalized, "password") || strings.Contains(normalized, "credential") || strings.Contains(normalized, "authorization") || strings.Contains(normalized, "proof")
}
