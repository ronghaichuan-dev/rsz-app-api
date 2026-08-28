package kids

import "testing"

// TestV1ValidateTaskOccurrencesQuery 验证 occurrence 分页严格保持 1..200 的合同边界和日历语义。
func TestV1ValidateTaskOccurrencesQuery(t *testing.T) {
	for _, limit := range []string{"1", "100", "200"} {
		if _, err := v1ValidateTaskOccurrencesQuery(v1TaskOccurrenceQuery(limit)); err != nil {
			t.Fatalf("合法 limit=%s 被拒绝: %v", limit, err)
		}
	}
	for _, limit := range []string{"0", "201", "500"} {
		if _, err := v1ValidateTaskOccurrencesQuery(v1TaskOccurrenceQuery(limit)); err == nil {
			t.Fatalf("非法 limit=%s 未被拒绝", limit)
		}
	}
	invalidQueries := []map[string][]string{
		{"start_date": {"2026-02-30"}, "end_date_exclusive": {"2026-03-01"}, "zone_id": {"Asia/Shanghai"}, "limit": {"1"}},
		{"start_date": {"2026-03-01"}, "end_date_exclusive": {"2026-03-01"}, "zone_id": {"Asia/Shanghai"}, "limit": {"1"}},
		{"start_date": {"2026-03-01"}, "end_date_exclusive": {"2026-03-02"}, "zone_id": {"Asia/NotAZone"}, "limit": {"1"}},
	}
	for _, query := range invalidQueries {
		if _, err := v1ValidateTaskOccurrencesQuery(query); err == nil {
			t.Fatalf("非法 occurrence query 未被拒绝: %#v", query)
		}
	}
}

// v1TaskOccurrenceQuery 构造满足接口必填 query 的最小合法基线。
func v1TaskOccurrenceQuery(limit string) map[string][]string {
	return map[string][]string{
		"start_date":         {"2026-03-01"},
		"end_date_exclusive": {"2026-04-01"},
		"zone_id":            {"Asia/Shanghai"},
		"limit":              {limit},
		"member_id":          {"member:v1:00000000-0000-4000-8000-000000000001"},
	}
}
