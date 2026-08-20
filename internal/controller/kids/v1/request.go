package v1

import (
	"context"
	"mime"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// RequestInput 在 Controller 层严格解析一个接口路由的 HTTP 参数和 JSON 请求体。
func RequestInput(ctx context.Context, operationID, method string) (v1.V1OperationInput, error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return v1.V1OperationInput{}, &v1.V1Error{Status: 422, Code: "VALIDATION_FAILED", Message: "v1 request is missing"}
	}
	body, bodyPresent, err := v1.DecodeV1Body(r.GetBody())
	if err != nil {
		return v1.V1OperationInput{}, &v1.V1Error{Status: 422, Code: "VALIDATION_FAILED", Message: "request body is invalid"}
	}
	if bodyPresent {
		contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if !strings.EqualFold(contentType, "application/json") {
			return v1.V1OperationInput{}, &v1.V1Error{Status: 422, Code: "VALIDATION_FAILED", Message: "request body content type is invalid"}
		}
	}
	in := v1.V1OperationInput{
		OperationID:    operationID,
		Method:         method,
		Path:           r.URL.Path,
		PathParameters: r.GetRouterMap(),
		Query:          r.URL.Query(),
		Headers: map[string]string{
			v1.V1RequestIDHeader:     strings.TrimSpace(r.Header.Get(v1.V1RequestIDHeader)),
			v1.V1VersionHeader:       strings.TrimSpace(r.Header.Get(v1.V1VersionHeader)),
			v1.V1ClientVersionHeader: strings.TrimSpace(r.Header.Get(v1.V1ClientVersionHeader)),
			v1.V1IdempotencyHeader:   strings.TrimSpace(r.Header.Get(v1.V1IdempotencyHeader)),
		},
		Body:           body,
		BodyPresent:    bodyPresent,
		AccessToken:    v1BearerToken(r.Header.Get("Authorization")),
		RequestID:      strings.TrimSpace(r.Header.Get(v1.V1RequestIDHeader)),
		IdempotencyKey: strings.TrimSpace(r.Header.Get(v1.V1IdempotencyHeader)),
	}
	if err = v1.ValidateV1Request(in); err != nil {
		return v1.V1OperationInput{}, err
	}
	return in, nil
}

// SuccessResponse 将 Service 输出转换为接口成功响应。
func SuccessResponse(in v1.V1OperationInput, out *v1.V1OperationOutput) *v1.V1Response {
	return &v1.V1Response{
		V1Version:    v1.V1Version,
		RequestID:    in.RequestID,
		ServerTimeMs: time.Now().UnixMilli(),
		Data:         out.Data,
		ChangeCursor: out.ChangeCursor,
		ETag:         out.ETag,
	}
}

// v1BearerToken 仅接受接口规定的 Bearer credential，避免把其他 scheme 当作访问令牌。
func v1BearerToken(authorization string) string {
	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return ""
	}
	return parts[1]
}
