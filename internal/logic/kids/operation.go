package kids

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commonlogin "rslytics-app-api/internal/common/login"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// v1Operation 是单个接口路由对应的领域操作。
type v1Operation func(context.Context, v1.V1OperationInput) (map[string]any, string, error)

// runV1 执行单个接口路由的公共认证、幂等和响应校验流程。
func (s *sKids) runV1(ctx context.Context, in v1.V1OperationInput, operation v1Operation) (*v1.V1OperationOutput, error) {
	if err := validateV1Input(in); err != nil {
		return nil, err
	}
	bodyFingerprint := v1BodyFingerprint(in.Body)
	routeFingerprint := v1RouteFingerprint(in)
	if in.OperationID == "refreshSession" {
		if output, conflict, pending, err := v1RefreshIdempotencyReplay(ctx, in, routeFingerprint, bodyFingerprint); err != nil {
			return nil, err
		} else if conflict {
			return nil, v1Error(409, "IDEMPOTENCY_CONFLICT", false, "idempotency key conflicts with an earlier request")
		} else if pending {
			return nil, v1Error(503, "UNAVAILABLE", true, "v1 request is still being committed")
		} else if output != nil {
			return output, nil
		}
	}
	if err := resolveV1Principal(ctx, &in); err != nil {
		return nil, err
	}
	principalScope := v1PrincipalScope(ctx, in)
	if v1OperationSupportsIdempotency(in) {
		if output, conflict, pending, err := v1IdempotencyBegin(ctx, principalScope, in, routeFingerprint, bodyFingerprint); err != nil {
			return nil, err
		} else if conflict {
			return nil, v1Error(409, "IDEMPOTENCY_CONFLICT", false, "idempotency key conflicts with an earlier request")
		} else if pending {
			return nil, v1Error(503, "UNAVAILABLE", true, "v1 request is still being committed")
		} else if output != nil {
			output = v1ReplayOutput(output)
			if err = v1.ValidateV1ResponseData(in.OperationID, output.Data); err != nil {
				v1LogResponseProjectionFailure(ctx, in, "idempotency_replay", err)
				return nil, v1Error(502, "PROTOCOL_ERROR", false, "stored v1 response violates the protocol")
			}
			return output, nil
		}
	}

	data, changeCursor, err := operation(ctx, in)
	if err != nil {
		if v1OperationSupportsIdempotency(in) {
			_ = v1IdempotencyAbort(ctx, principalScope, in, routeFingerprint, bodyFingerprint)
		}
		return nil, err
	}
	if err = v1.ValidateV1ResponseData(in.OperationID, data); err != nil {
		v1LogResponseProjectionFailure(ctx, in, "operation_result", err)
		if v1OperationSupportsIdempotency(in) {
			_ = v1IdempotencyAbort(ctx, principalScope, in, routeFingerprint, bodyFingerprint)
		}
		return nil, v1Error(502, "PROTOCOL_ERROR", false, "v1 response projection violates the protocol")
	}
	output := &v1.V1OperationOutput{Data: data, Status: 200, ChangeCursor: changeCursor}
	if v1OperationSupportsIdempotency(in) {
		if err = v1IdempotencySave(ctx, principalScope, in, routeFingerprint, bodyFingerprint, output); err != nil {
			return nil, err
		}
	}
	return output, nil
}

// v1LogResponseProjectionFailure 记录响应合同失败的关联信息，不记录 credential、proof 或完整请求正文。
func v1LogResponseProjectionFailure(ctx context.Context, in v1.V1OperationInput, stage string, err error) {
	traceID := ""
	if request := ghttp.RequestFromCtx(ctx); request != nil {
		traceID = request.GetCtxVar(consts.CtxTraceIDKey).String()
	}
	g.Log().Errorf(ctx, "event=kids_v1_response_projection_invalid operation_id=%s request_id=%s trace_id=%s stage=%s error=%v", in.OperationID, in.RequestID, traceID, stage, err)
}

// GetCurrentAccount 获取当前账号接口 bootstrap。
func (s *sKids) GetAccountBootstrap(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1GetCurrentAccount)
}

// SelectCurrentCircle 选择当前接口圈子。
func (s *sKids) SelectCircle(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1SelectCurrentCircle)
}

// CreateInviteGuestSession 创建接口游客会话。
func (s *sKids) CreateGuestSession(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CreateGuestSession)
}

// RefreshSession 刷新接口会话。
func (s *sKids) RefreshV1Session(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1RefreshSession)
}

// RevokeSession 撤销接口会话。
func (s *sKids) RevokeSession(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1RevokeSession)
}

// ListMyCircles 查询接口圈子列表。
func (s *sKids) ListV1Circles(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1ListV1Circles)
}

// GetCircleBootstrap 获取接口圈子 bootstrap。
func (s *sKids) GetV1CircleBootstrap(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1GetCircleBootstrap)
}

// UpdateCircleWithVersion 更新接口圈子。
func (s *sKids) UpdateCircleWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1UpdateCircle)
}

// DeleteCircleWithVersion 删除接口圈子。
func (s *sKids) DeleteCircleWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1DeleteCircle)
}

// LeaveCircleWithVersion 退出接口圈子。
func (s *sKids) LeaveCircleWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1LeaveCircle)
}

// CreateCircleMember 创建接口成员。
func (s *sKids) CreateCircleMember(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CreateCircleMember)
}

// UpsertCircleMember 更新接口成员。
func (s *sKids) UpsertCircleMember(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1UpsertMember)
}

// DeleteCircleMember 删除接口成员。
func (s *sKids) DeleteCircleMember(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1DeleteMember)
}

// UpsertCircleAdministrator 更新接口管理员。
func (s *sKids) UpsertCircleAdministrator(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1UpsertAdministrator)
}

// DeleteCircleAdministrator 删除接口管理员。
func (s *sKids) DeleteCircleAdministrator(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1DeleteAdministrator)
}

// CreateCircleInvite 创建接口邀请。
func (s *sKids) CreateCircleInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CreateInvite)
}

// RefreshCircleInvite 刷新接口邀请。
func (s *sKids) RefreshCircleInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1RefreshInvite)
}

// RevokeCircleInvite 撤销接口邀请。
func (s *sKids) RevokeCircleInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1RevokeInvite)
}

// RedeemAdministratorInvite 兑换管理员接口邀请。
func (s *sKids) RedeemAdministratorInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, func(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
		return s.v1RedeemInvite(ctx, in, "administrator")
	})
}

// RedeemMemberInvite 兑换成员接口邀请。
func (s *sKids) RedeemMemberInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, func(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
		return s.v1RedeemInvite(ctx, in, "member")
	})
}

// GetCurrentEntitlement 获取接口权益。
func (s *sKids) GetEntitlement(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1GetCurrentEntitlement)
}

// SubmitFeedback 提交接口反馈。
func (s *sKids) SubmitFeedbackV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1SubmitFeedback)
}

// CompleteOnboarding 完成接口 onboarding。
func (s *sKids) CompleteOnboardingV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CompleteOnboarding)
}

// UpsertTaskTagWithVersion 更新接口任务标签。
func (s *sKids) UpsertTaskTagWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1UpsertTaskTag)
}

// DeleteTaskTagWithVersion 删除接口任务标签。
func (s *sKids) DeleteTaskTagWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1DeleteTaskTag)
}

// UpsertTaskWithVersion 更新接口任务。
func (s *sKids) UpsertTaskWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1UpsertTask)
}

// DeleteTaskWithVersion 删除接口任务。
func (s *sKids) DeleteTaskWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1DeleteTask)
}

// ListTaskOccurrences 查询接口任务 occurrence。
func (s *sKids) ListTaskOccurrencesV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1TaskOccurrences)
}

// CompleteTaskWithVersion 完成接口任务。
func (s *sKids) CompleteTaskWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CompleteTask)
}

// CancelTaskCompletionWithVersion 取消接口任务完成记录。
func (s *sKids) CancelTaskCompletionWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CancelTaskCompletion)
}

// ListTaskCompletionDetails 查询接口任务完成明细。
func (s *sKids) ListTaskCompletionDetailsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1CompletionDetails)
}

// AdjustMemberStarsWithVersion 调整接口成员星星。
func (s *sKids) AdjustMemberStarsWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.v1AdjustMemberStars)
}

// Unavailable 返回尚未实现的接口 operation 错误。
func (s *sKids) Unavailable(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, func(context.Context, v1.V1OperationInput) (map[string]any, string, error) {
		return nil, "", v1Error(503, "UNAVAILABLE", true, "v1 operation is not yet available")
	})
}

// ExecuteV1 执行接口中暂未拆分专用 service 入口的 operation，保证每个路由均进入对应领域实现。
func (s *sKids) ExecuteV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, s.executeV1Operation)
}

// executeV1Operation 根据冻结 operationId 路由到接口领域实现，未知 operation 一律失败关闭。
func (s *sKids) executeV1Operation(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	switch in.OperationID {
	case "exchangeGoogleProof":
		return s.v1ExchangeGoogleProof(ctx, in)
	case "prepareAssetUpload":
		return s.v1PrepareAsset(ctx, in)
	case "commitAssetUpload":
		return s.v1CommitAsset(ctx, in)
	case "listExchangeHistory":
		return s.v1ExchangeHistory(ctx, in)
	case "redeemReward":
		return s.v1RedeemReward(ctx, in)
	case "deleteReward":
		return s.v1DeleteReward(ctx, in)
	case "upsertReward":
		return s.v1UpsertReward(ctx, in)
	case "getRewardEligibility":
		return s.v1RewardEligibility(ctx, in)
	case "pullCircleBootstrapDelta":
		return s.v1PullCircleBootstrapDelta(ctx, in)
	case "verifyPlayPurchase":
		return s.v1VerifyPlayPurchase(ctx, in)
	default:
		return nil, "", v1Error(404, "NOT_FOUND", false, "protocol operation is not registered")
	}
}

// v1ReplayOutput 复制首次快照，并将存在的接口回执标记为安全幂等重放。
func v1ReplayOutput(output *v1.V1OperationOutput) *v1.V1OperationOutput {
	if output == nil {
		return nil
	}
	encoded, err := json.Marshal(output.Data)
	if err != nil {
		return output
	}
	var copied map[string]any
	if err = json.Unmarshal(encoded, &copied); err != nil {
		return output
	}
	if receipt, ok := copied["receipt"].(map[string]any); ok {
		receipt["result_kind"] = "idempotent_replay"
	}
	return &v1.V1OperationOutput{Data: copied, Status: output.Status, ChangeCursor: output.ChangeCursor, ETag: output.ETag}
}

// v1OperationSupportsIdempotency 仅允许已经具备完整重放快照的 operation 进入幂等提交流程。
func v1OperationSupportsIdempotency(in v1.V1OperationInput) bool {
	return in.Method != "GET" && in.Method != "HEAD" && (in.OperationID == "createInviteGuestSession" || in.OperationID == "refreshSession" || in.OperationID == "revokeSession" || in.OperationID == "createInvite" || in.OperationID == "refreshInvite" || in.OperationID == "revokeInvite" || in.OperationID == "redeemAdministratorInvite" || in.OperationID == "redeemMemberInvite" || in.OperationID == "upsertTaskTag" || in.OperationID == "deleteTaskTag" || in.OperationID == "upsertTask" || in.OperationID == "deleteTask" || in.OperationID == "completeTask" || in.OperationID == "cancelTaskCompletion" || in.OperationID == "adjustMemberStars" || in.OperationID == "completeOnboarding" || in.OperationID == "selectCurrentCircle" || in.OperationID == "updateCircle" || in.OperationID == "deleteCircle" || in.OperationID == "createCircleMember" || in.OperationID == "upsertMember" || in.OperationID == "deleteMember" || in.OperationID == "upsertAdministrator" || in.OperationID == "deleteAdministrator" || in.OperationID == "leaveCircle" || in.OperationID == "submitFeedback")
}

// resolveV1Principal 从持久化 session 解析接口认证上下文，拒绝过期、撤销或类型不匹配的凭据。
func resolveV1Principal(ctx context.Context, in *v1.V1OperationInput) error {
	authContext := v1OperationAuthContext(in.OperationID)
	if authContext == "public" || (authContext == "public_or_guest" && strings.TrimSpace(in.AccessToken) == "") {
		return nil
	}
	tokenField := "access_token_hash"
	if authContext == "refresh" {
		tokenField = "refresh_token_hash"
	}
	var accessClaims *utils.V1SessionTokenClaims
	if authContext != "refresh" {
		var err error
		accessClaims, err = utils.ParseV1SessionToken(in.AccessToken, v1SessionSigningSecret(ctx), v1SessionEnvironment(ctx))
		if err != nil {
			return v1Error(401, "UNAUTHENTICATED", false, "session credential is invalid")
		}
	}
	sessionModel := utils.KidsDB(ctx).Model(consts.KidsV1SessionTable).Ctx(ctx).
		Where(tokenField, sha256Hex(in.AccessToken))
	if authContext != "refresh" {
		sessionModel = sessionModel.Where("status", "active")
	}
	session, err := sessionModel.One()
	if err != nil {
		return err
	}
	if session.IsEmpty() {
		return v1Error(401, "UNAUTHENTICATED", false, "session credential is invalid")
	}
	nowMs := time.Now().UnixMilli()
	expiresAtMs := session["access_expires_at_ms"].Int64()
	if authContext == "refresh" {
		expiresAtMs = session["refresh_expires_at_ms"].Int64()
	}
	if expiresAtMs <= nowMs {
		return v1Error(401, "TOKEN_EXPIRED", false, "session credential is expired")
	}
	in.SessionID = session["session_id"].String()
	in.PrincipalKind = session["principal_kind"].String()
	in.PrincipalID = session["account_id"].String()
	if accessClaims != nil && (accessClaims.SessionID != in.SessionID || accessClaims.AccountID != in.PrincipalID || accessClaims.PrincipalKind != in.PrincipalKind || accessClaims.ExpiresAtMs != expiresAtMs) {
		return v1Error(401, "UNAUTHENTICATED", false, "session credential is invalid")
	}
	if authContext == "refresh" {
		if in.PrincipalKind != "account" || in.SessionID != fmt.Sprint(in.Body["session_id"]) {
			return v1Error(401, "UNAUTHENTICATED", false, "refresh credential does not match the account session")
		}
		return nil
	}
	if authContext == "account" && in.PrincipalKind != "account" {
		return v1Error(403, "FORBIDDEN", false, "account session is required")
	}
	if authContext == "account" || in.PrincipalKind == "account" {
		account, err := utils.KidsDB(ctx).Model(consts.KidsV1AccountTable).Ctx(ctx).Where("account_id", in.PrincipalID).One()
		if err != nil {
			return err
		}
		if account.IsEmpty() {
			return v1Error(401, "UNAUTHENTICATED", false, "v1 account is missing")
		}
		if account["status"].String() != "active" {
			return v1Error(403, "ACCOUNT_DISABLED", false, "v1 account is disabled")
		}
	}
	if authContext == "guest_or_account" || authContext == "public_or_guest" || authContext == "asset_owner" || authContext == "membership" {
		return nil
	}
	return nil
}

// v1OperationAuthContext 返回冻结 operation 的认证语义，OpenAPI request 校验已保证 operation ID 合法。
func v1OperationAuthContext(operationID string) string {
	switch operationID {
	case "createInviteGuestSession":
		return "public"
	case "exchangeGoogleProof":
		return "public_or_guest"
	case "refreshSession":
		return "refresh"
	case "getCurrentAccount", "redeemAdministratorInvite", "redeemMemberInvite", "submitFeedback":
		return "guest_or_account"
	case "commitAssetUpload", "prepareAssetUpload":
		return "asset_owner"
	case "revokeSession", "selectCurrentCircle", "listMyCircles", "completeOnboarding", "getCurrentEntitlement", "verifyPlayPurchase":
		return "account"
	default:
		return "membership"
	}
}

// validateV1Input 校验接口输入的基础传输字段。
func validateV1Input(in v1.V1OperationInput) error {
	if in.OperationID == "" || in.Method == "" || in.RequestID == "" {
		return v1Error(422, "VALIDATION_FAILED", false, "v1 request metadata is incomplete")
	}
	return nil
}

// v1GetCurrentEntitlement 返回持久化权益快照；首次查询会建立 inactive 初始快照。
func (s *sKids) v1GetCurrentEntitlement(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	row, err := utils.KidsDB(ctx).Model(consts.KidsV1EntitlementTable).Ctx(ctx).Where("account_id", in.PrincipalID).One()
	if err != nil {
		return nil, "", err
	}
	if row.IsEmpty() {
		now := time.Now()
		entitlementID := v1ID("entitlement", in.PrincipalID)
		if _, err = utils.KidsDB(ctx).Model(consts.KidsV1EntitlementTable).Ctx(ctx).Data(gdb.Map{"entitlement_id": entitlementID, "account_id": in.PrincipalID, "status": "inactive", "version": 1, "verified_at": now, "created_at": now, "updated_at": now}).Insert(); err != nil {
			return nil, "", err
		}
		return map[string]any{"entitlement_id": entitlementID, "account_id": in.PrincipalID, "status": "inactive", "trial_eligible": true, "product_id": nil, "plan_code": nil, "base_plan_id": nil, "offer_id": nil, "billing_phase": "none", "valid_until_ms": nil, "verified_at_ms": now.UnixMilli(), "revoked_at_ms": nil, "version": int64(1)}, "", nil
	}
	billingPhase := "unknown"
	if row["status"].String() == "inactive" {
		billingPhase = "none"
	}
	return map[string]any{"entitlement_id": row["entitlement_id"].String(), "account_id": row["account_id"].String(), "status": row["status"].String(), "trial_eligible": row["status"].String() == "inactive", "product_id": nullableString(row["product_id"].String()), "plan_code": nil, "base_plan_id": nil, "offer_id": nil, "billing_phase": billingPhase, "valid_until_ms": v1NullableTimeMillis(row["valid_until_at"].Time()), "verified_at_ms": row["verified_at"].Time().UnixMilli(), "revoked_at_ms": v1NullableTimeMillis(row["revoked_at"].Time()), "version": row["version"].Int64()}, "", nil
}

// v1GetCircleBootstrap 返回指定圈子的同一快照分页，所有实体均来自接口事实表。
func (s *sKids) v1GetCircleBootstrap(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	circle, err := utils.KidsDB(ctx).Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").One()
	if err != nil {
		return nil, "", err
	}
	if circle.IsEmpty() {
		return nil, "", v1Error(404, "NOT_FOUND", false, "circle is missing")
	}
	limit, err := strconv.Atoi(v1QueryFirst(in.Query, "limit"))
	if err != nil || limit < 1 || limit > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "bootstrap limit is invalid")
	}
	offset, err := v1PageOffset(v1QueryFirst(in.Query, "cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "bootstrap cursor is invalid")
	}
	memberships, err := utils.KidsDB(ctx).Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").Order("created_at ASC,id ASC").Limit(offset, limit+1).All()
	if err != nil {
		return nil, "", err
	}
	hasMore := len(memberships) > limit
	if hasMore {
		memberships = memberships[:limit]
	}
	admins, err := utils.KidsDB(ctx).Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("circle_id", circleID).Where("status <> ?", "deleted").Order("created_at ASC,id ASC").All()
	if err != nil {
		return nil, "", err
	}
	members, err := utils.KidsDB(ctx).Model(consts.KidsV1MemberTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").Order("created_at ASC,id ASC").All()
	if err != nil {
		return nil, "", err
	}
	adminOut := make([]map[string]any, 0, len(admins))
	for _, row := range admins {
		adminOut = append(adminOut, v1AdministratorRecordProjection(row))
	}
	memberOut := make([]map[string]any, 0, len(members))
	for _, row := range members {
		memberOut = append(memberOut, v1MemberRecordProjection(row))
	}
	membershipOut := make([]map[string]any, 0, len(memberships))
	for _, row := range memberships {
		membershipOut = append(membershipOut, v1MembershipRecordProjection(row))
	}
	next := any(nil)
	if hasMore {
		next = v1PageCursor(offset + limit)
	}
	snapshot, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"circle": v1CircleRecordProjection(circle), "memberships": membershipOut, "administrators": adminOut, "members": memberOut, "next_cursor": next, "has_more": hasMore, "snapshot_cursor": snapshot}, "", nil
}

// v1UpdateCircle 按 expected_version 更新圈子名称和视觉引用，并记录完整 commit。
func (s *sKids) v1UpdateCircle(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageCircle); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	iconJSON, err := json.Marshal(in.Body["icon"])
	if err != nil {
		return nil, "", err
	}
	now := time.Now()
	var circle map[string]any
	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageCircle); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "circle is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "circle version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Data(gdb.Map{"name": in.Body["name"], "icon": string(iconJSON), "version": version, "updated_at": now}).Update(); e != nil {
			return e
		}
		row, e = tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).One()
		if e != nil {
			return e
		}
		circle = v1CircleRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"circle": circle})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "circle": circle}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1CreateCircleMember 创建初始成员并递增圈子版本，拒绝客户端伪造绑定账号。
func (s *sKids) v1CreateCircleMember(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_circle_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected circle version is invalid")
	}
	member, ok := in.Body["member"].(map[string]any)
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "member is invalid")
	}
	memberID := fmt.Sprint(member["member_id"])
	avatarJSON, err := json.Marshal(member["avatar"])
	if err != nil {
		return nil, "", err
	}
	now := time.Now()
	var out map[string]any
	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); e != nil {
			return e
		}
		circle, e := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if e != nil {
			return e
		}
		if circle.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "circle is missing")
		}
		cv := circle["version"].Int64()
		if cv != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &cv, Message: "circle version conflicts"}
		}
		existing, e := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).One()
		if e != nil {
			return e
		}
		if !existing.IsEmpty() {
			return v1Error(409, "IDEMPOTENCY_CONFLICT", false, "member already exists")
		}
		if _, e = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Data(gdb.Map{"member_id": memberID, "circle_id": circleID, "display_name": member["display_name"], "gender": member["gender"], "avatar": string(avatarJSON), "status": "active", "version": 1, "created_at": now, "updated_at": now}).Insert(); e != nil {
			return e
		}
		cv++
		if _, e = tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Data(gdb.Map{"version": cv, "updated_at": now}).Update(); e != nil {
			return e
		}
		row, _ := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).One()
		out = v1MemberRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"member": out})
		if ce != nil {
			return ce
		}
		_, ce = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Data(gdb.Map{
			"circle_id": circleID, "member_id": memberID, "balance": 0, "version": 1,
			"source_commit_id": receipt["commit_id"], "source_commit_sequence": receipt["commit_sequence"], "updated_at": now,
		}).Insert()
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "member": out, "sync_cursor": cursor}, cursor, nil
}

// v1UpsertMember 更新已有成员资料并检查成员版本。
func (s *sKids) v1UpsertMember(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	memberID := in.PathParameters["member_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	avatarJSON, err := json.Marshal(in.Body["avatar"])
	if err != nil {
		return nil, "", err
	}
	now := time.Now()
	var out map[string]any
	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "member is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "member version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Data(gdb.Map{"display_name": in.Body["display_name"], "gender": in.Body["gender"], "avatar": string(avatarJSON), "version": version, "updated_at": now}).Update(); e != nil {
			return e
		}
		row, e = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).One()
		if e != nil {
			return e
		}
		out = v1MemberRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"member": out})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "member": out}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1DeleteMember 将成员软删除并保留接口 tombstone；任务和奖励领域尚未接入时不会存在关联指派。
func (s *sKids) v1DeleteMember(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	memberID := in.PathParameters["member_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	now := time.Now()
	var tombstone map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers)
		if e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "member is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "member version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Data(gdb.Map{"status": "deleted", "version": version, "deleted_at": now, "updated_at": now}).Update(); e != nil {
			return e
		}
		actor, e := v1ActorSnapshotTx(ctx, tx, membership)
		if e != nil {
			return e
		}
		tombstone = v1EntityTombstone("member", memberID, version, now, actor)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"member_tombstone": tombstone})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "member_tombstone": tombstone, "affected_task_assignments": int64(0), "affected_reward_assignments": int64(0)}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1UpsertAdministrator 创建待邀请管理员或更新已有管理员资料，权限和 Owner 归属始终由服务端决定。
func (s *sKids) v1UpsertAdministrator(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	adminID := in.PathParameters["administrator_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	intent := fmt.Sprint(in.Body["profile_intent"])
	expected, hasExpected := v1ExpectedVersion(in.Body["expected_version"])
	if !hasExpected && in.Body["expected_version"] != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	avatarJSON, err := json.Marshal(in.Body["avatar"])
	if err != nil {
		return nil, "", err
	}
	now := time.Now()
	var out map[string]any
	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).Where("circle_id", circleID).LockUpdate().One()
		if e != nil {
			return e
		}
		if intent == "create_administrator" {
			if hasExpected {
				return v1Error(422, "VALIDATION_FAILED", false, "new administrator must use null expected version")
			}
			if !row.IsEmpty() {
				return v1Error(409, "IDEMPOTENCY_CONFLICT", false, "administrator already exists")
			}
			emptyPermissions, _ := json.Marshal([]string{})
			if _, e = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Data(gdb.Map{"administrator_id": adminID, "circle_id": circleID, "display_name": in.Body["display_name"], "avatar": string(avatarJSON), "role": "administrator", "permissions": string(emptyPermissions), "status": "pending_invite", "version": 1, "created_at": now, "updated_at": now}).Insert(); e != nil {
				return e
			}
			row, _ = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).One()
		} else if intent == "update_profile" {
			if !hasExpected {
				return v1Error(422, "VALIDATION_FAILED", false, "administrator update requires expected version")
			}
			if row.IsEmpty() || row["status"].String() == "deleted" {
				return v1Error(404, "NOT_FOUND", false, "administrator is missing")
			}
			current := row["version"].Int64()
			if current != expected {
				return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "administrator version conflicts"}
			}
			version := current + 1
			if _, e = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).Data(gdb.Map{"display_name": in.Body["display_name"], "avatar": string(avatarJSON), "version": version, "updated_at": now}).Update(); e != nil {
				return e
			}
			row, e = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).One()
			if e != nil {
				return e
			}
		} else {
			return v1Error(422, "VALIDATION_FAILED", false, "administrator intent is invalid")
		}
		out = v1AdministratorRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"administrator": out})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "administrator": out}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1DeleteAdministrator 软删除非 Owner 管理员，保留 Owner slot 以满足每圈至少一个 Owner 的约束。
func (s *sKids) v1DeleteAdministrator(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	adminID := in.PathParameters["administrator_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	now := time.Now()
	var tombstone map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers)
		if e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).Where("circle_id", circleID).Where("status <> ?", "deleted").LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "administrator is missing")
		}
		if row["role"].String() == "owner" {
			return v1Error(403, "FORBIDDEN", false, "owner administrator cannot be deleted")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "administrator version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).Data(gdb.Map{"status": "deleted", "version": version, "deleted_at": now, "updated_at": now}).Update(); e != nil {
			return e
		}
		actor, e := v1ActorSnapshotTx(ctx, tx, membership)
		if e != nil {
			return e
		}
		tombstone = v1EntityTombstone("administrator", adminID, version, now, actor)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"administrator_tombstone": tombstone})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "administrator_tombstone": tombstone}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1DeleteCircle 软删除圈子并为当前账号生成可访问的 fallback selection。
func (s *sKids) v1DeleteCircle(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageCircle); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	now := time.Now()
	var tombstone map[string]any
	var selection map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageCircle)
		if e != nil {
			return e
		}
		circle, e := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if e != nil {
			return e
		}
		if circle.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "circle is missing")
		}
		current := circle["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "circle version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Data(gdb.Map{"status": "deleted", "version": version, "deleted_at": now, "updated_at": now}).Update(); e != nil {
			return e
		}
		if e = v1SoftDeleteCircleDependentsTx(ctx, tx, circleID, in.PrincipalID, now); e != nil {
			return e
		}
		actor, e := v1ActorSnapshotTx(ctx, tx, membership)
		if e != nil {
			return e
		}
		tombstone = v1EntityTombstone("circle", circleID, version, now, actor)
		sel, e := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).LockUpdate().One()
		if e != nil {
			return e
		}
		if sel.IsEmpty() {
			fallback, fallbackErr := v1FallbackCircleTx(ctx, tx, in.PrincipalID, circleID, "")
			if fallbackErr != nil {
				return fallbackErr
			}
			selection = v1SelectionProjection(v1ID("selection", uuid.NewString()), in.PrincipalID, fallback, 1, now)
			_, e = tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Data(gdb.Map{"selection_id": selection["selection_id"], "account_id": in.PrincipalID, "current_circle_id": fallback, "version": 1, "created_at": now, "updated_at": now}).Insert()
		} else {
			sv := sel["version"].Int64() + 1
			fallback, fallbackErr := v1FallbackCircleTx(ctx, tx, in.PrincipalID, circleID, "")
			if fallbackErr != nil {
				return fallbackErr
			}
			selection = v1SelectionProjection(sel["selection_id"].String(), in.PrincipalID, fallback, sv, now)
			_, e = tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).Data(gdb.Map{"current_circle_id": fallback, "version": sv, "updated_at": now}).Update()
		}
		if e != nil {
			return e
		}
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"circle_tombstone": tombstone, "selection": selection})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "circle_tombstone": tombstone, "selection": selection}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1LeaveCircle 将当前账号 membership 标记为 left，并清理不可访问的当前选择。
func (s *sKids) v1LeaveCircle(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	expected, ok := v1ExpectedVersion(in.Body["expected_membership_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected membership version is invalid")
	}
	now := time.Now()
	var tombstone map[string]any
	var selection map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, "")
		if e != nil {
			return e
		}
		if membership["role"].String() == "owner" {
			return v1Error(403, "FORBIDDEN", false, "owner must delete circle instead of leaving")
		}
		current := membership["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "membership version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("membership_id", membership["membership_id"].String()).Data(gdb.Map{"status": "left", "version": version, "deleted_at": now, "updated_at": now}).Update(); e != nil {
			return e
		}
		actor, e := v1ActorSnapshotTx(ctx, tx, membership)
		if e != nil {
			return e
		}
		tombstone = v1EntityTombstone("membership", membership["membership_id"].String(), version, now, actor)
		sel, e := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).LockUpdate().One()
		if e != nil {
			return e
		}
		if sel.IsEmpty() {
			fallback, fallbackErr := v1FallbackCircleTx(ctx, tx, in.PrincipalID, circleID, fmt.Sprint(in.Body["preferred_fallback_circle_id"]))
			if fallbackErr != nil {
				return fallbackErr
			}
			selection = v1SelectionProjection(v1ID("selection", uuid.NewString()), in.PrincipalID, fallback, 1, now)
			_, e = tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Data(gdb.Map{"selection_id": selection["selection_id"], "account_id": in.PrincipalID, "current_circle_id": fallback, "version": 1, "created_at": now, "updated_at": now}).Insert()
		} else {
			sv := sel["version"].Int64() + 1
			fallback, fallbackErr := v1FallbackCircleTx(ctx, tx, in.PrincipalID, circleID, fmt.Sprint(in.Body["preferred_fallback_circle_id"]))
			if fallbackErr != nil {
				return fallbackErr
			}
			selection = v1SelectionProjection(sel["selection_id"].String(), in.PrincipalID, fallback, sv, now)
			_, e = tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).Data(gdb.Map{"current_circle_id": fallback, "version": sv, "updated_at": now}).Update()
		}
		if e != nil {
			return e
		}
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"membership_tombstone": tombstone, "selection": selection})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "membership_tombstone": tombstone, "selection": selection}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1RequireMembershipPermission 读取 active membership 并校验权限，避免信任客户端 role。
func v1RequireMembershipPermission(ctx context.Context, accountID, circleID, permission string) error {
	_, err := v1RequireMembership(ctx, accountID, circleID, permission)
	return err
}

// v1RequireMembership 读取当前账号在圈子的 active membership。
func v1RequireMembership(ctx context.Context, accountID, circleID, permission string) (gdb.Record, error) {
	row, err := utils.KidsDB(ctx).Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("account_id", accountID).Where("circle_id", circleID).Where("status", "active").One()
	if err != nil {
		return nil, err
	}
	if row.IsEmpty() {
		return nil, v1Error(403, "FORBIDDEN", false, "active circle membership is required")
	}
	if permission != "" && !v1PermissionsContain(row["permissions"].String(), permission) {
		return nil, v1Error(403, "FORBIDDEN", false, "membership permission is required")
	}
	return row, nil
}

// v1RequireMembershipTx 在已有事务中读取并锁定当前账号 membership。
func v1RequireMembershipTx(ctx context.Context, tx gdb.TX, accountID, circleID, permission string) (gdb.Record, error) {
	row, err := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("account_id", accountID).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
	if err != nil {
		return nil, err
	}
	if row.IsEmpty() {
		return nil, v1Error(403, "FORBIDDEN", false, "active circle membership is required")
	}
	if permission != "" && !v1PermissionsContain(row["permissions"].String(), permission) {
		return nil, v1Error(403, "FORBIDDEN", false, "membership permission is required")
	}
	return row, nil
}

// v1PermissionsContain 检查 JSON 权限数组中的单项权限。
func v1PermissionsContain(raw, permission string) bool {
	value := v1JSONValue(raw)
	items, ok := value.([]any)
	if !ok {
		return false
	}
	for _, item := range items {
		if fmt.Sprint(item) == permission {
			return true
		}
	}
	return false
}

// v1ActorSnapshotTx 从锁定 membership 生成删除审计所需的 actor snapshot。
func v1ActorSnapshotTx(ctx context.Context, tx gdb.TX, membership gdb.Record) (map[string]any, error) {
	name := membership["actor_id"].String()
	if membership["actor_type"].String() == "owner" || membership["actor_type"].String() == "administrator" {
		admin, e := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", name).One()
		if e != nil {
			return nil, e
		}
		if !admin.IsEmpty() {
			name = admin["display_name"].String()
		}
	}
	return map[string]any{"actor_type": membership["actor_type"].String(), "actor_id": membership["actor_id"].String(), "role": membership["role"].String(), "display_name_snapshot": name}, nil
}

// v1EntityTombstone 构造接口通用实体 tombstone。
func v1EntityTombstone(entityType, entityID string, version int64, deletedAt time.Time, actor map[string]any) map[string]any {
	return map[string]any{"entity_type": entityType, "entity_id": entityID, "version": version, "deleted_at_ms": deletedAt.UnixMilli(), "deleted_by": actor}
}

// v1FallbackCircleTx 在当前事务内选择账号仍可访问的 preferred 或最早 active 圈子。
func v1FallbackCircleTx(ctx context.Context, tx gdb.TX, accountID, excludedCircleID, preferredCircleID string) (string, error) {
	rows, err := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("account_id", accountID).Where("status", "active").Order("created_at ASC,id ASC").All()
	if err != nil {
		return "", err
	}
	var first string
	for _, row := range rows {
		circleID := row["circle_id"].String()
		if circleID == excludedCircleID {
			continue
		}
		circle, queryErr := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").One()
		if queryErr != nil {
			return "", queryErr
		}
		if circle.IsEmpty() {
			continue
		}
		if circleID == preferredCircleID {
			return circleID, nil
		}
		if first == "" {
			first = circleID
		}
	}
	return first, nil
}

// v1SoftDeleteCircleDependentsTx 软删除圈子下的授权和资料，并清理其他账号指向该圈子的选择。
func v1SoftDeleteCircleDependentsTx(ctx context.Context, tx gdb.TX, circleID, actingAccountID string, now time.Time) error {
	memberships, err := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("circle_id", circleID).Where("status", "active").All()
	if err != nil {
		return err
	}
	for _, row := range memberships {
		version := row["version"].Int64() + 1
		if _, err = tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("membership_id", row["membership_id"].String()).Data(gdb.Map{"status": "deleted", "version": version, "deleted_at": now, "updated_at": now}).Update(); err != nil {
			return err
		}
		accountID := row["account_id"].String()
		if accountID == actingAccountID {
			continue
		}
		selection, queryErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", accountID).Where("current_circle_id", circleID).LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if selection.IsEmpty() {
			continue
		}
		fallback, fallbackErr := v1FallbackCircleTx(ctx, tx, accountID, circleID, "")
		if fallbackErr != nil {
			return fallbackErr
		}
		if _, err = tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("selection_id", selection["selection_id"].String()).Data(gdb.Map{"current_circle_id": fallback, "version": selection["version"].Int64() + 1, "updated_at": now}).Update(); err != nil {
			return err
		}
	}
	admins, err := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("circle_id", circleID).Where("status <> ?", "deleted").All()
	if err != nil {
		return err
	}
	for _, row := range admins {
		if _, err = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", row["administrator_id"].String()).Data(gdb.Map{"status": "deleted", "version": row["version"].Int64() + 1, "deleted_at": now, "updated_at": now}).Update(); err != nil {
			return err
		}
	}
	members, err := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("circle_id", circleID).Where("status <> ?", "deleted").All()
	if err != nil {
		return err
	}
	for _, row := range members {
		if _, err = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", row["member_id"].String()).Data(gdb.Map{"status": "deleted", "version": row["version"].Int64() + 1, "deleted_at": now, "updated_at": now}).Update(); err != nil {
			return err
		}
	}
	return nil
}

// v1AdministratorRecordProjection 将管理员记录还原为接口管理员实体。
func v1AdministratorRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"administrator_id": row["administrator_id"].String(), "circle_id": row["circle_id"].String(), "bound_account_id": nullableString(row["bound_account_id"].String()), "display_name": row["display_name"].String(), "avatar": v1JSONValue(row["avatar"].String()), "role": row["role"].String(), "permissions": v1JSONValue(row["permissions"].String()), "status": row["status"].String(), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "updated_at_ms": row["updated_at"].Time().UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(row["deleted_at"].Time())}
}

// v1MemberRecordProjection 将成员记录还原为接口成员实体。
func v1MemberRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"member_id": row["member_id"].String(), "circle_id": row["circle_id"].String(), "bound_account_id": nullableString(row["bound_account_id"].String()), "display_name": row["display_name"].String(), "gender": row["gender"].String(), "avatar": v1JSONValue(row["avatar"].String()), "status": row["status"].String(), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "updated_at_ms": row["updated_at"].Time().UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(row["deleted_at"].Time())}
}

// v1ListV1Circles 从接口成员身份和圈子事实读取当前账号的分页圈子列表。
func (s *sKids) v1ListV1Circles(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	limit, err := strconv.Atoi(v1QueryFirst(in.Query, "limit"))
	if err != nil || limit < 1 || limit > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "circle list limit is invalid")
	}
	offset, err := v1PageOffset(v1QueryFirst(in.Query, "cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "circle list cursor is invalid")
	}
	selection, err := utils.KidsDB(ctx).Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).One()
	if err != nil {
		return nil, "", err
	}
	memberships, err := utils.KidsDB(ctx).Model(consts.KidsV1MembershipTable).Ctx(ctx).
		Where("account_id", in.PrincipalID).Where("status", "active").Order("created_at ASC, id ASC").Limit(offset, limit+1).All()
	if err != nil {
		return nil, "", err
	}
	hasMore := len(memberships) > limit
	if hasMore {
		memberships = memberships[:limit]
	}
	items := make([]map[string]any, 0, len(memberships))
	for _, membershipRecord := range memberships {
		circle, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1CircleTable).Ctx(ctx).
			Where("circle_id", membershipRecord["circle_id"].String()).Where("status", "active").One()
		if queryErr != nil {
			return nil, "", queryErr
		}
		if circle.IsEmpty() {
			continue
		}
		items = append(items, map[string]any{
			"circle": v1CircleRecordProjection(circle), "membership": v1MembershipRecordProjection(membershipRecord),
			"selected": !selection.IsEmpty() && selection["current_circle_id"].String() == circle["circle_id"].String(),
		})
	}
	nextCursor := any(nil)
	if hasMore {
		nextCursor = v1PageCursor(offset + limit)
	}
	snapshotCursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"items": items, "next_cursor": nextCursor, "has_more": hasMore, "snapshot_cursor": snapshotCursor}, "", nil
}

// v1SelectCurrentCircle 校验账号在目标圈子的有效成员身份后原子更新当前选择。
func (s *sKids) v1SelectCurrentCircle(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := fmt.Sprint(in.Body["circle_id"])
	if circleID == "" {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "circle ID is invalid")
	}
	expectedVersion, hasExpectedVersion := v1ExpectedVersion(in.Body["expected_version"])
	if !hasExpectedVersion && in.Body["expected_version"] != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	now := time.Now()
	var selection map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, queryErr := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).
			Where("circle_id", circleID).Where("account_id", in.PrincipalID).Where("status", "active").LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if membership.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "circle membership is missing")
		}
		record, queryErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if !record.IsEmpty() && hasExpectedVersion && record["version"].Int64() != expectedVersion {
			version := record["version"].Int64()
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &version, Message: "circle selection version conflicts"}
		}
		if record.IsEmpty() {
			selection = v1SelectionProjection(v1ID("selection", uuid.NewString()), in.PrincipalID, circleID, 1, now)
			if _, insertErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Data(gdb.Map{
				"selection_id": selection["selection_id"], "account_id": in.PrincipalID, "current_circle_id": circleID, "version": 1, "created_at": now, "updated_at": now,
			}).Insert(); insertErr != nil {
				return insertErr
			}
		} else {
			version := record["version"].Int64() + 1
			selection = v1SelectionProjection(record["selection_id"].String(), in.PrincipalID, circleID, version, now)
			if _, updateErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).Data(gdb.Map{
				"current_circle_id": circleID, "version": version, "updated_at": now,
			}).Update(); updateErr != nil {
				return updateErr
			}
		}
		var commitErr error
		receipt, commitErr = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"selection": selection})
		return commitErr
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "selection": selection}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1ExpectedVersion 将 JSON schema 已验证的版本值转换为内部整数。
func v1ExpectedVersion(value any) (int64, bool) {
	switch typed := value.(type) {
	case json.Number:
		version, err := typed.Int64()
		if err == nil && version >= 1 {
			return version, true
		}
	case float64:
		if typed >= 1 && typed == float64(int64(typed)) {
			return int64(typed), true
		}
	case int64:
		if typed >= 1 {
			return typed, true
		}
	case int:
		if typed >= 1 {
			return int64(typed), true
		}
	case nil:
		return 0, false
	}
	return 0, false
}

// v1PageOffset 解析由本服务签发的分页游标；首请求可省略游标。
func v1PageOffset(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	const prefix = "cur:v1:page_"
	if !strings.HasPrefix(cursor, prefix) {
		return 0, fmt.Errorf("unsupported cursor")
	}
	offset, err := strconv.Atoi(strings.TrimPrefix(cursor, prefix))
	if err != nil || offset < 0 {
		return 0, fmt.Errorf("invalid cursor offset")
	}
	return offset, nil
}

// v1PageCursor 为接口分页响应签发下一页的服务端不透明游标。
func v1PageCursor(offset int) string { return fmt.Sprintf("cur:v1:page_%08d", offset) }

// v1LatestCursor 读取最新持久化提交，作为当前查询快照边界。
func v1LatestCursor(ctx context.Context) (string, error) {
	commit, err := utils.KidsDB(ctx).Model(consts.KidsV1CommitTable).Ctx(ctx).OrderDesc("commit_sequence").One()
	if err != nil {
		return "", err
	}
	if commit.IsEmpty() {
		return v1CommitCursor(0), nil
	}
	return v1CommitCursor(commit["commit_sequence"].Int64()), nil
}

// v1LatestCursorTx 在当前数据库事务的同一读取快照中获取最新提交游标。
func v1LatestCursorTx(ctx context.Context, tx gdb.TX) (string, error) {
	commit, err := tx.Model(consts.KidsV1CommitTable).Ctx(ctx).OrderDesc("commit_sequence").One()
	if err != nil {
		return "", err
	}
	if commit.IsEmpty() {
		return v1CommitCursor(0), nil
	}
	return v1CommitCursor(commit["commit_sequence"].Int64()), nil
}

// v1CircleRecordProjection 将接口圈子数据库记录还原为冻结 wire 投影。
func v1CircleRecordProjection(record gdb.Record) map[string]any {
	return map[string]any{"circle_id": record["circle_id"].String(), "name": record["name"].String(), "icon": v1JSONValue(record["icon"].String()), "owner_administrator_id": record["owner_administrator_id"].String(), "status": record["status"].String(), "version": record["version"].Int64(), "created_at_ms": record["created_at"].Time().UnixMilli(), "updated_at_ms": record["updated_at"].Time().UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(record["deleted_at"].Time())}
}

// v1MembershipRecordProjection 将接口成员身份数据库记录还原为冻结 wire 投影。
func v1MembershipRecordProjection(record gdb.Record) map[string]any {
	return map[string]any{"membership_id": record["membership_id"].String(), "circle_id": record["circle_id"].String(), "account_id": record["account_id"].String(), "actor_type": record["actor_type"].String(), "actor_id": record["actor_id"].String(), "role": record["role"].String(), "permissions": v1JSONValue(record["permissions"].String()), "status": record["status"].String(), "version": record["version"].Int64(), "created_at_ms": record["created_at"].Time().UnixMilli(), "updated_at_ms": record["updated_at"].Time().UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(record["deleted_at"].Time())}
}

// v1JSONValue 解析数据库 JSON 列，解析失败时返回 nil 使响应校验安全拒绝。
func v1JSONValue(value string) any {
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return nil
	}
	return decoded
}

// v1NullableTimeMillis 将可空数据库时间转换为接口 epoch 毫秒字段。
func v1NullableTimeMillis(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value.UnixMilli()
}

// v1CompleteOnboarding 原子创建圈子、所有者、初始成员、授权身份和当前选择。
func (s *sKids) v1CompleteOnboarding(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleInput, ok := in.Body["circle"].(map[string]any)
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "circle input is invalid")
	}
	ownerInput, ok := in.Body["owner_administrator"].(map[string]any)
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "owner administrator input is invalid")
	}
	membersInput, ok := in.Body["initial_members"].([]any)
	if !ok || len(membersInput) == 0 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "initial members input is invalid")
	}
	circleID := fmt.Sprint(circleInput["circle_id"])
	if circleID == "" || circleID != fmt.Sprint(in.Body["selected_circle_id"]) {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "selected circle must be the created circle")
	}
	ownerID := fmt.Sprint(ownerInput["administrator_id"])
	if ownerID == "" {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "owner administrator ID is invalid")
	}
	memberIDs := make(map[string]struct{}, len(membersInput))
	for _, rawMember := range membersInput {
		member, ok := rawMember.(map[string]any)
		if !ok || fmt.Sprint(member["member_id"]) == "" {
			return nil, "", v1Error(422, "VALIDATION_FAILED", false, "initial member is invalid")
		}
		memberID := fmt.Sprint(member["member_id"])
		if _, duplicated := memberIDs[memberID]; duplicated {
			return nil, "", v1Error(422, "VALIDATION_FAILED", false, "initial member IDs must be unique")
		}
		memberIDs[memberID] = struct{}{}
	}

	now := time.Now()
	permissions := v1OwnerPermissions()
	selectionID := v1ID("selection", uuid.NewString())
	membershipID := v1ID("membership", uuid.NewString())
	iconJSON, err := json.Marshal(circleInput["icon"])
	if err != nil {
		return nil, "", err
	}
	avatarJSON, err := json.Marshal(ownerInput["avatar"])
	if err != nil {
		return nil, "", err
	}
	permissionsJSON, err := json.Marshal(permissions)
	if err != nil {
		return nil, "", err
	}

	circle := v1CircleProjection(circleID, circleInput, ownerID, now)
	owner := v1AdministratorProjection(circleID, ownerID, in.PrincipalID, ownerInput, permissions, now)
	membership := v1MembershipProjection(membershipID, circleID, in.PrincipalID, ownerID, permissions, now)
	selection := v1SelectionProjection(selectionID, in.PrincipalID, circleID, 1, now)
	members := make([]map[string]any, 0, len(membersInput))
	for _, rawMember := range membersInput {
		member := rawMember.(map[string]any)
		members = append(members, v1MemberProjection(circleID, member, now))
	}

	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		existingCircle, queryErr := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", circleID).LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if !existingCircle.IsEmpty() {
			return v1Error(422, "VALIDATION_FAILED", false, "circle already exists")
		}
		if _, insertErr := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Data(gdb.Map{
			"circle_id": circleID, "name": circleInput["name"], "icon": string(iconJSON), "owner_administrator_id": ownerID,
			"status": "active", "version": 1, "created_at": now, "updated_at": now,
		}).Insert(); insertErr != nil {
			return insertErr
		}
		if _, insertErr := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Data(gdb.Map{
			"administrator_id": ownerID, "circle_id": circleID, "bound_account_id": in.PrincipalID,
			"display_name": ownerInput["display_name"], "avatar": string(avatarJSON), "role": "owner", "permissions": string(permissionsJSON),
			"status": "active", "version": 1, "created_at": now, "updated_at": now,
		}).Insert(); insertErr != nil {
			return insertErr
		}
		for _, rawMember := range membersInput {
			member := rawMember.(map[string]any)
			memberAvatar, marshalErr := json.Marshal(member["avatar"])
			if marshalErr != nil {
				return marshalErr
			}
			if _, insertErr := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Data(gdb.Map{
				"member_id": member["member_id"], "circle_id": circleID, "display_name": member["display_name"], "gender": member["gender"],
				"avatar": string(memberAvatar), "status": "active", "version": 1, "created_at": now, "updated_at": now,
			}).Insert(); insertErr != nil {
				return insertErr
			}
		}
		if _, insertErr := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Data(gdb.Map{
			"membership_id": membershipID, "circle_id": circleID, "account_id": in.PrincipalID, "actor_type": "owner", "actor_id": ownerID,
			"role": "owner", "permissions": string(permissionsJSON), "status": "active", "version": 1, "created_at": now, "updated_at": now,
		}).Insert(); insertErr != nil {
			return insertErr
		}
		existingSelection, queryErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if existingSelection.IsEmpty() {
			if _, insertErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Data(gdb.Map{
				"selection_id": selectionID, "account_id": in.PrincipalID, "current_circle_id": circleID, "version": 1, "created_at": now, "updated_at": now,
			}).Insert(); insertErr != nil {
				return insertErr
			}
		} else {
			selection["selection_id"] = existingSelection["selection_id"].String()
			selection["version"] = existingSelection["version"].Int64() + 1
			if _, updateErr := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).Data(gdb.Map{
				"current_circle_id": circleID, "version": selection["version"], "updated_at": now,
			}).Update(); updateErr != nil {
				return updateErr
			}
		}
		var commitErr error
		receipt, commitErr = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{
			"circle": circle, "owner_administrator": owner, "initial_members": members, "membership": membership, "selection": selection,
		})
		if commitErr != nil {
			return commitErr
		}
		for _, member := range members {
			if _, insertErr := tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Data(gdb.Map{
				"circle_id": circleID, "member_id": member["member_id"], "balance": 0, "version": 1,
				"source_commit_id": receipt["commit_id"], "source_commit_sequence": receipt["commit_sequence"], "updated_at": now,
			}).Insert(); insertErr != nil {
				return insertErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	syncCursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{
		"receipt": receipt, "circle": circle, "owner_administrator": owner, "initial_members": members,
		"membership": membership, "selection": selection, "sync_cursor": syncCursor,
	}, syncCursor, nil
}

// v1OwnerPermissions 返回 Owner 在接口中拥有的固定圈子管理权限集合。
func v1OwnerPermissions() []string {
	return []string{"manage_circle", "manage_members", "manage_tasks", "manage_rewards", "adjust_stars"}
}

// v1CircleProjection 构造新建圈子的接口投影。
func v1CircleProjection(circleID string, input map[string]any, ownerID string, now time.Time) map[string]any {
	return map[string]any{"circle_id": circleID, "name": input["name"], "icon": input["icon"], "owner_administrator_id": ownerID, "status": "active", "version": int64(1), "created_at_ms": now.UnixMilli(), "updated_at_ms": now.UnixMilli(), "deleted_at_ms": nil}
}

// v1AdministratorProjection 构造新建所有者管理员的接口投影。
func v1AdministratorProjection(circleID, administratorID, accountID string, input map[string]any, permissions []string, now time.Time) map[string]any {
	return map[string]any{"administrator_id": administratorID, "circle_id": circleID, "bound_account_id": accountID, "display_name": input["display_name"], "avatar": input["avatar"], "role": "owner", "permissions": permissions, "status": "active", "version": int64(1), "created_at_ms": now.UnixMilli(), "updated_at_ms": now.UnixMilli(), "deleted_at_ms": nil}
}

// v1MemberProjection 构造 onboarding 初始成员的接口投影。
func v1MemberProjection(circleID string, input map[string]any, now time.Time) map[string]any {
	return map[string]any{"member_id": input["member_id"], "circle_id": circleID, "bound_account_id": nil, "display_name": input["display_name"], "gender": input["gender"], "avatar": input["avatar"], "status": "active", "version": int64(1), "created_at_ms": now.UnixMilli(), "updated_at_ms": now.UnixMilli(), "deleted_at_ms": nil}
}

// v1MembershipProjection 构造 Owner 对当前圈子的授权身份投影。
func v1MembershipProjection(membershipID, circleID, accountID, actorID string, permissions []string, now time.Time) map[string]any {
	return map[string]any{"membership_id": membershipID, "circle_id": circleID, "account_id": accountID, "actor_type": "owner", "actor_id": actorID, "role": "owner", "permissions": permissions, "status": "active", "version": int64(1), "created_at_ms": now.UnixMilli(), "updated_at_ms": now.UnixMilli(), "deleted_at_ms": nil}
}

// v1SelectionProjection 构造账号的当前圈子选择投影。
func v1SelectionProjection(selectionID, accountID, circleID string, version int64, now time.Time) map[string]any {
	return map[string]any{"selection_id": selectionID, "account_id": accountID, "current_circle_id": nullableString(circleID), "version": version, "updated_at_ms": now.UnixMilli()}
}

// v1GetCurrentAccount 根据持久化 principal 返回严格区分的账号或游客 bootstrap。
func (s *sKids) v1GetCurrentAccount(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	session, err := utils.KidsDB(ctx).Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", in.SessionID).One()
	if err != nil {
		return nil, "", err
	}
	if session.IsEmpty() {
		return nil, "", v1Error(401, "UNAUTHENTICATED", false, "guest session is missing")
	}
	if in.PrincipalKind == "account" {
		account, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1AccountTable).Ctx(ctx).Where("account_id", in.PrincipalID).One()
		if queryErr != nil {
			return nil, "", queryErr
		}
		binding, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1AccountBindingTable).Ctx(ctx).Where("account_id", in.PrincipalID).One()
		if queryErr != nil {
			return nil, "", queryErr
		}
		if account.IsEmpty() || binding.IsEmpty() {
			return nil, "", v1Error(503, "UNAVAILABLE", true, "account binding is unavailable")
		}
		memberships, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("account_id", in.PrincipalID).Where("status", "active").Order("created_at ASC,id ASC").All()
		if queryErr != nil {
			return nil, "", queryErr
		}
		membershipOut := make([]map[string]any, 0, len(memberships))
		for _, row := range memberships {
			membershipOut = append(membershipOut, v1MembershipRecordProjection(row))
		}
		selection, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).One()
		if queryErr != nil {
			return nil, "", queryErr
		}
		cursor, queryErr := v1LatestCursor(ctx)
		if queryErr != nil {
			return nil, "", queryErr
		}
		return map[string]any{"account": v1AccountRecordProjection(account), "session": v1SessionMetadataProjection(session), "account_binding": v1AccountBindingRecordProjection(binding), "principal_kind": "account", "memberships": membershipOut, "selection": v1SelectionRecordProjection(selection), "entitlement": nil, "bootstrap_cursor": cursor}, cursor, nil
	}
	if in.PrincipalKind != "invite_guest" {
		return nil, "", v1Error(403, "FORBIDDEN", false, "unsupported session principal")
	}
	expiresAtMs := session["access_expires_at_ms"].Int64()
	return map[string]any{
		"principal_kind":      "invite_guest",
		"session":             v1SessionMetadataProjection(session),
		"capabilities":        []string{"invite_redeem_administrator", "invite_redeem_member", "current_account_bootstrap", "submit_feedback", "feedback_asset_upload"},
		"guest_expires_at_ms": expiresAtMs,
		"bootstrap_cursor":    nil,
	}, "", nil
}

// v1RefreshSession 原子轮换账号 session 的 access 和 refresh credential，只保存其摘要。
func (s *sKids) v1RefreshSession(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	issuedAtMs := time.Now().UnixMilli()
	accessExpiryMs := issuedAtMs + int64(time.Hour/time.Millisecond)
	refreshExpiryMs := issuedAtMs + int64((30*24*time.Hour)/time.Millisecond)
	accessToken, err := v1IssueAccessToken(ctx, in.SessionID, in.PrincipalID, "account", issuedAtMs, accessExpiryMs)
	if err != nil {
		return nil, "", err
	}
	refreshToken := v1Secret()
	receiptID := v1ID("receipt", uuid.NewString())
	var sessionRecord gdb.Record
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		session, err := tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", in.SessionID).Where("account_id", in.PrincipalID).Where("principal_kind", "account").Where("refresh_token_hash", sha256Hex(in.AccessToken)).Where("status", "active").LockUpdate().One()
		if err != nil {
			return err
		}
		if session.IsEmpty() {
			return v1Error(401, "UNAUTHENTICATED", false, "account session is missing")
		}
		if _, err = tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", in.SessionID).Data(gdb.Map{
			"access_token_hash": sha256Hex(accessToken), "refresh_token_hash": sha256Hex(refreshToken), "issued_at_ms": issuedAtMs,
			"access_expires_at_ms": accessExpiryMs, "refresh_expires_at_ms": refreshExpiryMs, "status": "active", "revoked_at_ms": nil,
		}).Update(); err != nil {
			return err
		}
		sessionRecord, err = tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", in.SessionID).One()
		if err != nil {
			return err
		}
		if sessionRecord.IsEmpty() || !v1SessionIsActive(sessionRecord) {
			return v1Error(503, "UNAVAILABLE", true, "rotated account session is unavailable")
		}
		if _, err = tx.Model(consts.KidsV1ReceiptTable).Ctx(ctx).Data(gdb.Map{"receipt_id": receiptID, "commit_id": "", "operation_id": in.OperationID, "result_kind": "first_committed", "committed_at": time.UnixMilli(issuedAtMs)}).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"session": v1AuthSessionProjection(sessionRecord, accessToken, refreshToken), "rotation_receipt_id": receiptID}, "", nil
}

// v1RevokeSession 撤销当前账号拥有的指定 session，并对重复撤销返回稳定结果。
func (s *sKids) v1RevokeSession(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	sessionID := in.PathParameters["session_id"]
	now := time.Now()
	receiptID := v1ID("receipt", uuid.NewString())
	alreadyRevoked := false
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		row, err := tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", sessionID).Where("account_id", in.PrincipalID).LockUpdate().One()
		if err != nil {
			return err
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "session is missing")
		}
		if row["status"].String() == "revoked" {
			alreadyRevoked = true
			now = time.UnixMilli(row["revoked_at_ms"].Int64())
			return nil
		}
		if _, err = tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", sessionID).Data(gdb.Map{"status": "revoked", "revoked_at_ms": now.UnixMilli()}).Update(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1ReceiptTable).Ctx(ctx).Data(gdb.Map{"receipt_id": receiptID, "commit_id": "", "operation_id": in.OperationID, "result_kind": "first_committed", "committed_at": now}).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"session_id": sessionID, "revoked_at_ms": now.UnixMilli(), "already_revoked": alreadyRevoked, "receipt_id": receiptID}, "", nil
}

// v1AccountRecordProjection 还原接口账号投影。
func v1AccountRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"account_id": row["account_id"].String(), "status": row["status"].String(), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "updated_at_ms": row["updated_at"].Time().UnixMilli()}
}

// v1AccountBindingRecordProjection 还原接口账号绑定投影。
func v1AccountBindingRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"binding_id": row["binding_id"].String(), "account_id": row["account_id"].String(), "environment": row["environment"].String(), "migration_policy": row["migration_policy"].String(), "version": row["version"].Int64(), "issued_at_ms": row["issued_at"].Time().UnixMilli()}
}

// v1SessionEnvironment 返回写入 token 并参与验签隔离的当前部署环境。
func v1SessionEnvironment(ctx context.Context) string {
	environment := strings.TrimSpace(g.Cfg().MustGet(ctx, "app.env", "").String())
	if environment == "" {
		return "dev"
	}
	return environment
}

// v1SessionSigningSecret 返回 v1 access token 的签发和验签密钥。
func v1SessionSigningSecret(ctx context.Context) string {
	return commonlogin.SecretFromConfig(ctx)
}

// v1IssueAccessToken 签发与 session 行一一对应的毫秒级 access token。
func v1IssueAccessToken(ctx context.Context, sessionID, accountID, principalKind string, issuedAtMs, expiresAtMs int64) (string, error) {
	return utils.GenerateV1SessionToken(utils.V1SessionTokenClaims{
		SessionID: sessionID, AccountID: accountID, PrincipalKind: principalKind, Environment: v1SessionEnvironment(ctx), IssuedAtMs: issuedAtMs, ExpiresAtMs: expiresAtMs,
	}, v1SessionSigningSecret(ctx))
}

// v1SessionIsActive 判断数据库 session 状态是否仍可用于授权。
func v1SessionIsActive(row gdb.Record) bool {
	return row["status"].String() == "active"
}

// v1SessionMetadataProjection 还原不含 secret 的 session metadata。
func v1SessionMetadataProjection(row gdb.Record) map[string]any {
	return map[string]any{
		"session_id": row["session_id"].String(), "principal_type": row["principal_kind"].String(), "status": row["status"].String(),
		"issued_at_ms": row["issued_at_ms"].Int64(), "access_expires_at_ms": row["access_expires_at_ms"].Int64(), "refresh_expires_at_ms": nullableV1Int64(row["refresh_expires_at_ms"].Int64()),
	}
}

// nullableV1Int64 将可空 epoch milliseconds 映射为 wire 的 null。
func nullableV1Int64(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
}

// v1AuthSessionProjection 构造只在 token 签发或刷新时返回的 account AuthSession。
func v1AuthSessionProjection(session gdb.Record, accessToken, refreshToken string) map[string]any {
	return map[string]any{"token_type": "Bearer", "access_token": accessToken, "refresh_token": refreshToken, "metadata": v1SessionMetadataProjection(session)}
}

// v1SelectionRecordProjection 将不存在的 selection 表示为 JSON null。
func v1SelectionRecordProjection(row gdb.Record) any {
	if row.IsEmpty() {
		return nil
	}
	return v1SelectionProjection(row["selection_id"].String(), row["account_id"].String(), row["current_circle_id"].String(), row["version"].Int64(), row["updated_at"].Time())
}

// v1CreateInvite 创建绑定到既有管理员或成员的单次邀请码，并只在本响应中返回明文邀请码。
func (s *sKids) v1CreateInvite(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	ttl, ok := v1ExpectedVersion(in.Body["requested_ttl_seconds"])
	if !ok || ttl < 300 || ttl > 2592000 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "invite TTL is invalid")
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_target_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "target version is invalid")
	}
	role := fmt.Sprint(in.Body["target_role"])
	inviteID := fmt.Sprint(in.Body["invite_id"])
	adminID := nullableV1String(in.Body["target_administrator_id"])
	memberID := nullableV1String(in.Body["target_member_id"])
	scope, err := v1StringSlice(in.Body["permission_scope"])
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "permission scope is invalid")
	}
	if inviteID == "" || (role == "administrator" && (adminID == "" || memberID != "" || len(scope) == 0)) || (role == "member" && (memberID == "" || adminID != "" || len(scope) != 0)) {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "invite target is invalid")
	}
	now := time.Now()
	expiresAt := now.Add(time.Duration(ttl) * time.Second)
	code, err := v1InviteCode()
	if err != nil {
		return nil, "", err
	}
	scopeJSON, err := json.Marshal(scope)
	if err != nil {
		return nil, "", err
	}
	var invite map[string]any
	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err = v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
			return err
		}
		if role == "administrator" {
			row, e := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", adminID).Where("circle_id", circleID).Where("status", "pending_invite").LockUpdate().One()
			if e != nil {
				return e
			}
			if row.IsEmpty() {
				return v1Error(404, "NOT_FOUND", false, "invite administrator target is missing")
			}
			if row["version"].Int64() != expected {
				v := row["version"].Int64()
				return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &v, Message: "target version conflicts"}
			}
		} else {
			row, e := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
			if e != nil {
				return e
			}
			if row.IsEmpty() {
				return v1Error(404, "NOT_FOUND", false, "invite member target is missing")
			}
			if row["version"].Int64() != expected {
				v := row["version"].Int64()
				return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &v, Message: "target version conflicts"}
			}
		}
		existing, e := tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).One()
		if e != nil {
			return e
		}
		if !existing.IsEmpty() {
			return v1Error(409, "IDEMPOTENCY_CONFLICT", false, "invite ID already exists")
		}
		if _, e = tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Data(gdb.Map{"invite_id": inviteID, "circle_id": circleID, "target_role": role, "target_administrator_id": adminID, "target_member_id": memberID, "permission_scope": string(scopeJSON), "code_hash": sha256Hex(code), "single_use": true, "generation": 1, "status": "active", "expires_at": expiresAt, "version": 1, "created_by_account_id": in.PrincipalID, "created_at": now, "updated_at": now}).Insert(); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).One()
		if e != nil {
			return e
		}
		invite = v1InviteRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"invite": invite})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "invite": invite, "invite_code": code}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1RefreshInvite 作废旧邀请码摘要、递增版本与 generation，并返回一次新的邀请码。
func (s *sKids) v1RefreshInvite(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	inviteID := in.PathParameters["invite_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	ttl, ok := v1ExpectedVersion(in.Body["requested_ttl_seconds"])
	if !ok || ttl < 300 || ttl > 2592000 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "invite TTL is invalid")
	}
	now := time.Now()
	code, err := v1InviteCode()
	if err != nil {
		return nil, "", err
	}
	var invite map[string]any
	var receipt map[string]any
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).Where("circle_id", circleID).LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "invite is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "invite version conflicts"}
		}
		if row["status"].String() != "active" {
			return v1Error(409, "VERSION_CONFLICT", false, "invite is not active")
		}
		version := current + 1
		generation := row["generation"].Int64() + 1
		if _, e = tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).Data(gdb.Map{"code_hash": sha256Hex(code), "expires_at": now.Add(time.Duration(ttl) * time.Second), "generation": generation, "version": version, "updated_at": now}).Update(); e != nil {
			return e
		}
		row, e = tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).One()
		if e != nil {
			return e
		}
		invite = v1InviteRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"invite": invite})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "invite": invite, "invite_code": code}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1RevokeInvite 软撤销 active invite，撤销响应不会重新暴露邀请码秘密。
func (s *sKids) v1RevokeInvite(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	inviteID := in.PathParameters["invite_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected version is invalid")
	}
	now := time.Now()
	var invite map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageMembers); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).Where("circle_id", circleID).LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "invite is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "invite version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).Data(gdb.Map{"status": "revoked", "revoked_at": now, "version": version, "updated_at": now}).Update(); e != nil {
			return e
		}
		row, e = tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", inviteID).One()
		if e != nil {
			return e
		}
		invite = v1InviteRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"invite": invite})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "invite": invite}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1InviteRecordProjection 将邀请码事实投影为不含 code secret 的接口 outgoing invite。
func v1InviteRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"invite_id": row["invite_id"].String(), "circle_id": row["circle_id"].String(), "target_role": row["target_role"].String(), "target_administrator_id": nullableString(row["target_administrator_id"].String()), "target_member_id": nullableString(row["target_member_id"].String()), "permission_scope": v1JSONValue(row["permission_scope"].String()), "single_use": row["single_use"].Bool(), "generation": row["generation"].Int64(), "status": row["status"].String(), "expires_at_ms": row["expires_at"].Time().UnixMilli(), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "revoked_at_ms": v1NullableTimeMillis(row["revoked_at"].Time()), "used_at_ms": v1NullableTimeMillis(row["used_at"].Time())}
}

// v1StringSlice 把已通过 JSON schema 的字符串数组转换成便于数据库持久化的切片。
func v1StringSlice(value any) ([]string, error) {
	if encoded, ok := value.(string); ok {
		value = v1JSONValue(encoded)
	}
	if encoded, ok := value.([]byte); ok {
		value = v1JSONValue(string(encoded))
	}
	if values, ok := value.([]string); ok {
		return append([]string(nil), values...), nil
	}
	values, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("not array")
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		item, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("not string")
		}
		out = append(out, item)
	}
	return out, nil
}

// v1InviteCode 使用无歧义大写字母和数字生成八位敏感邀请码。
func v1InviteCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	raw := make([]byte, 8)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	for i := range raw {
		raw[i] = alphabet[int(raw[i])%len(alphabet)]
	}
	return string(raw), nil
}

// v1RedeemInvite 原子核销邀请码、绑定 account 到目标 actor，并创建该圈子的 active membership。
func (s *sKids) v1RedeemInvite(ctx context.Context, in v1.V1OperationInput, targetRole string) (map[string]any, string, error) {
	if in.PrincipalKind != "account" {
		return nil, "", v1Error(503, "UNAVAILABLE", true, "guest account upgrade requires verified Google proof")
	}
	code := fmt.Sprint(in.Body["invite_code"])
	now := time.Now()
	var receipt map[string]any
	var circle map[string]any
	var actor map[string]any
	var membership map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		invite, err := tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("code_hash", sha256Hex(code)).LockUpdate().One()
		if err != nil {
			return err
		}
		if invite.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "invite code is invalid")
		}
		if invite["target_role"].String() != targetRole {
			return v1Error(403, "FORBIDDEN", false, "invite role does not match redemption")
		}
		if invite["status"].String() != "active" {
			return v1Error(409, "VERSION_CONFLICT", false, "invite is no longer active")
		}
		if !invite["expires_at"].Time().After(now) {
			return v1Error(409, "VERSION_CONFLICT", false, "invite has expired")
		}
		circleRow, err := tx.Model(consts.KidsV1CircleTable).Ctx(ctx).Where("circle_id", invite["circle_id"].String()).Where("status", "active").LockUpdate().One()
		if err != nil {
			return err
		}
		if circleRow.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "invite circle is missing")
		}
		existingMembership, err := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("account_id", in.PrincipalID).Where("circle_id", invite["circle_id"].String()).LockUpdate().One()
		if err != nil {
			return err
		}
		if !existingMembership.IsEmpty() {
			return v1Error(409, "VERSION_CONFLICT", false, "account already has circle membership")
		}
		permissions := []string{}
		actorID := ""
		if targetRole == "administrator" {
			permissions, err = v1StringSlice(invite["permission_scope"].Interface())
			if err != nil {
				return err
			}
			actorID = invite["target_administrator_id"].String()
			row, e := tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", actorID).Where("circle_id", invite["circle_id"].String()).Where("status", "pending_invite").LockUpdate().One()
			if e != nil {
				return e
			}
			if row.IsEmpty() {
				return v1Error(404, "NOT_FOUND", false, "administrator target is missing")
			}
			version := row["version"].Int64() + 1
			if _, e = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", actorID).Data(gdb.Map{"bound_account_id": in.PrincipalID, "permissions": mustV1JSON(permissions), "status": "active", "version": version, "updated_at": now}).Update(); e != nil {
				return e
			}
			row, e = tx.Model(consts.KidsV1AdministratorTable).Ctx(ctx).Where("administrator_id", actorID).One()
			if e != nil {
				return e
			}
			actor = v1AdministratorRecordProjection(row)
		} else {
			actorID = invite["target_member_id"].String()
			row, e := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", actorID).Where("circle_id", invite["circle_id"].String()).Where("status", "active").LockUpdate().One()
			if e != nil {
				return e
			}
			if row.IsEmpty() {
				return v1Error(404, "NOT_FOUND", false, "member target is missing")
			}
			if row["bound_account_id"].String() != "" {
				return v1Error(409, "VERSION_CONFLICT", false, "member is already bound")
			}
			version := row["version"].Int64() + 1
			if _, e = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", actorID).Data(gdb.Map{"bound_account_id": in.PrincipalID, "version": version, "updated_at": now}).Update(); e != nil {
				return e
			}
			row, e = tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", actorID).One()
			if e != nil {
				return e
			}
			actor = v1MemberRecordProjection(row)
		}
		membershipID := v1ID("membership", uuid.NewString())
		role := targetRole
		if _, err = tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Data(gdb.Map{"membership_id": membershipID, "circle_id": invite["circle_id"].String(), "account_id": in.PrincipalID, "actor_type": role, "actor_id": actorID, "role": role, "permissions": mustV1JSON(permissions), "status": "active", "version": 1, "created_at": now, "updated_at": now}).Insert(); err != nil {
			return err
		}
		membership = v1MembershipProjectionForActor(membershipID, invite["circle_id"].String(), in.PrincipalID, role, actorID, permissions, now)
		if _, err = tx.Model(consts.KidsV1InviteTable).Ctx(ctx).Where("invite_id", invite["invite_id"].String()).Data(gdb.Map{"status": "used", "used_at": now, "version": invite["version"].Int64() + 1, "updated_at": now}).Update(); err != nil {
			return err
		}
		selection, e := tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Where("account_id", in.PrincipalID).LockUpdate().One()
		if e != nil {
			return e
		}
		if selection.IsEmpty() {
			if _, e = tx.Model(consts.KidsV1CircleSelectionTable).Ctx(ctx).Data(gdb.Map{"selection_id": v1ID("selection", uuid.NewString()), "account_id": in.PrincipalID, "current_circle_id": invite["circle_id"].String(), "version": 1, "created_at": now, "updated_at": now}).Insert(); e != nil {
				return e
			}
		}
		circle = v1CircleRecordProjection(circleRow)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, invite["circle_id"].String(), in.OperationID, map[string]any{"circle": circle, "membership": membership, "actor": actor})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	data := map[string]any{"auth_session": nil, "membership": membership, "circle": circle, "redemption_receipt": receipt, "sync_cursor": cursor}
	if targetRole == "administrator" {
		data["administrator"] = actor
	} else {
		data["member"] = actor
	}
	return data, cursor, nil
}

// v1MembershipProjectionForActor 构造管理员或成员加入后的接口 membership 投影。
func v1MembershipProjectionForActor(membershipID, circleID, accountID, role, actorID string, permissions []string, now time.Time) map[string]any {
	return map[string]any{"membership_id": membershipID, "circle_id": circleID, "account_id": accountID, "actor_type": role, "actor_id": actorID, "role": role, "permissions": permissions, "status": "active", "version": int64(1), "created_at_ms": now.UnixMilli(), "updated_at_ms": now.UnixMilli(), "deleted_at_ms": nil}
}

// mustV1JSON 序列化内部生成的简单结构，失败时返回空 JSON 数组以保证数据库 JSON 合法。
func mustV1JSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

// v1ListCircles 从当前 kids 圈子投影构造接口分页数据。
func (s *sKids) v1ListCircles(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	userID := v1UserID(ctx)
	if userID == 0 {
		return nil, "", v1Error(401, "UNAUTHENTICATED", false, "account credential is required")
	}
	out, err := s.ListCircles(ctx, v1.CircleListInput{UserId: userID})
	if err != nil {
		return nil, "", err
	}
	items := make([]map[string]any, 0, len(out.Managed)+len(out.Joined))
	for _, circle := range append(out.Managed, out.Joined...) {
		items = append(items, map[string]any{
			"circle_id": v1ID("circle", circle.Id), "name": circle.Name,
			"icon": nullableString(circle.Icon), "role": circle.Role, "status": "active", "version": 1,
		})
	}
	return map[string]any{"items": items, "next_cursor": nil, "has_more": false, "snapshot_cursor": v1Cursor()}, "", nil
}

// v1Statistics 返回符合接口时间和分页语义的统计骨架。
func (s *sKids) v1Statistics(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	if in.OperationID == "compareStatistics" {
		return map[string]any{"base_member": map[string]any{}, "compare_member": map[string]any{}, "snapshot_cursor": v1Cursor()}, "", nil
	}
	return map[string]any{"member_id": v1QueryFirst(in.Query, "member_id"), "metric": v1QueryFirst(in.Query, "metric"), "period_type": v1QueryFirst(in.Query, "period_type"), "bucket_unit": v1QueryFirst(in.Query, "bucket_unit"), "start_at_ms": v1QueryFirst(in.Query, "start_at_ms"), "end_at_ms": v1QueryFirst(in.Query, "end_at_ms"), "zone_id": v1QueryFirst(in.Query, "zone_id"), "buckets": []any{}, "summary": map[string]any{"total": 0, "peak_value": 0, "non_zero_bucket_count": 0}, "as_of_cursor": v1Cursor()}, "", nil
}

// v1Entitlement 从接口权益快照读取当前权益。
func (s *sKids) v1Entitlement(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	return map[string]any{"account_id": v1AccountID(ctx, in), "entitlement_id": nil, "status": "none", "plan_code": nil, "base_plan_id": nil, "offer_id": nil, "product_id": nil, "billing_phase": nil, "trial_eligible": false, "valid_until_ms": nil, "verified_at_ms": nil, "revoked_at_ms": nil, "version": 1}, "", nil
}

// v1SubmitFeedback 持久化结构化反馈并创建接口回执。
func (s *sKids) v1SubmitFeedback(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	feedbackID := v1ID("feedback", uuid.NewString())
	content := fmt.Sprint(in.Body["content"])
	if strings.TrimSpace(content) == "" {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "feedback content is required")
	}
	attachmentJSON, _ := json.Marshal(in.Body["attachment_asset_ids"])
	now := time.Now()
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(consts.KidsV1FeedbackTable).Ctx(ctx).Data(gdb.Map{
			"feedback_id": feedbackID, "account_id": v1AccountID(ctx, in), "category": in.Body["category"], "content": content,
			"contact_type": nullableV1String(in.Body["contact_type"]), "contact": nullableV1String(in.Body["contact"]),
			"privacy_consent_version": in.Body["privacy_consent_version"], "attachment_asset_ids": string(attachmentJSON),
		}).Insert(); err != nil {
			return err
		}
		var commitErr error
		receipt, commitErr = v1CreateCommitTx(ctx, tx, "", in.OperationID, map[string]any{"feedback_ids": []string{feedbackID}})
		return commitErr
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"feedback_id": feedbackID, "status": "accepted", "received_at_ms": now.UnixMilli(), "version": 1, "receipt": receipt}, "", nil
}

// v1Mutation 统一生成持久化接口回执，领域迁移期间保持稳定 wire 结构。
func (s *sKids) v1Mutation(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	receipt := v1Receipt(in.OperationID)
	data := map[string]any{"receipt": receipt}
	changeCursor := v1Cursor()
	if in.OperationID == "createInviteGuestSession" {
		return s.v1CreateGuestSession(ctx, in)
	}
	if in.OperationID == "exchangeGoogleProof" {
		return s.v1ExchangeGoogleProof(ctx, in)
	}
	if in.OperationID == "refreshSession" {
		return s.v1RefreshSession(ctx, in)
	}
	if in.OperationID == "revokeSession" {
		return s.v1RevokeSession(ctx, in)
	}
	if in.OperationID == "prepareAssetUpload" {
		return s.v1PrepareAsset(ctx, in)
	}
	if in.OperationID == "commitAssetUpload" {
		return s.v1CommitAsset(ctx, in)
	}
	data["change_cursor"] = changeCursor
	return data, changeCursor, nil
}

// v1CreateGuestSession 创建接口游客会话并把 token 摘要持久化。
func (s *sKids) v1CreateGuestSession(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	issuedAtMs := time.Now().UnixMilli()
	sessionID := v1ID("session", uuid.NewString())
	accessExpiryMs := issuedAtMs + int64(time.Hour/time.Millisecond)
	accessToken, err := v1IssueAccessToken(ctx, sessionID, "", "invite_guest", issuedAtMs, accessExpiryMs)
	if err != nil {
		return nil, "", err
	}
	refreshToken := v1Secret()
	grant := v1Secret()
	refreshExpiryMs := issuedAtMs + int64((24*time.Hour)/time.Millisecond)
	_, err = utils.KidsDB(ctx).Model(consts.KidsV1SessionTable).Ctx(ctx).Data(gdb.Map{
		"session_id": sessionID, "principal_kind": "invite_guest", "status": "active", "access_token_hash": sha256Hex(accessToken), "refresh_token_hash": sha256Hex(refreshToken), "guest_upgrade_grant_hash": sha256Hex(grant), "issued_at_ms": issuedAtMs, "access_expires_at_ms": accessExpiryMs, "refresh_expires_at_ms": refreshExpiryMs,
	}).Insert()
	if err != nil {
		return nil, "", err
	}
	session, err := utils.KidsDB(ctx).Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", sessionID).One()
	if err != nil {
		return nil, "", err
	}
	if session.IsEmpty() {
		return nil, "", v1Error(503, "UNAVAILABLE", true, "guest session is unavailable")
	}
	return map[string]any{
		"session": map[string]any{
			"token_type":    "Bearer",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"metadata":      v1SessionMetadataProjection(session),
		},
		"capabilities":        []string{"invite_redeem_administrator", "invite_redeem_member", "current_account_bootstrap", "submit_feedback", "feedback_asset_upload"},
		"guest_expires_at_ms": session["access_expires_at_ms"].Int64(),
		"guest_upgrade_grant": grant,
	}, "", nil
}

// v1ExchangeGoogleProof 校验 Google 签名和 nonce，并在同一事务创建账号绑定与会话。
func (s *sKids) v1ExchangeGoogleProof(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	identityToken, nonce := strings.TrimSpace(fmt.Sprint(in.Body["google_id_token"])), strings.TrimSpace(fmt.Sprint(in.Body["proof_nonce"]))
	if identityToken == "" || nonce == "" {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "Google proof is required")
	}
	clientID := strings.TrimSpace(g.Cfg().MustGet(ctx, "auth.google.clientId", "").String())
	if clientID == "" {
		return nil, "", v1Error(503, "UNAVAILABLE", true, "Google issuer configuration is unavailable")
	}
	identity, err := utils.VerifyGoogleIdentityTokenWithNonce(ctx, identityToken, clientID, nonce)
	if err != nil {
		g.Log().Warningf(ctx, "Google 身份凭据校验失败 request_id=%s", in.RequestID)
		return nil, "", v1Error(401, "UNAUTHENTICATED", false, "Google proof is invalid")
	}
	accountID := v1ID("account", "google:"+identity.OpenId)
	issuedAtMs := time.Now().UnixMilli()
	accessExpiryMs := issuedAtMs + int64(time.Hour/time.Millisecond)
	refreshExpiryMs := issuedAtMs + int64((30*24*time.Hour)/time.Millisecond)
	sessionID := v1ID("session", uuid.NewString())
	accessToken, err := v1IssueAccessToken(ctx, sessionID, accountID, "account", issuedAtMs, accessExpiryMs)
	if err != nil {
		return nil, "", err
	}
	refreshToken := v1Secret()
	var account, binding map[string]any
	var memberships []map[string]any
	var persistedSession gdb.Record
	mergeOutcome := "new_scope"
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if grant := nullableV1String(in.Body["guest_upgrade_grant"]); grant != "" {
			if in.PrincipalKind != "invite_guest" || in.SessionID == "" {
				return v1Error(401, "UNAUTHENTICATED", false, "guest upgrade session is required")
			}
			guest, queryErr := tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", in.SessionID).Where("principal_kind", "invite_guest").LockUpdate().One()
			if queryErr != nil {
				return queryErr
			}
			if guest.IsEmpty() || guest["guest_upgrade_grant_hash"].String() != sha256Hex(grant) || !v1SessionIsActive(guest) || guest["access_expires_at_ms"].Int64() <= issuedAtMs {
				return v1Error(401, "UNAUTHENTICATED", false, "guest upgrade grant is invalid")
			}
			if _, queryErr = tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("id", guest["id"].Int64()).Data(gdb.Map{"status": "revoked", "revoked_at_ms": issuedAtMs}).Update(); queryErr != nil {
				return queryErr
			}
			mergeOutcome = "no_local_facts"
		}
		accountRow, queryErr := tx.Model(consts.KidsV1AccountTable).Ctx(ctx).Where("account_id", accountID).LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if accountRow.IsEmpty() {
			issuedAt := time.UnixMilli(issuedAtMs)
			if _, queryErr = tx.Model(consts.KidsV1AccountTable).Ctx(ctx).Data(gdb.Map{"account_id": accountID, "status": "active", "version": 1, "created_at": issuedAt, "updated_at": issuedAt}).Insert(); queryErr != nil {
				return queryErr
			}
			account = map[string]any{"account_id": accountID, "status": "active", "version": int64(1), "created_at_ms": issuedAtMs, "updated_at_ms": issuedAtMs}
		} else {
			if accountRow["status"].String() != "active" {
				return v1Error(403, "ACCOUNT_DISABLED", false, "account is disabled")
			}
			account = v1AccountRecordProjection(accountRow)
			mergeOutcome = "no_local_facts"
		}
		bindingRow, queryErr := tx.Model(consts.KidsV1AccountBindingTable).Ctx(ctx).Where("account_id", accountID).LockUpdate().One()
		if queryErr != nil {
			return queryErr
		}
		if bindingRow.IsEmpty() {
			bindingID := v1ID("binding", uuid.NewString())
			issuedAt := time.UnixMilli(issuedAtMs)
			if _, queryErr = tx.Model(consts.KidsV1AccountBindingTable).Ctx(ctx).Data(gdb.Map{"binding_id": bindingID, "account_id": accountID, "environment": consts.KidsV1AccountBindingEnvironmentLive, "migration_policy": "no_merge", "version": 1, "issued_at": issuedAt, "created_at": issuedAt, "updated_at": issuedAt}).Insert(); queryErr != nil {
				return queryErr
			}
			binding = map[string]any{"binding_id": bindingID, "account_id": accountID, "environment": consts.KidsV1AccountBindingEnvironmentLive, "migration_policy": "no_merge", "version": int64(1), "issued_at_ms": issuedAtMs}
		} else {
			if bindingRow["environment"].String() != consts.KidsV1AccountBindingEnvironmentLive {
				if _, queryErr = tx.Model(consts.KidsV1AccountBindingTable).Ctx(ctx).Where("id", bindingRow["id"].Int64()).Data(gdb.Map{"environment": consts.KidsV1AccountBindingEnvironmentLive, "updated_at": time.UnixMilli(issuedAtMs)}).Update(); queryErr != nil {
					return queryErr
				}
			}
			binding = v1AccountBindingRecordProjection(bindingRow)
			binding["environment"] = consts.KidsV1AccountBindingEnvironmentLive
		}
		if _, queryErr = tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Data(gdb.Map{"session_id": sessionID, "account_id": accountID, "principal_kind": "account", "status": "active", "access_token_hash": sha256Hex(accessToken), "refresh_token_hash": sha256Hex(refreshToken), "issued_at_ms": issuedAtMs, "access_expires_at_ms": accessExpiryMs, "refresh_expires_at_ms": refreshExpiryMs}).Insert(); queryErr != nil {
			return queryErr
		}
		persistedSession, queryErr = tx.Model(consts.KidsV1SessionTable).Ctx(ctx).Where("session_id", sessionID).One()
		if queryErr != nil {
			return queryErr
		}
		if persistedSession.IsEmpty() {
			return v1Error(503, "UNAVAILABLE", true, "account session is unavailable")
		}
		rows, queryErr := tx.Model(consts.KidsV1MembershipTable).Ctx(ctx).Where("account_id", accountID).Where("status", "active").Order("created_at ASC,id ASC").All()
		if queryErr != nil {
			return queryErr
		}
		memberships = make([]map[string]any, 0, len(rows))
		for _, row := range rows {
			memberships = append(memberships, v1MembershipRecordProjection(row))
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	cursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"account": account, "session": v1AuthSessionProjection(persistedSession, accessToken, refreshToken), "account_binding": binding, "memberships": memberships, "merge_outcome": mergeOutcome, "bootstrap_cursor": cursor}, "", nil
}

// v1PrepareAsset 为受控演示传输登记上传元数据，不暴露本地路径或对象地址。
func (s *sKids) v1PrepareAsset(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	uploadID := fmt.Sprint(in.Body["upload_id"])
	byteSize, ok := v1Integer(in.Body["byte_size"])
	if !ok || byteSize < 1 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "asset byte size is invalid")
	}
	circleID := nullableV1String(in.Body["circle_id"])
	if circleID != "" {
		if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
			return nil, "", err
		}
	}
	now, expiresAt := time.Now(), time.Now().Add(15*time.Minute)
	_, err := utils.KidsDB(ctx).Model(consts.KidsV1AssetUploadTable).Ctx(ctx).Data(gdb.Map{
		"upload_id": uploadID, "account_id": in.PrincipalID, "circle_id": circleID, "purpose": in.Body["purpose"], "content_type": in.Body["content_type"],
		"byte_size": byteSize, "sha256": in.Body["sha256"], "version": 1, "status": "prepared", "expires_at": expiresAt, "created_at": now, "updated_at": now,
	}).Insert()
	if err != nil {
		return nil, "", v1Error(409, "VERSION_CONFLICT", false, "upload identifier already exists")
	}
	return map[string]any{"upload_target": map[string]any{"upload_id": uploadID, "upload_mode": "demo_transport", "method": "PUT", "upload_url": nil, "expires_at_ms": expiresAt.UnixMilli(), "version": int64(1)}}, "", nil
}

// v1CommitAsset 校验已登记摘要与大小后，将演示传输的资产状态原子提交。
func (s *sKids) v1CommitAsset(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	uploadID := in.PathParameters["upload_id"]
	expected, ok := v1ExpectedVersion(in.Body["expected_upload_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "upload version is invalid")
	}
	byteSize, ok := v1Integer(in.Body["byte_size"])
	if !ok || byteSize < 1 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "asset byte size is invalid")
	}
	now := time.Now()
	var asset, receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		upload, err := tx.Model(consts.KidsV1AssetUploadTable).Ctx(ctx).Where("upload_id", uploadID).Where("account_id", in.PrincipalID).LockUpdate().One()
		if err != nil {
			return err
		}
		if upload.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "asset upload is missing")
		}
		version := upload["version"].Int64()
		if version != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &version, Message: "upload version conflicts"}
		}
		if upload["status"].String() != "prepared" || now.After(upload["expires_at"].Time()) {
			return v1Error(409, "VERSION_CONFLICT", false, "asset upload is no longer available")
		}
		if upload["sha256"].String() != fmt.Sprint(in.Body["sha256"]) || upload["byte_size"].Int64() != byteSize {
			return v1Error(422, "VALIDATION_FAILED", false, "asset metadata differs from prepared upload")
		}
		assetID := v1ID("asset", uploadID)
		if _, err = tx.Model(consts.KidsV1AssetTable).Ctx(ctx).Data(gdb.Map{"asset_id": assetID, "upload_id": uploadID, "circle_id": upload["circle_id"], "purpose": upload["purpose"], "content_type": upload["content_type"], "byte_size": byteSize, "sha256": upload["sha256"], "state": "committed", "version": 1, "committed_at": now, "created_at": now}).Insert(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1AssetUploadTable).Ctx(ctx).Where("id", upload["id"].Int64()).Data(gdb.Map{"status": "committed", "version": version + 1, "updated_at": now}).Update(); err != nil {
			return err
		}
		receipt, err = v1CreateCommitTx(ctx, tx, upload["circle_id"].String(), in.OperationID, map[string]any{})
		if err != nil {
			return err
		}
		asset = map[string]any{"asset_id": assetID, "purpose": upload["purpose"].String(), "content_type": upload["content_type"].String(), "byte_size": byteSize, "sha256": upload["sha256"].String(), "state": "committed", "version": int64(1), "committed_at_ms": now.UnixMilli()}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "asset": asset}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1IdempotencyBegin 在执行副作用前占用幂等键，阻止并发请求重复创建会话。
func v1IdempotencyBegin(ctx context.Context, scope string, in v1.V1OperationInput, routeHash, bodyHash string) (*v1.V1OperationOutput, bool, bool, error) {
	record, err := utils.KidsDB(ctx).Model(consts.KidsV1IdempotencyTable).Ctx(ctx).Where("principal_scope", scope).Where("idempotency_key", in.IdempotencyKey).One()
	if err != nil {
		return nil, false, false, err
	}
	if !record.IsEmpty() {
		return v1IdempotencyRecord(record, routeHash, bodyHash)
	}
	_, err = utils.KidsDB(ctx).Model(consts.KidsV1IdempotencyTable).Ctx(ctx).Data(gdb.Map{
		"principal_scope": scope, "idempotency_key": in.IdempotencyKey, "operation_id": in.OperationID,
		"route_fingerprint": routeHash, "body_fingerprint": bodyHash, "response_status": 0, "response_body": "{}",
	}).Insert()
	if err == nil {
		return nil, false, false, nil
	}
	record, readErr := utils.KidsDB(ctx).Model(consts.KidsV1IdempotencyTable).Ctx(ctx).Where("principal_scope", scope).Where("idempotency_key", in.IdempotencyKey).One()
	if readErr != nil || record.IsEmpty() {
		return nil, false, false, err
	}
	return v1IdempotencyRecord(record, routeHash, bodyHash)
}

// v1IdempotencyRecord 将既有幂等记录归类为冲突、处理中或可重放响应。
func v1IdempotencyRecord(record gdb.Record, routeHash, bodyHash string) (*v1.V1OperationOutput, bool, bool, error) {
	if record["route_fingerprint"].String() != routeHash || record["body_fingerprint"].String() != bodyHash {
		return nil, true, false, nil
	}
	if record["response_status"].Int() == 0 {
		return nil, false, true, nil
	}
	var data map[string]any
	if err := json.Unmarshal(record["response_body"].Bytes(), &data); err != nil {
		return nil, false, false, err
	}
	return &v1.V1OperationOutput{Data: data, Status: record["response_status"].Int(), ChangeCursor: record["response_change_cursor"].String(), ETag: record["response_etag"].String()}, false, false, nil
}

// v1RefreshIdempotencyReplay 在旧 refresh credential 已失效后仍按固定 session 作用域返回首次轮换结果。
func v1RefreshIdempotencyReplay(ctx context.Context, in v1.V1OperationInput, routeHash, bodyHash string) (*v1.V1OperationOutput, bool, bool, error) {
	session, err := utils.KidsDB(ctx).Model(consts.KidsV1SessionTable).Ctx(ctx).
		Where("session_id", fmt.Sprint(in.Body["session_id"])).Where("status", "active").One()
	if err != nil || session.IsEmpty() {
		return nil, false, false, err
	}
	record, err := utils.KidsDB(ctx).Model(consts.KidsV1IdempotencyTable).Ctx(ctx).
		Where("principal_scope", v1RefreshIdempotencyScope(in)).Where("idempotency_key", in.IdempotencyKey).One()
	if err != nil || record.IsEmpty() {
		return nil, false, false, err
	}
	output, conflict, pending, err := v1IdempotencyRecord(record, routeHash, bodyHash)
	if err != nil || output == nil {
		return output, conflict, pending, err
	}
	output = v1ReplayOutput(output)
	if err = v1.ValidateV1ResponseData(in.OperationID, output.Data); err != nil {
		return nil, false, false, v1Error(503, "UNAVAILABLE", true, "stored v1 response is unavailable")
	}
	return output, false, false, nil
}

// v1IdempotencySave 保存接口首次响应，保证断线重放返回同一结果。
func v1IdempotencySave(ctx context.Context, scope string, in v1.V1OperationInput, routeHash, bodyHash string, out *v1.V1OperationOutput) error {
	body, err := json.Marshal(out.Data)
	if err != nil {
		return err
	}
	_, err = utils.KidsDB(ctx).Model(consts.KidsV1IdempotencyTable).Ctx(ctx).
		Where("principal_scope", scope).Where("idempotency_key", in.IdempotencyKey).Where("response_status", 0).
		Data(gdb.Map{"response_status": out.Status, "response_body": string(body), "response_change_cursor": nullableV1String(out.ChangeCursor), "response_etag": nullableV1String(out.ETag)}).Update()
	return err
}

// v1IdempotencySaveTx 在业务写入事务内固化首次响应，避免写入已提交但重放快照未保存。
func v1IdempotencySaveTx(ctx context.Context, tx gdb.TX, scope string, in v1.V1OperationInput, routeHash, bodyHash string, out *v1.V1OperationOutput) error {
	record, err := tx.Model(consts.KidsV1IdempotencyTable).Ctx(ctx).
		Where("principal_scope", scope).Where("idempotency_key", in.IdempotencyKey).LockUpdate().One()
	if err != nil {
		return err
	}
	if record.IsEmpty() || record["response_status"].Int() != 0 {
		return v1Error(503, "UNAVAILABLE", true, "idempotency response store is unavailable")
	}
	if record["route_fingerprint"].String() != routeHash || record["body_fingerprint"].String() != bodyHash {
		return v1Error(409, "IDEMPOTENCY_CONFLICT", false, "idempotency key conflicts with an earlier request")
	}
	body, err := json.Marshal(out.Data)
	if err != nil {
		return err
	}
	_, err = tx.Model(consts.KidsV1IdempotencyTable).Ctx(ctx).Where("id", record["id"].Int64()).Data(gdb.Map{
		"response_status": out.Status, "response_body": string(body), "response_change_cursor": nullableV1String(out.ChangeCursor), "response_etag": nullableV1String(out.ETag),
	}).Update()
	return err
}

// v1IdempotencyAbort 清理未产生副作用的占位记录，使同一请求可安全重试。
func v1IdempotencyAbort(ctx context.Context, scope string, in v1.V1OperationInput, routeHash, bodyHash string) error {
	_, err := utils.KidsDB(ctx).Model(consts.KidsV1IdempotencyTable).Ctx(ctx).
		Where("principal_scope", scope).Where("idempotency_key", in.IdempotencyKey).
		Where("route_fingerprint", routeHash).Where("body_fingerprint", bodyHash).Where("response_status", 0).Delete()
	return err
}

// v1CreateCommitTx 在业务写入的同一事务中分配提交序列并持久化 commit 和 receipt。
func v1CreateCommitTx(ctx context.Context, tx gdb.TX, circleID, operationID string, changes map[string]any) (map[string]any, error) {
	sequenceRecord, err := tx.Model(consts.KidsV1SequenceTable).Ctx(ctx).Where("id", 1).LockUpdate().One()
	if err != nil {
		return nil, err
	}
	if sequenceRecord.IsEmpty() {
		if _, err = tx.Model(consts.KidsV1SequenceTable).Ctx(ctx).Data(gdb.Map{"id": 1, "next_commit_sequence": 1}).InsertIgnore(); err != nil {
			return nil, err
		}
		sequenceRecord, err = tx.Model(consts.KidsV1SequenceTable).Ctx(ctx).Where("id", 1).LockUpdate().One()
		if err != nil {
			return nil, err
		}
		if sequenceRecord.IsEmpty() {
			return nil, v1Error(503, "UNAVAILABLE", true, "v1 commit sequence is unavailable")
		}
	}
	sequence := sequenceRecord["next_commit_sequence"].Int64()
	if sequence < 1 {
		return nil, fmt.Errorf("v1 commit sequence is invalid")
	}
	if _, err = tx.Model(consts.KidsV1SequenceTable).Ctx(ctx).Where("id", 1).Data(gdb.Map{"next_commit_sequence": sequence + 1}).Update(); err != nil {
		return nil, err
	}
	now := time.Now()
	commitID := v1ID("commit", uuid.NewString())
	receiptID := v1ID("receipt", uuid.NewString())
	changePayload, err := json.Marshal(v1SyncChanges(changes))
	if err != nil {
		return nil, err
	}
	if _, err = tx.Model(consts.KidsV1CommitTable).Ctx(ctx).Data(gdb.Map{
		"commit_id": commitID, "circle_id": circleID, "commit_sequence": sequence, "change_payload": string(changePayload),
	}).Insert(); err != nil {
		return nil, err
	}
	if _, err = tx.Model(consts.KidsV1ReceiptTable).Ctx(ctx).Data(gdb.Map{
		"receipt_id": receiptID, "commit_id": commitID, "operation_id": operationID, "result_kind": "first_committed", "committed_at": now,
	}).Insert(); err != nil {
		return nil, err
	}
	return map[string]any{
		"receipt_id": receiptID, "commit_id": commitID, "commit_sequence": sequence,
		"result_kind": "first_committed", "committed_at_ms": now.UnixMilli(),
	}, nil
}

// v1UpdateCommitChangesTx 在同一事务中回填依赖 commit sequence 的完整同步投影。
func v1UpdateCommitChangesTx(ctx context.Context, tx gdb.TX, commitID string, changes map[string]any) error {
	changePayload, err := json.Marshal(v1SyncChanges(changes))
	if err != nil {
		return err
	}
	_, err = tx.Model(consts.KidsV1CommitTable).Ctx(ctx).Where("commit_id", commitID).Data(gdb.Map{"change_payload": string(changePayload)}).Update()
	return err
}

// v1Receipt 构造符合接口字段集的稳定回执，持久化由具体事务负责。
func v1Receipt(operationID string) map[string]any {
	now := time.Now()
	receiptID := v1ID("receipt", uuid.NewString())
	return map[string]any{"receipt_id": receiptID, "commit_id": v1ID("commit", uuid.NewString()), "commit_sequence": now.UnixMilli(), "result_kind": "first_committed", "committed_at_ms": now.UnixMilli()}
}

// v1PrincipalScope 计算接口幂等查询作用域。
func v1PrincipalScope(ctx context.Context, in v1.V1OperationInput) string {
	if in.OperationID == "refreshSession" {
		return v1RefreshIdempotencyScope(in)
	}
	if in.PrincipalID != "" {
		return in.PrincipalKind + ":" + in.PrincipalID + ":" + in.SessionID
	}
	return "public"
}

// v1RefreshIdempotencyScope 以请求 session ID 固定 refresh 重放作用域，避免轮换后丢失首次结果。
func v1RefreshIdempotencyScope(in v1.V1OperationInput) string {
	return "refresh:" + fmt.Sprint(in.Body["session_id"])
}

// v1AccountID 返回服务端解析的账号标识。
func v1AccountID(ctx context.Context, in v1.V1OperationInput) string {
	if in.PrincipalID != "" {
		return in.PrincipalID
	}
	if id := v1UserID(ctx); id > 0 {
		return v1ID("account", id)
	}
	return "public"
}

// v1UserID 从既有 kids JWT 上下文读取服务端用户 ID。
func v1UserID(ctx context.Context) uint64 {
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}

// v1BodyFingerprint 生成请求体摘要。
func v1BodyFingerprint(body map[string]any) string {
	encoded, _ := json.Marshal(body)
	return sha256Hex(string(encoded))
}

// v1RouteFingerprint 生成方法、operation 和 path 参数摘要。
func v1RouteFingerprint(in v1.V1OperationInput) string {
	encoded, _ := json.Marshal(map[string]any{"method": in.Method, "operation_id": in.OperationID, "path_parameters": in.PathParameters})
	return sha256Hex(string(encoded))
}

// v1ID 将内部标识稳定映射为接口要求的 UUID 型 opaque ID。
func v1ID(kind string, value any) string {
	prefixes := map[string]string{
		"account": "account", "adjustment": "adjustment", "admin": "admin", "asset": "asset", "binding": "binding",
		"circle": "circle", "commit": "commit", "completion": "completion", "entitlement": "entitlement", "exchange": "exchange",
		"feedback": "feedback", "invite": "invite", "ledger": "star-transaction", "member": "member", "membership": "membership", "notification": "notification",
		"receipt": "receipt", "reward": "reward", "selection": "selection", "session": "session", "task": "task", "task-tag": "task-tag", "upload": "upload",
	}
	prefix, ok := prefixes[kind]
	if !ok {
		prefix = kind
	}
	identifier := uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("clearwave:%s:%v", kind, value)))
	return fmt.Sprintf("%s:v1:%s", prefix, identifier.String())
}

// v1Cursor 创建满足接口 pattern 的不透明同步游标。
func v1Cursor() string { return "cur:v1:" + strings.ReplaceAll(uuid.NewString(), "-", "") }

// v1CommitCursor 基于持久化提交序列生成稳定同步游标。
func v1CommitCursor(sequence int64) string { return fmt.Sprintf("cur:v1:%08x", sequence) }

// sha256Hex 计算字符串 SHA-256 摘要。
func sha256Hex(value string) string {
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}

// nullableString 将空字符串映射为接口 nullable 字段。
func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

// nullableV1String 将 JSON null 映射为可安全写入接口文本列的空字符串。
func nullableV1String(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

// v1QueryFirst 读取已经通过 schema 校验的单值查询参数。
func v1QueryFirst(query map[string][]string, key string) string {
	values := query[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// v1Error 创建稳定的接口错误。
func v1Error(status int, code string, retryable bool, message string) error {
	return &v1.V1Error{Status: status, Code: code, Retryable: retryable, Message: message}
}

// v1Secret 生成仅用于接口 credential 的随机秘密。
func v1Secret() string {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return uuid.NewString()
	}
	return hex.EncodeToString(value)
}
