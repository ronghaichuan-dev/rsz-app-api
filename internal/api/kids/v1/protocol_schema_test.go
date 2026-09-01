package v1

import (
	"encoding/json"
	"testing"
)

// TestOpenAPISpecIncludesDetailedSchemas 验证 Swagger 使用的接口说明包含具体请求和成功响应定义。
func TestOpenAPISpecIncludesDetailedSchemas(t *testing.T) {
	var spec map[string]any
	if err := json.Unmarshal(OpenAPISpec(), &spec); err != nil {
		t.Fatalf("接口说明解析失败: %v", err)
	}
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatal("接口说明缺少 paths")
	}
	path, ok := paths["/v1/assets/uploads:prepare"].(map[string]any)
	if !ok {
		t.Fatal("接口说明缺少资产上传路径")
	}
	operation, ok := path["post"].(map[string]any)
	if !ok || operation["requestBody"] == nil || operation["responses"] == nil {
		t.Fatal("接口说明缺少请求体或响应定义")
	}
}

// TestV1ErrorIncludesField 验证服务端日志会保留校验失败字段路径而不包含请求原文。
func TestV1ErrorIncludesField(t *testing.T) {
	field := "/body/proof_nonce"
	err := (&V1Error{Code: "VALIDATION_FAILED", Field: &field, Message: "string is too short"}).Error()
	if err != "/body/proof_nonce: string is too short" {
		t.Fatalf("校验错误日志不包含字段路径: %s", err)
	}
}

// TestValidateV1Request 验证冻结接口会接受完整合法请求，并拒绝未知字段和非法 query。
func TestValidateV1Request(t *testing.T) {
	body, present, err := DecodeV1Body([]byte(`{
        "upload_id":"upload:v1:00000000-0000-4000-8000-00000000000f",
        "purpose":"task_proof",
        "circle_id":"circle:v1:00000000-0000-4000-8000-000000000001",
        "content_type":"image/jpeg",
        "byte_size":1024,
        "sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    }`))
	if err != nil || !present {
		t.Fatalf("解析合法请求失败: present=%v, err=%v", present, err)
	}
	in := V1OperationInput{
		OperationID: "prepareAssetUpload",
		Method:      "POST",
		Path:        "/v1/assets/uploads:prepare",
		Headers: map[string]string{
			V1RequestIDHeader:     "request:valid-0001",
			V1VersionHeader:       V1Version,
			V1ClientVersionHeader: "1.0.0",
			V1IdempotencyHeader:   "idem:v1:0123456789abcdef",
		},
		PathParameters: map[string]string{},
		Query:          map[string][]string{},
		Body:           body,
		BodyPresent:    true,
		AccessToken:    "fixture-access-token",
	}
	if err = ValidateV1Request(in); err != nil {
		t.Fatalf("合法请求被拒绝: %v", err)
	}
	in.Body["client_role"] = "admin"
	if err = ValidateV1Request(in); err == nil {
		t.Fatal("未知请求字段未被拒绝")
	}
}

// TestDecodeV1BodyRejectsDuplicateKey 验证重复 JSON 键不会被静默覆盖。
func TestDecodeV1BodyRejectsDuplicateKey(t *testing.T) {
	if _, _, err := DecodeV1Body([]byte(`{"purpose":"task_proof","purpose":"reward_image"}`)); err == nil {
		t.Fatal("重复 JSON 键未被拒绝")
	}
}

// TestValidateV1RequestPreservesRepeatedQuery 验证数组 query 保留全部值且单值参数拒绝重复。
func TestValidateV1RequestPreservesRepeatedQuery(t *testing.T) {
	in := V1OperationInput{
		OperationID: "getMemberBalances",
		Method:      "GET",
		Path:        "/v1/circles/circle:v1:00000000-0000-4000-8000-000000000001/member-balances",
		Headers: map[string]string{
			V1RequestIDHeader:     "request:valid-0002",
			V1VersionHeader:       V1Version,
			V1ClientVersionHeader: "1.0.0",
		},
		PathParameters: map[string]string{"circle_id": "circle:v1:00000000-0000-4000-8000-000000000001"},
		Query: map[string][]string{
			"member_id": {
				"member:v1:00000000-0000-4000-8000-000000000001",
				"member:v1:00000000-0000-4000-8000-000000000002",
			},
		},
		Body:        map[string]any{},
		AccessToken: "fixture-access-token",
	}
	if err := ValidateV1Request(in); err != nil {
		t.Fatalf("合法数组 query 被拒绝: %v", err)
	}
	in.OperationID = "listMyCircles"
	in.Method = "GET"
	in.Path = "/v1/circles"
	in.PathParameters = map[string]string{}
	in.Query = map[string][]string{"limit": {"20", "30"}}
	if err := ValidateV1Request(in); err == nil {
		t.Fatal("重复单值 query 未被拒绝")
	}
}

// TestValidateV1TaskOccurrencesLimitBoundary 验证 OpenAPI 请求校验拒绝超过 200 的 occurrence 分页请求。
func TestValidateV1TaskOccurrencesLimitBoundary(t *testing.T) {
	for _, limit := range []string{"1", "100", "200"} {
		in := v1TaskOccurrencesValidationInput(limit)
		if err := ValidateV1Request(in); err != nil {
			t.Fatalf("合法 occurrence limit=%s 被 OpenAPI 校验拒绝: %v", limit, err)
		}
	}
	for _, limit := range []string{"201", "500"} {
		in := v1TaskOccurrencesValidationInput(limit)
		err := ValidateV1Request(in)
		v1Error, ok := err.(*V1Error)
		if !ok || v1Error.Status != 422 || v1Error.Code != "VALIDATION_FAILED" {
			t.Fatalf("非法 occurrence limit=%s 未返回 422/VALIDATION_FAILED: %v", limit, err)
		}
	}
	in := v1TaskOccurrencesValidationInput("1")
	in.Query["member_id"] = []string{"not-a-member-id"}
	if err := ValidateV1Request(in); err == nil {
		t.Fatal("非法 optional member_id 未被 OpenAPI 校验拒绝")
	}
}

// v1TaskOccurrencesValidationInput 构造只用于 OpenAPI 请求校验的 occurrence 请求。
func v1TaskOccurrencesValidationInput(limit string) V1OperationInput {
	circleID := "circle:v1:00000000-0000-4000-8000-000000000001"
	return V1OperationInput{
		OperationID: "listTaskOccurrences",
		Method:      "GET",
		Path:        "/v1/circles/" + circleID + "/task-occurrences",
		Headers: map[string]string{
			V1RequestIDHeader:     "request:occurrence-limit-0001",
			V1VersionHeader:       V1Version,
			V1ClientVersionHeader: "1.0.0",
		},
		PathParameters: map[string]string{"circle_id": circleID},
		Query: map[string][]string{
			"start_date":         {"2026-03-01"},
			"end_date_exclusive": {"2026-04-01"},
			"zone_id":            {"Asia/Shanghai"},
			"limit":              {limit},
			"member_id":          {"member:v1:00000000-0000-4000-8000-000000000001"},
		},
		Body:        map[string]any{},
		AccessToken: "fixture-access-token",
	}
}

// TestValidateV1ResponseData 验证成功数据也按 operation 专属 schema 检查。
func TestValidateV1ResponseData(t *testing.T) {
	data := map[string]any{
		"session": map[string]any{
			"token_type":    "Bearer",
			"access_token":  "0123456789abcdef",
			"refresh_token": "fedcba9876543210",
			"metadata": map[string]any{
				"session_id":            "session:v1:00000000-0000-4000-8000-000000000001",
				"principal_type":        "invite_guest",
				"status":                "active",
				"issued_at_ms":          int64(1),
				"access_expires_at_ms":  int64(2),
				"refresh_expires_at_ms": int64(3),
			},
		},
		"capabilities":        []string{"invite_redeem_administrator", "invite_redeem_member"},
		"guest_expires_at_ms": int64(2),
		"guest_upgrade_grant": "0123456789abcdef01234567",
	}
	if err := ValidateV1ResponseData("createInviteGuestSession", data); err != nil {
		t.Fatalf("合法游客会话响应被拒绝: %v", err)
	}
	data["session"].(map[string]any)["metadata"].(map[string]any)["principal_type"] = "account"
	if err := ValidateV1ResponseData("createInviteGuestSession", data); err == nil {
		t.Fatal("错误 principal_type 的游客会话响应未被拒绝")
	}
}

// TestValidateV1ResponseDataGuestBootstrap 验证游客 bootstrap 只能返回 Guest 分支的字段。
func TestValidateV1ResponseDataGuestBootstrap(t *testing.T) {
	data := map[string]any{
		"principal_kind": "invite_guest",
		"session": map[string]any{
			"session_id":            "session:v1:00000000-0000-4000-8000-000000000001",
			"principal_type":        "invite_guest",
			"status":                "active",
			"issued_at_ms":          int64(1),
			"access_expires_at_ms":  int64(2),
			"refresh_expires_at_ms": int64(3),
		},
		"capabilities":        []string{"invite_redeem_administrator", "invite_redeem_member"},
		"guest_expires_at_ms": int64(2),
		"bootstrap_cursor":    nil,
	}
	if err := ValidateV1ResponseData("getCurrentAccount", data); err != nil {
		t.Fatalf("合法游客 bootstrap 响应被拒绝: %v", err)
	}
	data["account"] = map[string]any{}
	if err := ValidateV1ResponseData("getCurrentAccount", data); err == nil {
		t.Fatal("Guest bootstrap 中的账号字段未被拒绝")
	}
}
