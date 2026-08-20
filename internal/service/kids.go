// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "rslytics-app-api/internal/api/kids/v1"
)

type (
	IKids interface {
		// GetStatisticsV1 以接口任务完成和星星账本事实返回成员统计序列。
		GetStatisticsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CompareStatisticsV1 在相同快照、时区和 bucket 规则下比较两名成员。
		CompareStatisticsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// GetAnalyticsSummary 查询单个儿童在指定时间范围内的任务或星星统计。
		GetAnalyticsSummary(ctx context.Context, in v1.AnalyticsSummaryInput) (*v1.AnalyticsSummaryOutput, error)
		// ListCompletedTaskDetails 查询指定儿童在时间范围内的已完成任务明细。
		ListCompletedTaskDetails(ctx context.Context, in v1.CompletedTaskListInput) (*v1.CompletedTaskListOutput, error)
		// CompareAnalytics 查询两个儿童在同一时间范围内的任务或星星统计对比。
		CompareAnalytics(ctx context.Context, in v1.AnalyticsCompareInput) (*v1.AnalyticsCompareOutput, error)
		// CreateCircle 创建管理员圈子，并把当前用户加入为管理员。
		CreateCircle(ctx context.Context, in v1.CircleCreateInput) (*v1.CircleCreateOutput, error)
		// CreateInviteCode 创建管理员或成员邀请码，默认有效期为 72 小时。
		CreateInviteCode(ctx context.Context, in v1.InviteCodeCreateInput) (*v1.InviteCodeCreateOutput, error)
		// PreviewInviteCode 查询邀请码对应的圈子和角色信息。
		PreviewInviteCode(ctx context.Context, in v1.InviteCodePreviewInput) (*v1.InviteCodePreviewOutput, error)
		// JoinCircle 使用邀请码加入圈子，并按邀请码角色建立圈子用户关系。
		JoinCircle(ctx context.Context, in v1.CircleJoinInput) (*v1.CircleJoinOutput, error)
		// ListCircles 查询当前用户管理和加入的群组。
		ListCircles(ctx context.Context, in v1.CircleListInput) (*v1.CircleListOutput, error)
		// UpdateCircle 校验管理权限后更新群组名称和图标。
		UpdateCircle(ctx context.Context, in v1.CircleUpdateInput) (*v1.CircleUpdateOutput, error)
		// DeleteCircle 仅群组所有者可以软删除群组。
		DeleteCircle(ctx context.Context, in v1.CircleDeleteInput) (*v1.CircleDeleteOutput, error)
		// LeaveCircle 允许非所有者退出已加入群组。
		LeaveCircle(ctx context.Context, in v1.CircleLeaveInput) (*v1.CircleLeaveOutput, error)
		// ListCircleMembers 查询群组管理员和成员绑定状态。
		ListCircleMembers(ctx context.Context, in v1.CircleMemberListInput) (*v1.CircleMemberListOutput, error)
		// RemoveCircleAdmin 仅群组所有者可以移除管理员。
		RemoveCircleAdmin(ctx context.Context, in v1.CircleAdminRemoveInput) (*v1.CircleAdminRemoveOutput, error)
		// SaveDeviceNotification 保存设备通知授权和偏好。
		SaveDeviceNotification(ctx context.Context, in v1.DeviceNotificationSaveInput) (*v1.DeviceNotificationSaveOutput, error)
		// ListFamilyMembers 从数据库读取家庭成员，并按成人和儿童分组返回。
		ListFamilyMembers(ctx context.Context, in v1.FamilyMemberListInput) (*v1.FamilyMemberListOutput, error)
		// CreateFamilyMember 校验并持久化一个家庭成员。
		CreateFamilyMember(ctx context.Context, in v1.FamilyMemberCreateInput) (*v1.FamilyMemberCreateOutput, error)
		// GetFamilyMember 查询家庭成员详情。
		GetFamilyMember(ctx context.Context, in v1.FamilyMemberDetailInput) (*v1.FamilyMemberDetailOutput, error)
		// UpdateFamilyMember 校验管理权限后更新家庭成员资料。
		UpdateFamilyMember(ctx context.Context, in v1.FamilyMemberUpdateInput) (*v1.FamilyMemberUpdateOutput, error)
		// DeleteFamilyMember 校验管理权限后软删除家庭成员。
		DeleteFamilyMember(ctx context.Context, in v1.FamilyMemberDeleteInput) (*v1.FamilyMemberDeleteOutput, error)
		// Get 保留 GoFrame service 示例方法，当前不承载业务逻辑。
		Get(ctx context.Context) error
		// ListNotifications 从数据库查询通知列表，支持按成员和未读状态筛选。
		ListNotifications(ctx context.Context, in v1.NotificationListInput) (*v1.NotificationListOutput, error)
		// ReadNotification 将指定通知持久化标记为已读。
		ReadNotification(ctx context.Context, in v1.NotificationReadInput) (*v1.NotificationReadOutput, error)
		// GetCurrentAccount 获取当前账号接口 bootstrap。
		GetAccountBootstrap(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// SelectCurrentCircle 选择当前接口圈子。
		SelectCircle(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CreateInviteGuestSession 创建接口游客会话。
		CreateGuestSession(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// RefreshSession 刷新接口会话。
		RefreshV1Session(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// RevokeSession 撤销接口会话。
		RevokeSession(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// ListMyCircles 查询接口圈子列表。
		ListV1Circles(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// GetCircleBootstrap 获取接口圈子 bootstrap。
		GetV1CircleBootstrap(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// UpdateCircleWithVersion 更新接口圈子。
		UpdateCircleWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// DeleteCircleWithVersion 删除接口圈子。
		DeleteCircleWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// LeaveCircleWithVersion 退出接口圈子。
		LeaveCircleWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CreateCircleMember 创建接口成员。
		CreateCircleMember(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// UpsertCircleMember 更新接口成员。
		UpsertCircleMember(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// DeleteCircleMember 删除接口成员。
		DeleteCircleMember(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// UpsertCircleAdministrator 更新接口管理员。
		UpsertCircleAdministrator(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// DeleteCircleAdministrator 删除接口管理员。
		DeleteCircleAdministrator(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CreateCircleInvite 创建接口邀请。
		CreateCircleInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// RefreshCircleInvite 刷新接口邀请。
		RefreshCircleInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// RevokeCircleInvite 撤销接口邀请。
		RevokeCircleInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// RedeemAdministratorInvite 兑换管理员接口邀请。
		RedeemAdministratorInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// RedeemMemberInvite 兑换成员接口邀请。
		RedeemMemberInvite(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// GetCurrentEntitlement 获取接口权益。
		GetEntitlement(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// SubmitFeedback 提交接口反馈。
		SubmitFeedbackV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CompleteOnboarding 完成接口 onboarding。
		CompleteOnboardingV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// UpsertTaskTagWithVersion 更新接口任务标签。
		UpsertTaskTagWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// DeleteTaskTagWithVersion 删除接口任务标签。
		DeleteTaskTagWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// UpsertTaskWithVersion 更新接口任务。
		UpsertTaskWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// DeleteTaskWithVersion 删除接口任务。
		DeleteTaskWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// ListTaskOccurrences 查询接口任务 occurrence。
		ListTaskOccurrencesV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CompleteTaskWithVersion 完成接口任务。
		CompleteTaskWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// CancelTaskCompletionWithVersion 取消接口任务完成记录。
		CancelTaskCompletionWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// ListTaskCompletionDetails 查询接口任务完成明细。
		ListTaskCompletionDetailsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// AdjustMemberStarsWithVersion 调整接口成员星星。
		AdjustMemberStarsWithVersion(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// Unavailable 返回尚未实现的接口 operation 错误。
		Unavailable(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// ExecuteV1 执行尚未拆分为专用 service 方法的接口 operation。
		ExecuteV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// GetStarRanking 查询群组内儿童星星排行榜。
		GetStarRanking(ctx context.Context, in v1.StarRankingInput) (*v1.StarRankingOutput, error)
		// ListRewardPresets 从数据库读取启用的奖励预设，并支持按关键字搜索。
		ListRewardPresets(ctx context.Context, in v1.RewardPresetListInput) (*v1.RewardPresetListOutput, error)
		// ListRewards 从数据库读取奖励列表，支持按群组、儿童和兑换进度筛选。
		ListRewards(ctx context.Context, in v1.RewardListInput) (*v1.RewardListOutput, error)
		// GetReward 查询单个未删除奖励详情。
		GetReward(ctx context.Context, in v1.RewardDetailInput) (*v1.RewardDetailOutput, error)
		// CreateReward 校验奖励入参并持久化创建奖励。
		CreateReward(ctx context.Context, in v1.RewardCreateInput) (*v1.RewardCreateOutput, error)
		// UpdateReward 校验奖励入参并更新奖励配置和指派儿童。
		UpdateReward(ctx context.Context, in v1.RewardUpdateInput) (*v1.RewardUpdateOutput, error)
		// DeleteReward 软删除奖励并保留兑换历史。
		DeleteReward(ctx context.Context, in v1.RewardDeleteInput) (*v1.RewardDeleteOutput, error)
		// RedeemReward 在事务中完成奖励兑换、扣减库存、校验重复周期并写入星星流水。
		RedeemReward(ctx context.Context, in v1.RewardRedeemInput) (*v1.RewardRedeemOutput, error)
		// ListRewardRedeemRecords 从数据库查询奖励兑换历史，支持月份和日期范围筛选。
		ListRewardRedeemRecords(ctx context.Context, in v1.RewardRecordListInput) (*v1.RewardRecordListOutput, error)
		// GetMemberBalancesV1 从接口余额投影读取指定成员的当前星星余额。
		GetMemberBalancesV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// ListStarTransactionsV1 查询接口账本中指定成员的追加式星星流水。
		ListStarTransactionsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error)
		// GetStarBalance 从星星流水表读取指定儿童的当前星星余额。
		GetStarBalance(ctx context.Context, in v1.StarBalanceInput) (*v1.StarBalanceOutput, error)
		// AdjustStars 在事务中写入手动调整流水，并返回调整后的星星余额。
		AdjustStars(ctx context.Context, in v1.StarAdjustInput) (*v1.StarAdjustOutput, error)
		// ListStarRecords 从数据库按儿童、类型和时间范围筛选星星流水。
		ListStarRecords(ctx context.Context, in v1.StarRecordListInput) (*v1.StarRecordListOutput, error)
		// ListTaskPresets 从数据库读取启用的任务预设，并支持按关键字搜索。
		ListTaskPresets(ctx context.Context, in v1.TaskPresetListInput) (*v1.TaskPresetListOutput, error)
		// GetTask 查询未删除任务详情，并返回分配状态和标签信息。
		GetTask(ctx context.Context, in v1.TaskDetailInput) (*v1.TaskDetailOutput, error)
		// ListTasks 按日期、孩子、标签和状态从数据库查询任务列表。
		ListTasks(ctx context.Context, in v1.TaskListInput) (*v1.TaskListOutput, error)
		// CreateTask 校验任务入参，并在事务中按重复规则持久化任务和任务分配关系。
		CreateTask(ctx context.Context, in v1.TaskCreateInput) (*v1.TaskCreateOutput, error)
		// UpdateTask 校验任务入参并更新未删除任务，同时重建任务分配关系。
		UpdateTask(ctx context.Context, in v1.TaskUpdateInput) (*v1.TaskUpdateOutput, error)
		// DeleteTask 软删除任务，保留完成历史和星星流水审计记录。
		DeleteTask(ctx context.Context, in v1.TaskDeleteInput) (*v1.TaskDeleteOutput, error)
		// CompleteTask 在事务中按完成模式完成任务、保存照片凭证，并写入星星流水。
		CompleteTask(ctx context.Context, in v1.TaskCompleteInput) (*v1.TaskCompleteOutput, error)
		// CancelTask 取消指定儿童的任务完成状态，并写入负向星星流水回滚奖励。
		CancelTask(ctx context.Context, in v1.TaskCancelInput) (*v1.TaskCancelOutput, error)
		// ListTaskTags 从数据库读取未删除的任务标签列表。
		ListTaskTags(ctx context.Context, in v1.TaskTagListInput) (*v1.TaskTagListOutput, error)
		// CreateTaskTag 校验标签名称并持久化创建任务标签。
		CreateTaskTag(ctx context.Context, in v1.TaskTagCreateInput) (*v1.TaskTagCreateOutput, error)
		// UpdateTaskTag 校验标签名称并更新任务标签。
		UpdateTaskTag(ctx context.Context, in v1.TaskTagUpdateInput) (*v1.TaskTagUpdateOutput, error)
		// DeleteTaskTag 软删除任务标签，历史任务仍保留标签 ID。
		DeleteTaskTag(ctx context.Context, in v1.TaskTagDeleteInput) (*v1.TaskTagDeleteOutput, error)
		// UploadFile 登记上传文件并返回文件地址。
		UploadFile(ctx context.Context, in v1.UploadFileInput) (*v1.UploadFileOutput, error)
		// Login 根据登录方式完成 kids 用户登录，并委托公共登录模块持久化用户、授权身份和 token。
		Login(ctx context.Context, in v1.UserLoginInput) (*v1.UserLoginOutput, error)
		// GetProfile 从数据库读取当前用户资料；正式鉴权接入前允许通过 userId 查询。
		GetProfile(ctx context.Context, in v1.ProfileGetInput) (*v1.ProfileGetOutput, error)
		// UpdateProfile 更新当前用户昵称和头像。
		UpdateProfile(ctx context.Context, in v1.ProfileUpdateInput) (*v1.ProfileUpdateOutput, error)
	}
)

var (
	localKids IKids
)

func Kids() IKids {
	if localKids == nil {
		panic("implement not found for interface IKids, forgot register?")
	}
	return localKids
}

func RegisterKids(i IKids) {
	localKids = i
}
