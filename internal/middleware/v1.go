package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "rslytics-app-api/internal/api/kids/v1"
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
		writeV1Error(r, &v1.V1Error{Status: 503, Code: "UNAVAILABLE", Retryable: true})
		return
	}
	if response, ok := r.GetHandlerResponse().(*v1.V1Response); ok {
		writeV1Headers(r, response.RequestID)
		r.Response.Status = 200
		r.Response.WriteJson(response)
		return
	}
	writeV1Error(r, &v1.V1Error{Status: 503, Code: "UNAVAILABLE", Retryable: true})
}

// writeV1Error 只输出 OpenAPI 冻结的错误 envelope，不泄漏内部异常。
func writeV1Error(r *ghttp.Request, v1Err *v1.V1Error) {
	writeV1Headers(r, r.Header.Get(v1.V1RequestIDHeader))
	r.Response.Status = v1Err.Status
	r.Response.WriteJson(map[string]any{
		"contract_version": v1.V1Version,
		"request_id":       r.Header.Get(v1.V1RequestIDHeader),
		"server_time_ms":   time.Now().UnixMilli(),
		"error": map[string]any{
			"code": v1Err.Code, "retryable": v1Err.Retryable,
			"retry_after_ms": v1Err.RetryAfterMs, "field": v1Err.Field,
			"current_version": v1Err.Version, "trace_id": nil,
		},
	})
}

// writeV1Headers 写入每个合同成功或失败响应都必须回显的冻结协议头。
func writeV1Headers(r *ghttp.Request, requestID string) {
	r.Response.Header().Set(v1.V1VersionHeader, v1.V1Version)
	r.Response.Header().Set(v1.V1RequestIDHeader, requestID)
}

// isV1Route 排除既有 kids、health 和文档路由，只匹配合同根路由。
func isV1Route(path string) bool {
	if strings.HasPrefix(path, "/v1/kids/") || path == "/v1/health" {
		return false
	}
	return strings.HasPrefix(path, "/v1/")
}
