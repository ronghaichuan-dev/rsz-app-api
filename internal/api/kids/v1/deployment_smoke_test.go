package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// TestV1DeploymentMemberBalancesContract 使用隔离账号验证余额快照、权限、故障 envelope 与故障恢复后的首个成功读取。
func TestV1DeploymentMemberBalancesContract(t *testing.T) {
	fixture, ok := loadV1MemberBalanceDeploymentFixture(t)
	if !ok {
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}

	assertMemberBalances := func(name string, token string, memberIDs []string, expectedStatus int) {
		t.Helper()
		requestID := "smoke:member-balances:" + name
		response, body := v1DeploymentMemberBalanceRequest(t, client, fixture.baseURL, fixture.circleID, token, requestID, memberIDs)
		if response.StatusCode != expectedStatus {
			t.Fatalf("余额读取状态不正确 case=%s got=%d want=%d body=%s", name, response.StatusCode, expectedStatus, string(body))
		}
		if expectedStatus != http.StatusOK {
			v1SmokeAssertErrorEnvelope(t, "getMemberBalances", requestID, response, body)
			return
		}
		v1DeploymentAssertMemberBalanceSuccess(t, requestID, response, body, fixture.circleID, memberIDs)
	}

	assertMemberBalances("single", fixture.accessToken, fixture.memberIDs[:1], http.StatusOK)
	assertMemberBalances("batch", fixture.accessToken, fixture.memberIDs, http.StatusOK)
	assertMemberBalances("zero", fixture.accessToken, []string{fixture.zeroBalanceMemberID}, http.StatusOK)
	assertMemberBalances("ledger", fixture.accessToken, []string{fixture.ledgerMemberID}, http.StatusOK)
	assertMemberBalances("forbidden", fixture.forbiddenAccessToken, fixture.memberIDs[:1], http.StatusForbidden)
	assertMemberBalances("missing", fixture.accessToken, []string{"member:v1:00000000-0000-4000-8000-000000000404"}, http.StatusNotFound)
	assertMemberBalances("invalid", fixture.accessToken, []string{"invalid-member-id"}, http.StatusUnprocessableEntity)

	if fixture.unavailableURL == "" {
		return
	}
	requestID := "smoke:member-balances:unavailable"
	request, err := http.NewRequest(http.MethodGet, fixture.unavailableURL, nil)
	if err != nil {
		t.Fatalf("构造余额依赖故障演练请求失败: %v", err)
	}
	request.Header.Set(V1RequestIDHeader, requestID)
	request.Header.Set(V1VersionHeader, V1Version)
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("余额依赖故障演练请求失败: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("读取余额依赖故障演练响应失败: %v", err)
	}
	if response.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("余额依赖故障演练应返回 503，实际 status=%d body=%s", response.StatusCode, string(body))
	}
	v1SmokeAssertErrorEnvelope(t, "getMemberBalances", requestID, response, body)
	if v1SmokeErrorCode(body) != "UNAVAILABLE" {
		t.Fatalf("余额依赖故障演练应返回 UNAVAILABLE，body=%s", string(body))
	}
	assertMemberBalances("recovered", fixture.accessToken, []string{fixture.zeroBalanceMemberID}, http.StatusOK)
}

// v1MemberBalanceDeploymentFixture 保存部署回归所需的隔离账号、圈子和可控成员信息。
type v1MemberBalanceDeploymentFixture struct {
	baseURL              string
	accessToken          string
	forbiddenAccessToken string
	circleID             string
	memberIDs            []string
	zeroBalanceMemberID  string
	ledgerMemberID       string
	unavailableURL       string
}

// loadV1MemberBalanceDeploymentFixture 读取完整的部署回归环境；未配置时跳过，避免误用真实账号。
func loadV1MemberBalanceDeploymentFixture(t *testing.T) (v1MemberBalanceDeploymentFixture, bool) {
	t.Helper()
	fixture := v1MemberBalanceDeploymentFixture{
		baseURL:              strings.TrimRight(strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_BASE_URL")), "/"),
		accessToken:          strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_ACCESS_TOKEN")),
		forbiddenAccessToken: strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_FORBIDDEN_ACCESS_TOKEN")),
		circleID:             strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_CIRCLE_ID")),
		memberIDs:            v1DeploymentMemberIDs(os.Getenv("KIDS_DEPLOY_SMOKE_MEMBER_IDS")),
		zeroBalanceMemberID:  strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_ZERO_BALANCE_MEMBER_ID")),
		ledgerMemberID:       strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_LEDGER_MEMBER_ID")),
		unavailableURL:       strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_UNAVAILABLE_URL")),
	}
	if fixture.baseURL == "" && fixture.accessToken == "" && fixture.forbiddenAccessToken == "" && fixture.circleID == "" && len(fixture.memberIDs) == 0 && fixture.zeroBalanceMemberID == "" && fixture.ledgerMemberID == "" {
		t.Skip("未配置隔离余额部署回归资产，跳过真实余额读取验收")
		return v1MemberBalanceDeploymentFixture{}, false
	}
	if fixture.baseURL == "" || fixture.accessToken == "" || fixture.forbiddenAccessToken == "" || fixture.circleID == "" || len(fixture.memberIDs) != 2 || fixture.zeroBalanceMemberID == "" || fixture.ledgerMemberID == "" {
		t.Fatal("余额部署回归配置不完整；请配置隔离账号、圈子、两个成员、零余额成员和已有流水成员")
	}
	return fixture, true
}

// v1DeploymentMemberIDs 解析逗号分隔的两个受控成员标识。
func v1DeploymentMemberIDs(value string) []string {
	var memberIDs []string
	for _, memberID := range strings.Split(value, ",") {
		if memberID = strings.TrimSpace(memberID); memberID != "" {
			memberIDs = append(memberIDs, memberID)
		}
	}
	return memberIDs
}

// v1DeploymentMemberBalanceRequest 向部署环境发送一次真实的 canonical 余额读取请求。
func v1DeploymentMemberBalanceRequest(t *testing.T, client *http.Client, baseURL, circleID, token, requestID string, memberIDs []string) (*http.Response, []byte) {
	t.Helper()
	query := make(url.Values, len(memberIDs))
	for _, memberID := range memberIDs {
		query.Add("member_id", memberID)
	}
	request, err := http.NewRequest(http.MethodGet, baseURL+"/v1/circles/"+url.PathEscape(circleID)+"/member-balances?"+query.Encode(), nil)
	if err != nil {
		t.Fatalf("构造余额读取请求失败: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set(V1RequestIDHeader, requestID)
	request.Header.Set(V1VersionHeader, V1Version)
	request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("余额读取请求失败: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("读取余额响应失败: %v", err)
	}
	return response, body
}

// v1DeploymentAssertMemberBalanceSuccess 校验 200 schema、trace 关联、请求集合一一对应以及可用于调整的版本。
func v1DeploymentAssertMemberBalanceSuccess(t *testing.T, requestID string, response *http.Response, body []byte, circleID string, memberIDs []string) {
	t.Helper()
	var envelope struct {
		ContractVersion string         `json:"contract_version"`
		RequestID       string         `json:"request_id"`
		Data            map[string]any `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("余额成功响应不是 JSON: %v", err)
	}
	if envelope.ContractVersion != V1Version || envelope.RequestID != requestID || response.Header.Get("X-Request-Id") != requestID || response.Header.Get("X-Trace-Id") == "" {
		t.Fatalf("余额成功响应缺少关联字段: body=%s", string(body))
	}
	if err := ValidateV1ResponseData("getMemberBalances", envelope.Data); err != nil {
		t.Fatalf("余额成功响应不符合合同: %v body=%s", err, string(body))
	}
	rawItems, ok := envelope.Data["items"].([]any)
	if !ok || len(rawItems) != len(memberIDs) {
		t.Fatalf("余额 items 数量不正确: body=%s", string(body))
	}
	expected := make(map[string]struct{}, len(memberIDs))
	for _, memberID := range memberIDs {
		expected[memberID] = struct{}{}
	}
	seen := make(map[string]struct{}, len(rawItems))
	for _, rawItem := range rawItems {
		item, ok := rawItem.(map[string]any)
		if !ok || item["circle_id"] != circleID {
			t.Fatalf("余额 item 不是请求圈子的成员: body=%s", string(body))
		}
		memberID, _ := item["member_id"].(string)
		if _, ok = expected[memberID]; !ok {
			t.Fatalf("余额 item 不属于请求集合: member_id=%s", memberID)
		}
		if _, duplicated := seen[memberID]; duplicated {
			t.Fatalf("余额 item 出现重复成员: member_id=%s", memberID)
		}
		seen[memberID] = struct{}{}
		version, ok := item["version"].(float64)
		if !ok || version < 1 || version != float64(int64(version)) {
			t.Fatalf("余额 item 缺少可用于调整的 canonical version: member_id=%s", memberID)
		}
	}
}

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
