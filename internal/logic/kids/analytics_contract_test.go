package kids

import (
	"testing"
	"time"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

// TestV1StatisticsEmptySeries 验证没有任何统计事实时仍返回完整、可解码的零值日 bucket 序列。
func TestV1StatisticsEmptySeries(t *testing.T) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("加载测试时区失败: %v", err)
	}
	start := time.Date(2024, time.March, 4, 0, 0, 0, 0, location)
	end := time.Date(2024, time.March, 11, 0, 0, 0, 0, location)
	data := v1StatisticsSeriesData(
		"member:v1:00000000-0000-4000-8000-000000000001", "tasks", "weekly", "day", "monday", location.String(),
		start.UTC().UnixMilli(), end.UTC().UnixMilli(), start, end, location, nil, "cur:v1:00000000",
	)
	if err = v1.ValidateV1ResponseData("getStatistics", data); err != nil {
		t.Fatalf("零事实统计响应不符合冻结合同: %v", err)
	}
	buckets := data["buckets"].([]map[string]any)
	if len(buckets) != 7 {
		t.Fatalf("默认周统计应返回 7 个日 bucket，实际=%d", len(buckets))
	}
	for index, bucket := range buckets {
		if bucket["index"] != index || bucket["value"] != int64(0) {
			t.Fatalf("零事实 bucket 不正确 index=%d bucket=%v", index, bucket)
		}
	}
	if buckets[0]["start_at_ms"] != start.UTC().UnixMilli() || buckets[len(buckets)-1]["end_at_ms"] != end.UTC().UnixMilli() {
		t.Fatalf("统计 bucket 没有保持请求 range 边界: buckets=%v", buckets)
	}
	summary := data["summary"].(map[string]any)
	if summary["total"] != int64(0) || summary["peak_value"] != int64(0) || summary["non_zero_bucket_count"] != int64(0) {
		t.Fatalf("零事实统计汇总不正确: %v", summary)
	}
}

// TestV1StatisticsBucketRangeLimit 验证超过协议最大 bucket 数量的确定性请求会在查询前被拒绝。
func TestV1StatisticsBucketRangeLimit(t *testing.T) {
	location := time.UTC
	start := time.Date(2000, time.January, 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 0, 10001)
	if err := v1ValidateStatisticsBucketRange(start, end, location, "day", "monday"); err == nil {
		t.Fatal("超过最大 bucket 数量的统计范围不应继续聚合")
	}
}
