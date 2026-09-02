package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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
	v1DeploymentAssertStarTransactionPage(t, client, fixture, "empty", fixture.zeroBalanceMemberID)
	v1DeploymentAssertStarTransactionPage(t, client, fixture, "existing", fixture.ledgerMemberID)
	if strings.EqualFold(strings.TrimSpace(os.Getenv("KIDS_DEPLOY_SMOKE_ADJUST_ENABLED")), "true") {
		v1DeploymentAdjustAndReverse(t, client, fixture)
	}

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

// v1DeploymentAssertStarTransactionPage 验证空流水和已有流水均返回可被客户端解码的规范分页响应。
func v1DeploymentAssertStarTransactionPage(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, name, memberID string) {
	t.Helper()
	query := url.Values{"member_id": {memberID}, "limit": {"20"}}
	requestID := "smoke:star-transactions:" + name
	request, err := http.NewRequest(http.MethodGet, fixture.baseURL+"/v1/circles/"+url.PathEscape(fixture.circleID)+"/star-transactions?"+query.Encode(), nil)
	if err != nil {
		t.Fatalf("构造星星流水读取请求失败: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+fixture.accessToken)
	request.Header.Set(V1RequestIDHeader, requestID)
	request.Header.Set(V1VersionHeader, V1Version)
	request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("星星流水读取请求失败: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("读取星星流水响应失败: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("星星流水读取应返回 200 case=%s got=%d body=%s", name, response.StatusCode, string(body))
	}
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	if err = json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("解析星星流水成功响应失败: %v", err)
	}
	if err = ValidateV1ResponseData("listStarTransactions", envelope.Data); err != nil {
		t.Fatalf("星星流水成功响应不符合合同 case=%s err=%v body=%s", name, err, string(body))
	}
	items := envelope.Data["items"].([]any)
	if name == "empty" && len(items) != 0 {
		t.Fatalf("零余额成员应返回空流水页，实际 items=%v", items)
	}
	if name == "existing" && len(items) == 0 {
		t.Fatal("已有流水成员不应返回空流水页")
	}
	for _, rawItem := range items {
		item := rawItem.(map[string]any)
		if item["circle_id"] != fixture.circleID || item["member_id"] != memberID {
			t.Fatalf("星星流水越过请求范围 case=%s item=%v", name, item)
		}
	}
	if envelope.Data["has_more"].(bool) != (envelope.Data["next_cursor"] != nil) {
		t.Fatalf("星星流水分页字段不一致 case=%s data=%v", name, envelope.Data)
	}
}

// v1DeploymentAdjustAndReverse 使用受控零余额成员验证 adjustment 200、幂等重放、账本读取和 append-only 反向调整。
func v1DeploymentAdjustAndReverse(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture) {
	t.Helper()
	before := v1DeploymentReadBalance(t, client, fixture, "before-adjust")
	first := v1DeploymentAdjust(t, client, fixture, "increase", before.memberID, before.version, 1)
	replay := v1DeploymentAdjustReplay(t, client, fixture, first)
	if replay.ledgerID != first.ledgerID || replay.receiptID != first.receiptID || replay.balanceVersion != first.balanceVersion {
		t.Fatal("同一 adjustment 幂等重放没有返回首次 canonical bundle")
	}
	if first.balance != before.balance+1 || first.balanceVersion != before.version+1 {
		t.Fatal("正向 adjustment 没有按 canonical version 更新余额")
	}
	v1DeploymentAssertAdjustmentLedger(t, client, fixture, first)
	v1DeploymentAssertAdjustmentSyncLedger(t, client, fixture, first)
	reversal := v1DeploymentAdjust(t, client, fixture, "reverse", first.memberID, first.balanceVersion, -1)
	if reversal.balance != before.balance || reversal.balanceVersion != first.balanceVersion+1 {
		t.Fatal("反向 adjustment 没有以 append-only 方式恢复测试余额")
	}
	v1DeploymentAssertAdjustmentLedger(t, client, fixture, reversal)
	v1DeploymentAssertAdjustmentSyncLedger(t, client, fixture, reversal)
	after := v1DeploymentReadBalance(t, client, fixture, "after-adjust")
	if after.balance != before.balance || after.version != reversal.balanceVersion {
		t.Fatal("重启前读取的余额投影与 adjustment 结果不一致")
	}
}

// v1DeploymentBalance 保存读取到的单成员 canonical 余额版本。
type v1DeploymentBalance struct {
	memberID string
	balance  int64
	version  int64
}

// v1DeploymentAdjustment 保存服务端返回的 adjustment canonical bundle 关键关联字段。
type v1DeploymentAdjustment struct {
	adjustmentID   string
	memberID       string
	ledgerID       string
	ledger         map[string]any
	commitID       string
	commitSequence int64
	receiptID      string
	balance        int64
	balanceVersion int64
	requestBody    []byte
	idempotencyKey string
}

// v1DeploymentReadBalance 读取受控成员的单条余额并提取后续 adjustment 所需 version。
func v1DeploymentReadBalance(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, name string) v1DeploymentBalance {
	t.Helper()
	requestID := "smoke:adjustment:" + name
	response, body := v1DeploymentMemberBalanceRequest(t, client, fixture.baseURL, fixture.circleID, fixture.accessToken, requestID, []string{fixture.zeroBalanceMemberID})
	if response.StatusCode != http.StatusOK {
		t.Fatalf("读取 adjustment 前余额失败 status=%d body=%s", response.StatusCode, string(body))
	}
	var envelope struct {
		Data struct {
			Items []struct {
				MemberID string `json:"member_id"`
				Balance  int64  `json:"balance"`
				Version  int64  `json:"version"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil || len(envelope.Data.Items) != 1 {
		t.Fatalf("解析 adjustment 前余额失败: err=%v body=%s", err, string(body))
	}
	item := envelope.Data.Items[0]
	return v1DeploymentBalance{memberID: item.MemberID, balance: item.Balance, version: item.Version}
}

// v1DeploymentAdjust 发起一次受权 adjustment，并校验完整 200 canonical bundle。
func v1DeploymentAdjust(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, name, memberID string, version, delta int64) v1DeploymentAdjustment {
	t.Helper()
	adjustmentID := "adjustment:v1:" + uuid.NewString()
	requestBody, err := json.Marshal(map[string]any{"adjustment_id": adjustmentID, "member_id": memberID, "delta": delta, "reason": "deployment smoke adjustment", "expected_balance_version": version})
	if err != nil {
		t.Fatalf("序列化 adjustment 请求失败: %v", err)
	}
	idempotencyKey := "idem:v1:" + uuid.NewString()
	return v1DeploymentPostAdjustment(t, client, fixture, "smoke:adjustment:"+name, requestBody, idempotencyKey, adjustmentID)
}

// v1DeploymentAdjustReplay 使用完全相同的请求体和幂等键重放首次 adjustment。
func v1DeploymentAdjustReplay(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, first v1DeploymentAdjustment) v1DeploymentAdjustment {
	t.Helper()
	return v1DeploymentPostAdjustment(t, client, fixture, "smoke:adjustment:replay", first.requestBody, first.idempotencyKey, first.adjustmentID)
}

// v1DeploymentPostAdjustment 执行 adjustment HTTP 请求并校验错误不会被伪装为 200。
func v1DeploymentPostAdjustment(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, requestID string, requestBody []byte, idempotencyKey, adjustmentID string) v1DeploymentAdjustment {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, fixture.baseURL+"/v1/circles/"+url.PathEscape(fixture.circleID)+"/star-adjustments", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("构造 adjustment 请求失败: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+fixture.accessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(V1RequestIDHeader, requestID)
	request.Header.Set(V1VersionHeader, V1Version)
	request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
	request.Header.Set(V1IdempotencyHeader, idempotencyKey)
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("adjustment 请求失败: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("读取 adjustment 响应失败: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("adjustment 应返回 200，实际 status=%d body=%s", response.StatusCode, string(body))
	}
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	if err = json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("解析 adjustment 成功响应失败: %v", err)
	}
	if err = ValidateV1ResponseData("adjustMemberStars", envelope.Data); err != nil {
		t.Fatalf("adjustment 成功 bundle 不符合合同: %v body=%s", err, string(body))
	}
	ledger := envelope.Data["ledger_entry"].(map[string]any)
	balance := envelope.Data["balance"].(map[string]any)
	receipt := envelope.Data["receipt"].(map[string]any)
	return v1DeploymentAdjustment{adjustmentID: adjustmentID, memberID: balance["member_id"].(string), ledgerID: ledger["ledger_id"].(string), ledger: ledger, commitID: receipt["commit_id"].(string), commitSequence: int64(receipt["commit_sequence"].(float64)), receiptID: receipt["receipt_id"].(string), balance: int64(balance["balance"].(float64)), balanceVersion: int64(balance["version"].(float64)), requestBody: requestBody, idempotencyKey: idempotencyKey}
}

// v1DeploymentAssertAdjustmentLedger 验证调整后读取能够看到对应的不可变流水，并保持与变更响应一致。
func v1DeploymentAssertAdjustmentLedger(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, adjustment v1DeploymentAdjustment) {
	t.Helper()
	query := url.Values{"member_id": {fixture.zeroBalanceMemberID}, "source_type": {"adjustment"}, "limit": {"20"}}
	request, err := http.NewRequest(http.MethodGet, fixture.baseURL+"/v1/circles/"+url.PathEscape(fixture.circleID)+"/star-transactions?"+query.Encode(), nil)
	if err != nil {
		t.Fatalf("构造 adjustment 流水读取请求失败: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+fixture.accessToken)
	request.Header.Set(V1RequestIDHeader, "smoke:adjustment:ledger")
	request.Header.Set(V1VersionHeader, V1Version)
	request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("adjustment 流水读取请求失败: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("adjustment 流水读取失败 status=%d err=%v body=%s", response.StatusCode, err, string(body))
	}
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	if err = json.Unmarshal(body, &envelope); err != nil || ValidateV1ResponseData("listStarTransactions", envelope.Data) != nil {
		t.Fatalf("adjustment 流水响应不符合合同: err=%v body=%s", err, string(body))
	}
	for _, rawItem := range envelope.Data["items"].([]any) {
		item := rawItem.(map[string]any)
		source := item["source"].(map[string]any)
		if source["source_id"] == adjustment.adjustmentID {
			v1DeploymentAssertCanonicalLedger(t, adjustment.ledger, item, "listStarTransactions")
			return
		}
	}
	t.Fatalf("流水读取未包含 adjustment_id=%s", adjustment.adjustmentID)
}

// v1DeploymentAssertAdjustmentSyncLedger 验证调整对应的同步提交保留与变更响应相同的不可变流水投影。
func v1DeploymentAssertAdjustmentSyncLedger(t *testing.T, client *http.Client, fixture v1MemberBalanceDeploymentFixture, adjustment v1DeploymentAdjustment) {
	t.Helper()
	query := url.Values{"change_cursor": {fmt.Sprintf("cur:v1:%08x", adjustment.commitSequence-1)}, "limit": {"20"}}
	request, err := http.NewRequest(http.MethodGet, fixture.baseURL+"/v1/circles/"+url.PathEscape(fixture.circleID)+"/sync?"+query.Encode(), nil)
	if err != nil {
		t.Fatalf("构造 adjustment 同步读取请求失败: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+fixture.accessToken)
	request.Header.Set(V1RequestIDHeader, "smoke:adjustment:sync")
	request.Header.Set(V1VersionHeader, V1Version)
	request.Header.Set(V1ClientVersionHeader, "deployment-smoke-v1")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("adjustment 同步读取请求失败: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("adjustment 同步读取失败 status=%d err=%v body=%s", response.StatusCode, err, string(body))
	}
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	if err = json.Unmarshal(body, &envelope); err != nil || ValidateV1ResponseData("pullCircleBootstrapDelta", envelope.Data) != nil {
		t.Fatalf("adjustment 同步响应不符合合同: err=%v body=%s", err, string(body))
	}
	for _, rawCommit := range envelope.Data["commits"].([]any) {
		commit := rawCommit.(map[string]any)
		if commit["commit_id"] != adjustment.commitID {
			continue
		}
		for _, rawLedger := range commit["changes"].(map[string]any)["ledger_entries"].([]any) {
			ledger := rawLedger.(map[string]any)
			if ledger["ledger_id"] == adjustment.ledgerID {
				v1DeploymentAssertCanonicalLedger(t, adjustment.ledger, ledger, "pullCircleBootstrapDelta")
				return
			}
		}
		t.Fatalf("同步提交未包含 ledger_id=%s", adjustment.ledgerID)
	}
	t.Fatalf("同步响应未包含 commit_id=%s", adjustment.commitID)
}

// v1DeploymentAssertCanonicalLedger 验证同一流水标识在不同接口中逐字段保持同一不可变事实。
func v1DeploymentAssertCanonicalLedger(t *testing.T, expected, actual map[string]any, source string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("同一 ledger_id 的规范投影不一致 source=%s expected=%v actual=%v", source, expected, actual)
	}
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

// TestV1DeploymentValidationSmoke 对已配置的部署环境执行 47 个 operation 的真实路由、校验和可观测性验收。
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
