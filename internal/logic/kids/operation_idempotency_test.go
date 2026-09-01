package kids

import (
	"testing"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1OperationSupportsIdempotency 验证幂等能力只对登记过的写操作生效。
func TestV1OperationSupportsIdempotency(t *testing.T) {
	testCases := []struct {
		name        string
		operationID string
		method      string
		want        bool
	}{
		{name: "已登记兑换写操作", operationID: "redeemReward", method: "POST", want: true},
		{name: "已登记任务写操作", operationID: "completeTask", method: "POST", want: true},
		{name: "读取操作", operationID: "redeemReward", method: "GET", want: false},
		{name: "未登记写操作", operationID: "unknownMutation", method: "POST", want: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			in := v1.V1OperationInput{OperationID: testCase.operationID, Method: testCase.method}
			if got := v1OperationSupportsIdempotency(in); got != testCase.want {
				t.Fatalf("幂等操作判断不正确 operation_id=%s method=%s got=%t want=%t", testCase.operationID, testCase.method, got, testCase.want)
			}
		})
	}
}
