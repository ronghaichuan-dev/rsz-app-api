package v1

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	// V1Version 是 Clearwave 接口主版本。
	V1Version = "1"
	// V1RequestIDHeader 是接口请求 ID 请求头。
	V1RequestIDHeader = "X-Request-Id"
	// V1VersionHeader 是接口版本请求头。
	V1VersionHeader = "X-Contract-Version"
	// V1ClientVersionHeader 是客户端版本请求头。
	V1ClientVersionHeader = "X-Client-Version"
	// V1IdempotencyHeader 是接口幂等键请求头。
	V1IdempotencyHeader = "Idempotency-Key"
)

// V1OperationInput 是 Controller 按 OpenAPI 严格解析后传给单个接口 Service 方法的输入。
type V1OperationInput struct {
	OperationID    string
	Method         string
	Path           string
	PathParameters map[string]string
	Query          map[string][]string
	Headers        map[string]string
	Body           map[string]any
	BodyPresent    bool
	AccessToken    string
	PrincipalKind  string
	PrincipalID    string
	SessionID      string
	RequestID      string
	IdempotencyKey string
}

// V1OperationOutput 是接口业务返回的数据和 HTTP 状态。
type V1OperationOutput struct {
	Data         map[string]any
	Status       int
	ChangeCursor string
	ETag         string
}

// V1Response 是接口成功 envelope，由统一响应 middleware 直接输出。
type V1Response struct {
	V1Version    string `json:"contract_version"`
	RequestID    string `json:"request_id"`
	ServerTimeMs int64  `json:"server_time_ms"`
	Data         any    `json:"data"`
	ChangeCursor string `json:"change_cursor,omitempty"`
	ETag         string `json:"etag,omitempty"`
}

// V1Error 是接口稳定错误。
type V1Error struct {
	Status       int
	Code         string
	Retryable    bool
	RetryAfterMs *int64
	Field        *string
	Version      *int64
	Message      string
}

// Error 实现 error 接口。
func (e *V1Error) Error() string {
	message := e.Code
	if e.Message != "" {
		message = e.Message
	}
	if e.Field != nil && *e.Field != "" {
		return fmt.Sprintf("%s: %s", *e.Field, message)
	}
	return message
}

// V1Req 是接口请求 DTO 的共同空基类。
type V1Req struct{}

// 以下请求 DTO 声明 GoFrame 冻结路由；Controller 以嵌入的 OpenAPI 严格解析请求字段，
// 再将 path、query、header 与 JSON body 转换为 V1OperationInput。
type GetCurrentAccountReq struct {
	g.Meta `path:"/account/bootstrap" method:"get" tags:"Auth" summary:"获取当前账号"`
}
type SelectCurrentCircleReq struct {
	g.Meta `path:"/account/circle-selection" method:"put" tags:"Circle" summary:"选择当前圈子"`
}
type CommitAssetUploadReq struct {
	g.Meta `path:"/assets/uploads/{upload_id}:commit" method:"post" tags:"Asset" summary:"提交资产上传"`
}
type PrepareAssetUploadReq struct {
	g.Meta `path:"/assets/uploads:prepare" method:"post" tags:"Asset" summary:"准备资产上传"`
}
type ExchangeGoogleProofReq struct {
	g.Meta `path:"/auth/google:exchange" method:"post" tags:"Auth" summary:"交换 Google proof"`
}
type CreateInviteGuestSessionReq struct {
	g.Meta `path:"/auth/guest-sessions" method:"post" tags:"Auth" summary:"创建邀请码游客会话"`
}
type RevokeSessionReq struct {
	g.Meta `path:"/auth/sessions/{session_id}:revoke" method:"post" tags:"Auth" summary:"撤销会话"`
}
type RefreshSessionReq struct {
	g.Meta `path:"/auth/sessions:refresh" method:"post" tags:"Auth" summary:"刷新会话"`
}
type ListMyCirclesReq struct {
	g.Meta `path:"/circles" method:"get" tags:"Circle" summary:"列出我的圈子"`
}
type DeleteCircleReq struct {
	g.Meta `path:"/circles/{circle_id}" method:"delete" tags:"Circle" summary:"删除圈子"`
}
type UpdateCircleReq struct {
	g.Meta `path:"/circles/{circle_id}" method:"patch" tags:"Circle" summary:"更新圈子"`
}
type DeleteAdministratorReq struct {
	g.Meta `path:"/circles/{circle_id}/administrators/{administrator_id}" method:"delete" tags:"Circle" summary:"删除管理员"`
}
type UpsertAdministratorReq struct {
	g.Meta `path:"/circles/{circle_id}/administrators/{administrator_id}" method:"put" tags:"Circle" summary:"更新管理员"`
}
type GetCircleBootstrapReq struct {
	g.Meta `path:"/circles/{circle_id}/bootstrap" method:"get" tags:"Circle" summary:"获取圈子 bootstrap"`
}
type ListExchangeHistoryReq struct {
	g.Meta `path:"/circles/{circle_id}/exchanges" method:"get" tags:"ReadModels" summary:"列出兑换历史"`
}
type CreateInviteReq struct {
	g.Meta `path:"/circles/{circle_id}/invites" method:"post" tags:"Invite" summary:"创建邀请"`
}
type RefreshInviteReq struct {
	g.Meta `path:"/circles/{circle_id}/invites/{invite_id}:refresh" method:"post" tags:"Invite" summary:"刷新邀请"`
}
type RevokeInviteReq struct {
	g.Meta `path:"/circles/{circle_id}/invites/{invite_id}:revoke" method:"post" tags:"Invite" summary:"撤销邀请"`
}
type GetMemberBalancesReq struct {
	g.Meta `path:"/circles/{circle_id}/member-balances" method:"get" tags:"Ledger" summary:"获取成员余额"`
}
type CreateCircleMemberReq struct {
	g.Meta `path:"/circles/{circle_id}/members" method:"post" tags:"Circle" summary:"创建圈子成员"`
}
type DeleteMemberReq struct {
	g.Meta `path:"/circles/{circle_id}/members/{member_id}" method:"delete" tags:"Circle" summary:"删除圈子成员"`
}
type UpsertMemberReq struct {
	g.Meta `path:"/circles/{circle_id}/members/{member_id}" method:"put" tags:"Circle" summary:"更新圈子成员"`
}
type RedeemRewardReq struct {
	g.Meta `path:"/circles/{circle_id}/reward-redemptions" method:"post" tags:"Redemption" summary:"兑换奖励"`
}
type DeleteRewardReq struct {
	g.Meta `path:"/circles/{circle_id}/rewards/{reward_id}" method:"delete" tags:"Reward" summary:"删除奖励"`
}
type UpsertRewardReq struct {
	g.Meta `path:"/circles/{circle_id}/rewards/{reward_id}" method:"put" tags:"Reward" summary:"更新奖励"`
}
type GetRewardEligibilityReq struct {
	g.Meta `path:"/circles/{circle_id}/rewards/{reward_id}/eligibility" method:"get" tags:"Reward" summary:"获取奖励资格"`
}
type AdjustMemberStarsReq struct {
	g.Meta `path:"/circles/{circle_id}/star-adjustments" method:"post" tags:"Ledger" summary:"调整成员星星"`
}
type ListStarTransactionsReq struct {
	g.Meta `path:"/circles/{circle_id}/star-transactions" method:"get" tags:"Ledger" summary:"列出星星流水"`
}
type GetStatisticsReq struct {
	g.Meta `path:"/circles/{circle_id}/statistics" method:"get" tags:"ReadModels" summary:"获取统计"`
}
type CompareStatisticsReq struct {
	g.Meta `path:"/circles/{circle_id}/statistics:compare" method:"get" tags:"ReadModels" summary:"对比统计"`
}
type PullCircleBootstrapDeltaReq struct {
	g.Meta `path:"/circles/{circle_id}/sync" method:"get" tags:"Sync" summary:"拉取圈子增量"`
}
type ListTaskCompletionDetailsReq struct {
	g.Meta `path:"/circles/{circle_id}/task-completions" method:"get" tags:"ReadModels" summary:"列出任务完成明细"`
}
type CompleteTaskReq struct {
	g.Meta `path:"/circles/{circle_id}/task-completions" method:"post" tags:"TaskLedger" summary:"完成任务"`
}
type CancelTaskCompletionReq struct {
	g.Meta `path:"/circles/{circle_id}/task-completions/{completion_id}:cancel" method:"post" tags:"TaskLedger" summary:"取消任务完成"`
}
type ListTaskOccurrencesReq struct {
	g.Meta `path:"/circles/{circle_id}/task-occurrences" method:"get" tags:"Task" summary:"列出任务 occurrence"`
}
type DeleteTaskTagReq struct {
	g.Meta `path:"/circles/{circle_id}/task-tags/{task_tag_id}" method:"delete" tags:"Task" summary:"删除任务标签"`
}
type UpsertTaskTagReq struct {
	g.Meta `path:"/circles/{circle_id}/task-tags/{task_tag_id}" method:"put" tags:"Task" summary:"更新任务标签"`
}
type DeleteTaskReq struct {
	g.Meta `path:"/circles/{circle_id}/tasks/{task_id}" method:"delete" tags:"Task" summary:"删除任务"`
}
type UpsertTaskReq struct {
	g.Meta `path:"/circles/{circle_id}/tasks/{task_id}" method:"put" tags:"Task" summary:"更新任务"`
}
type LeaveCircleReq struct {
	g.Meta `path:"/circles/{circle_id}:leave" method:"post" tags:"Circle" summary:"退出圈子"`
}
type CompleteOnboardingReq struct {
	g.Meta `path:"/circles:onboard" method:"post" tags:"Circle" summary:"完成 onboarding"`
}
type GetCurrentEntitlementReq struct {
	g.Meta `path:"/entitlements/current" method:"get" tags:"Entitlement" summary:"获取当前权益"`
}
type VerifyPlayPurchaseReq struct {
	g.Meta `path:"/entitlements/google-play:verify" method:"post" tags:"Entitlement" summary:"验证 Play purchase"`
}
type RedeemAdministratorInviteReq struct {
	g.Meta `path:"/invites/administrator:redeem" method:"post" tags:"Invite" summary:"兑换管理员邀请"`
}
type RedeemMemberInviteReq struct {
	g.Meta `path:"/invites/member:redeem" method:"post" tags:"Invite" summary:"兑换成员邀请"`
}
type SubmitFeedbackReq struct {
	g.Meta `path:"/support/feedback" method:"post" tags:"Support" summary:"提交反馈"`
}
