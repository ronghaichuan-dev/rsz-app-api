package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// TestV1DeploymentValidationSmoke 对已配置的部署环境执行 46 个 operation 的真实路由、校验和可观测性验收。
func TestV1DeploymentValidationSmoke(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_BASE_URL")), "/")
	if baseURL == "" {
		t.Skip("未配置 KIDS_DEPLOY_SMOKE_BASE_URL，跳过部署级 smoke")
	}
	spec, err := loadV1Spec()
	if err != nil {
		t.Fatalf("加载 v1 OpenAPI 失败: %v", err)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	for path, rawPath := range spec["paths"].(map[string]any) {
		for method, rawOperation := range rawPath.(map[string]any) {
			if method == "parameters" {
				continue
			}
			operation := rawOperation.(map[string]any)
			operationID := operation["operationId"].(string)
			t.Run(operationID, func(t *testing.T) {
				requestID := "smoke:" + operationID + ":0001"
				request, requestErr := http.NewRequest(strings.ToUpper(method), baseURL+v1SmokePath(path), bytes.NewReader([]byte("{}")))
				if requestErr != nil {
					t.Fatalf("构造请求失败: %v", requestErr)
				}
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set(V1RequestIDHeader, requestID)
				request.Header.Set(V1VersionHeader, V1Version)
				request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
				request.Header.Set(V1IdempotencyHeader, "idem:v1:00000000-0000-4000-8000-000000000001")
				response, requestErr := client.Do(request)
				if requestErr != nil {
					t.Fatalf("部署请求失败: %v", requestErr)
				}
				defer response.Body.Close()
				body, requestErr := io.ReadAll(io.LimitReader(response.Body, 1<<20))
				if requestErr != nil {
					t.Fatalf("读取响应失败: %v", requestErr)
				}
				if response.StatusCode < 400 || response.StatusCode >= 500 {
					t.Fatalf("部署校验请求应返回受控 4xx，operation=%s status=%d", operationID, response.StatusCode)
				}
				v1SmokeAssertErrorEnvelope(t, operationID, requestID, response, body)
			})
		}
	}
}

// TestV1DeploymentTaskOccurrencesRejectsOversizedLimits 在真实部署环境固定 occurrence 的 1..200 分页边界。
func TestV1DeploymentTaskOccurrencesRejectsOversizedLimits(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_BASE_URL")), "/")
	if baseURL == "" {
		t.Skip("未配置 KIDS_DEPLOY_SMOKE_BASE_URL，跳过部署级 occurrence 负向 smoke")
	}
	client := &http.Client{Timeout: 15 * time.Second}
	circleID := "circle:v1:00000000-0000-4000-8000-000000000001"
	for _, limit := range []string{"201", "500"} {
		t.Run("limit_"+limit, func(t *testing.T) {
			requestID := "smoke:task-occurrences:" + limit
			request, err := http.NewRequest(http.MethodGet, baseURL+"/v1/circles/"+circleID+"/task-occurrences?start_date=2026-03-01&end_date_exclusive=2026-04-01&zone_id=Asia%2FShanghai&limit="+limit, nil)
			if err != nil {
				t.Fatalf("构造 occurrence 负向请求失败: %v", err)
			}
			request.Header.Set("Authorization", "Bearer deployment-smoke-invalid-token")
			request.Header.Set(V1RequestIDHeader, requestID)
			request.Header.Set(V1VersionHeader, V1Version)
			request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
			response, err := client.Do(request)
			if err != nil {
				t.Fatalf("occurrence 负向请求失败: %v", err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
			if err != nil {
				t.Fatalf("读取 occurrence 负向响应失败: %v", err)
			}
			if response.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("limit=%s 应返回 422，实际 status=%d body=%s", limit, response.StatusCode, string(body))
			}
			v1SmokeAssertErrorEnvelope(t, "listTaskOccurrences", requestID, response, body)
			if v1SmokeErrorCode(body) != "VALIDATION_FAILED" {
				t.Fatalf("limit=%s 应返回 VALIDATION_FAILED，实际 body=%s", limit, string(body))
			}
		})
	}
}

// v1SmokePath 将 OpenAPI path 参数替换为符合冻结 ID 规则的合成值，不使用真实用户数据。
func v1SmokePath(path string) string {
	replacements := map[string]string{
		"{circle_id}":        "circle:v1:00000000-0000-4000-8000-000000000001",
		"{session_id}":       "session:v1:00000000-0000-4000-8000-000000000001",
		"{administrator_id}": "admin:v1:00000000-0000-4000-8000-000000000001",
		"{member_id}":        "member:v1:00000000-0000-4000-8000-000000000001",
		"{reward_id}":        "reward:v1:00000000-0000-4000-8000-000000000001",
		"{task_id}":          "task:v1:00000000-0000-4000-8000-000000000001",
		"{task_tag_id}":      "task-tag:v1:00000000-0000-4000-8000-000000000001",
		"{completion_id}":    "completion:v1:00000000-0000-4000-8000-000000000001",
		"{invite_id}":        "invite:v1:00000000-0000-4000-8000-000000000001",
		"{upload_id}":        "upload:v1:00000000-0000-4000-8000-000000000001",
	}
	for placeholder, value := range replacements {
		path = strings.ReplaceAll(path, placeholder, value)
	}
	return path
}

// v1SmokeAssertErrorEnvelope 校验部署失败响应仍符合 ErrorEnvelope 与 request/trace 可观测性约定。
func v1SmokeAssertErrorEnvelope(t *testing.T, operationID, requestID string, response *http.Response, body []byte) {
	t.Helper()
	var envelope map[string]any
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("响应不是 JSON ErrorEnvelope，operation=%s: %v", operationID, err)
	}
	if envelope["contract_version"] != V1Version || envelope["request_id"] != requestID {
		t.Fatalf("响应关联字段错误，operation=%s body=%s", operationID, string(body))
	}
	errorBody, ok := envelope["error"].(map[string]any)
	if !ok || errorBody["code"] == "" {
		t.Fatalf("响应缺少受控错误码，operation=%s", operationID)
	}
	if _, ok = errorBody["retryable"].(bool); !ok {
		t.Fatalf("响应缺少 retryable，operation=%s", operationID)
	}
	traceID, ok := errorBody["trace_id"].(string)
	if !ok || len(traceID) < 8 || response.Header.Get("X-Trace-Id") != traceID {
		t.Fatalf("响应缺少可关联 traceId，operation=%s", operationID)
	}
}

// v1SmokeErrorCode 从受控 ErrorEnvelope 提取稳定错误码，供部署级断言使用。
func v1SmokeErrorCode(body []byte) string {
	var envelope struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ""
	}
	return envelope.Error.Code
}
