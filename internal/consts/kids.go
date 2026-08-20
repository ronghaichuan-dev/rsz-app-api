package consts

// kids 数据库表名和基础格式常量。
const (
	// KidsCircleTable 是 kids 圈子表名。
	KidsCircleTable = "kids_circle"
	// KidsCircleUserTable 是 kids 用户圈子关系表名。
	KidsCircleUserTable = "kids_circle_user"
	// KidsInviteCodeTable 是 kids 邀请码表名。
	KidsInviteCodeTable = "kids_invite_code"
	// KidsUserTable 是 kids 用户表名。
	KidsUserTable = "kids_user"
	// KidsUserAuthTable 是 kids 用户授权身份表名。
	KidsUserAuthTable = "kids_user_auth"
	// KidsUserTokenTable 是 kids 用户访问令牌表名。
	KidsUserTokenTable = "kids_user_token"
	// KidsFamilyMemberTable 是 kids 家庭成员表名。
	KidsFamilyMemberTable = "kids_family_member"
	// KidsTaskPresetTable 是 kids 任务预设表名。
	KidsTaskPresetTable = "kids_task_preset"
	// KidsTaskTable 是 kids 任务表名。
	KidsTaskTable = "kids_task"
	// KidsTaskAssigneeTable 是 kids 任务分配表名。
	KidsTaskAssigneeTable = "kids_task_assignee"
	// KidsTaskTagTable 是 kids 任务标签表名。
	KidsTaskTagTable = "kids_task_tag"
	// KidsRewardTable 是 kids 奖励表名。
	KidsRewardTable = "kids_reward"
	// KidsRewardAssigneeTable 是 kids 奖励指派表名。
	KidsRewardAssigneeTable = "kids_reward_assignee"
	// KidsRewardPresetTable 是 kids 奖励预设表名。
	KidsRewardPresetTable = "kids_reward_preset"
	// KidsRewardRedeemRecordTable 是 kids 奖励兑换记录表名。
	KidsRewardRedeemRecordTable = "kids_reward_redeem_record"
	// KidsUploadFileTable 是 kids 上传文件表名。
	KidsUploadFileTable = "kids_upload_file"
	// KidsDeviceNotificationTable 是 kids 设备通知设置表名。
	KidsDeviceNotificationTable = "kids_device_notification"
	// KidsStarRecordTable 是 kids 星星流水表名。
	KidsStarRecordTable = "kids_star_record"
	// KidsNotificationTable 是 kids 通知表名。
	KidsNotificationTable = "kids_notification"
	// KidsV1SessionTable 是合同会话表名。
	KidsV1SessionTable = "kids_contract_session"
	// KidsV1IdempotencyTable 是合同幂等记录表名。
	KidsV1IdempotencyTable = "kids_contract_idempotency"
	// KidsV1CommitTable 是合同同步提交表名。
	KidsV1CommitTable = "kids_contract_commit"
	// KidsV1ReceiptTable 是合同写入回执表名。
	KidsV1ReceiptTable = "kids_contract_receipt"
	// KidsV1AssetUploadTable 是合同资产上传表名。
	KidsV1AssetUploadTable = "kids_contract_asset_upload"
	// KidsV1AssetTable 是合同已提交资产表名。
	KidsV1AssetTable = "kids_contract_asset"
	// KidsV1FeedbackTable 是合同反馈表名。
	KidsV1FeedbackTable = "kids_contract_feedback"
	// KidsV1EntitlementTable 是合同权益表名。
	KidsV1EntitlementTable = "kids_contract_entitlement"
	// KidsV1SequenceTable 是合同提交序列表名。
	KidsV1SequenceTable = "kids_contract_sequence"
	// KidsV1AccountTable 是合同账号表名。
	KidsV1AccountTable = "kids_contract_account"
	// KidsV1AccountBindingTable 是合同账号绑定表名。
	KidsV1AccountBindingTable = "kids_contract_account_binding"
	// KidsV1CircleTable 是合同圈子表名。
	KidsV1CircleTable = "kids_contract_circle"
	// KidsV1AdministratorTable 是合同管理员表名。
	KidsV1AdministratorTable = "kids_contract_administrator"
	// KidsV1MemberTable 是合同成员表名。
	KidsV1MemberTable = "kids_contract_member"
	// KidsV1MembershipTable 是合同成员身份表名。
	KidsV1MembershipTable = "kids_contract_membership"
	// KidsV1CircleSelectionTable 是合同圈子选择表名。
	KidsV1CircleSelectionTable = "kids_contract_circle_selection"
	// KidsV1InviteTable 是合同邀请码表名。
	KidsV1InviteTable = "kids_contract_invite"
	// KidsV1TaskTagTable 是合同任务标签表名。
	KidsV1TaskTagTable = "kids_contract_task_tag"
	// KidsV1TaskTable 是合同任务定义事实表名。
	KidsV1TaskTable = "kids_contract_task"
	// KidsV1TaskAssignmentTable 是合同任务成员分配事实表名。
	KidsV1TaskAssignmentTable = "kids_contract_task_assignment"
	// KidsV1TaskOccurrenceTable 是合同任务 occurrence 事实表名。
	KidsV1TaskOccurrenceTable = "kids_contract_task_occurrence"
	// KidsV1TaskCompletionTable 是合同任务完成审计事实表名。
	KidsV1TaskCompletionTable = "kids_contract_task_completion"
	// KidsV1TaskCancellationTable 是合同任务取消审计事实表名。
	KidsV1TaskCancellationTable = "kids_contract_task_cancellation"
	// KidsV1LedgerTable 是合同星星流水事实表名。
	KidsV1LedgerTable = "kids_contract_ledger"
	// KidsV1BalanceTable 是合同成员星星余额投影表名。
	KidsV1BalanceTable = "kids_contract_balance"
	// DefaultDBGroup 是当前微服务配置文件中默认数据库分组名。
	DefaultDBGroup = "default"
	// MySQLTimeLayout 是 MySQL DATETIME 字段格式。
	MySQLTimeLayout = "2006-01-02 15:04:05"
	// DateLayout 是业务日期字段格式。
	DateLayout = "2006-01-02"
)

// kids 认证相关常量。
const (
	// KidsAccessTokenTTL 是 kids 访问令牌默认有效期。
	KidsAccessTokenTTL = DefaultAccessTokenTTL
	// KidsV1PermissionManageCircle 表示合同圈子配置管理权限。
	KidsV1PermissionManageCircle = "manage_circle"
	// KidsV1PermissionManageMembers 表示合同成员与管理员管理权限。
	KidsV1PermissionManageMembers = "manage_members"
	// KidsV1PermissionManageTasks 表示合同任务管理权限。
	KidsV1PermissionManageTasks = "manage_tasks"
	// KidsV1PermissionManageRewards 表示合同奖励管理权限。
	KidsV1PermissionManageRewards = "manage_rewards"
	// KidsV1PermissionAdjustStars 表示合同星星调整权限。
	KidsV1PermissionAdjustStars = "adjust_stars"
)

// kids 任务规则相关常量。
const (
	// TaskRepeatNone 表示任务不重复。
	TaskRepeatNone = "none"
	// TaskRepeatDaily 表示任务每天重复。
	TaskRepeatDaily = "daily"
	// TaskRepeatWeekly 表示任务每周重复。
	TaskRepeatWeekly = "weekly"
	// TaskRepeatMonthly 表示任务每月重复。
	TaskRepeatMonthly = "monthly"
	// TaskRepeatYearly 表示任务每年重复。
	TaskRepeatYearly = "yearly"
	// TaskRepeatCustom 表示任务使用自定义重复规则。
	TaskRepeatCustom = "custom"
	// TaskRepeatEndNever 表示重复任务不设置结束条件。
	TaskRepeatEndNever = "never"
	// TaskRepeatEndDate 表示重复任务按日期结束。
	TaskRepeatEndDate = "date"
	// TaskRepeatEndCount 表示重复任务按次数结束。
	TaskRepeatEndCount = "count"
	// TaskTimeLimitAllDay 表示任务全天可完成。
	TaskTimeLimitAllDay = "all_day"
	// TaskTimeLimitRange 表示任务只能在指定时间段完成。
	TaskTimeLimitRange = "range"
	// TaskReminderNone 表示任务不提醒。
	TaskReminderNone = "none"
	// TaskReminderAtTime 表示任务在固定时间提醒。
	TaskReminderAtTime = "at_time"
	// TaskReminderBeforeStart 表示任务开始前提醒。
	TaskReminderBeforeStart = "before_start"
	// MaxTaskRepeatGenerateDays 是一次创建重复任务最多预生成的天数。
	MaxTaskRepeatGenerateDays = 365
)

// kids 奖励相关常量。
const (
	// RewardRepeatNone 表示每个儿童只能兑换一次。
	RewardRepeatNone = "none"
	// RewardRepeatDaily 表示每天可兑换一次。
	RewardRepeatDaily = "daily"
	// RewardRepeatWeekly 表示每周可兑换一次。
	RewardRepeatWeekly = "weekly"
	// RewardRepeatMonthly 表示每月可兑换一次。
	RewardRepeatMonthly = "monthly"
	// RewardRepeatCustom 表示按自定义天数间隔兑换。
	RewardRepeatCustom = "custom"
)
