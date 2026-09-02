package v1

import "testing"

// TestV1ContractSmokeCoverage 确保 CI 对冻结的全部 47 个 operation 保持显式协议覆盖入口。
func TestV1ContractSmokeCoverage(t *testing.T) {
	spec, err := loadV1Spec()
	if err != nil {
		t.Fatalf("加载 v1 OpenAPI 失败: %v", err)
	}
	operations := make(map[string]struct{})
	for _, rawPath := range spec["paths"].(map[string]any) {
		path := rawPath.(map[string]any)
		for method, rawOperation := range path {
			if method == "parameters" {
				continue
			}
			operation, ok := rawOperation.(map[string]any)
			if !ok {
				t.Fatalf("OpenAPI operation 不是对象: method=%s", method)
			}
			operationID, _ := operation["operationId"].(string)
			if operationID == "" {
				t.Fatalf("OpenAPI operation 缺少 operationId: method=%s", method)
			}
			if _, duplicated := operations[operationID]; duplicated {
				t.Fatalf("OpenAPI operationId 重复: %s", operationID)
			}
			operations[operationID] = struct{}{}
			if _, _, err = findV1OperationByID(spec, operationID); err != nil {
				t.Fatalf("operation 无法被请求校验器定位: %s: %v", operationID, err)
			}
		}
	}
	if len(operations) != 47 {
		t.Fatalf("v1 smoke 覆盖的 operation 数量错误: got=%d want=47", len(operations))
	}
}
