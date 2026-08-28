package middleware

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
)

// V1Envelope 为 Clearwave v1 路由输出冻结的 success/error envelope。
func V1Envelope(r *ghttp.Request) {
	r.Middleware.Next()
	if !isV1Route(r.URL.Path) || r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}
	if err := r.GetError(); err != nil {
		var v1Err *v1.V1Error
		if errors.As(err, &v1Err) {
			writeV1Error(r, v1Err)
			return
		}
		writeV1Error(r, &v1.V1Error{Status: 502, Code: "PROTOCOL_ERROR", Retryable: false})
		return
	}
	if response, ok := r.GetHandlerResponse().(*v1.V1Response); ok {
		writeV1Headers(r, response.RequestID)
		r.Response.Status = 200
		r.Response.WriteJson(response)
		return
	}
	writeV1Error(r, &v1.V1Error{Status: 502, Code: "PROTOCOL_ERROR", Retryable: false})
}

// writeV1Error 只输出 OpenAPI 冻结的错误 envelope，不泄漏内部异常。
func writeV1Error(r *ghttp.Request, v1Err *v1.V1Error) {
	if v1Err.Code == "UNAVAILABLE" && v1Err.Retryable && v1Err.RetryAfterMs == nil {
		retryAfterMs := int64(1000)
		v1Err.RetryAfterMs = &retryAfterMs
	}
	writeV1Headers(r, r.Header.Get(v1.V1RequestIDHeader))
	r.Response.Status = v1Err.Status
	if v1Err.RetryAfterMs != nil {
		r.Response.Header().Set("Retry-After", strconv.FormatInt((*v1Err.RetryAfterMs+999)/1000, 10))
	}
	logV1Failure(r, v1Err.Code)
	r.Response.WriteJson(map[string]any{
		"contract_version": v1.V1Version,
		"request_id":       r.Header.Get(v1.V1RequestIDHeader),
		"server_time_ms":   time.Now().UnixMilli(),
		"error": map[string]any{
			"code": v1Err.Code, "retryable": v1Err.Retryable,
			"retry_after_ms": v1Err.RetryAfterMs, "field": v1Err.Field,
			"current_version": v1Err.Version, "trace_id": v1TraceID(r),
		},
	})
}

// writeV1Headers 写入每个接口成功或失败响应都必须回显的冻结接口头。
func writeV1Headers(r *ghttp.Request, requestID string) {
	r.Response.Header().Set(v1.V1VersionHeader, v1.V1Version)
	r.Response.Header().Set(v1.V1RequestIDHeader, requestID)
	r.Response.Header().Set("X-Trace-Id", v1TraceID(r))
}

// v1TraceID 返回由请求中间件生成的受控 trace ID。
func v1TraceID(r *ghttp.Request) string {
	return r.GetCtxVar(consts.CtxTraceIDKey).String()
}

// logV1Failure 只记录允许的诊断字段，避免把 proof、token 或用户输入写入日志。
func logV1Failure(r *ghttp.Request, code string) {
	g.Log().Warningf(r.Context(), "v1 请求失败 operation_id=%s status=%d error_code=%s request_id=%s trace_id=%s", r.GetCtxVar(consts.CtxV1OperationIDKey).String(), r.Response.Status, code, r.Header.Get(v1.V1RequestIDHeader), v1TraceID(r))
}

// isV1Route 排除既有 kids、health 和文档路由，只匹配接口根路由。
func isV1Route(path string) bool {
	if strings.HasPrefix(path, "/v1/kids/") || path == "/v1/health" {
		return false
	}
	return strings.HasPrefix(path, "/v1/")
}
