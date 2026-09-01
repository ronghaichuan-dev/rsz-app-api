package kids

import "testing"

// TestV1MemberInviteRedemptionChanges 验证成员兑换邀请码会把完整 Member 事实写入同步提交的 members 变更集。
func TestV1MemberInviteRedemptionChanges(t *testing.T) {
	member := map[string]any{
		"member_id":        "member:v1:00000000-0000-4000-8000-000000000001",
		"circle_id":        "circle:v1:00000000-0000-4000-8000-000000000001",
		"bound_account_id": "account:v1:00000000-0000-4000-8000-000000000001",
		"display_name":     "成员",
		"gender":           "female",
		"avatar":           map[string]any{"kind": "preset", "preset_id": "gif_female_6"},
		"status":           "active",
		"version":          int64(2),
		"created_at_ms":    int64(1),
		"updated_at_ms":    int64(2),
		"deleted_at_ms":    nil,
	}
	changes := v1SyncChanges(v1InviteRedemptionChanges(
		"member",
		map[string]any{"circle_id": member["circle_id"]},
		map[string]any{"membership_id": "membership:v1:00000000-0000-4000-8000-000000000001"},
		member,
	))
	members, ok := changes["members"].([]any)
	if !ok || len(members) != 1 {
		t.Fatalf("成员兑换同步提交缺少 members 变更: %+v", changes["members"])
	}
	persisted, ok := members[0].(map[string]any)
	if !ok || persisted["member_id"] != member["member_id"] || persisted["bound_account_id"] != member["bound_account_id"] || persisted["status"] != "active" || persisted["avatar"] == nil {
		t.Fatalf("成员兑换同步提交未包含完整 active Member 事实: %+v", members[0])
	}
}
