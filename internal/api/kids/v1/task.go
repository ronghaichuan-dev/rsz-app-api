package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	TaskCompletionModeSingle   = "single"
	TaskCompletionModeRotation = "rotation"
	TaskCompletionModeAnyone   = "anyone"
	TaskCompletionModeEveryone = "everyone"
)

type TaskPreset struct {
	Id          uint64 `json:"id" dc:"预设ID"`
	Title       string `json:"title" dc:"预设任务标题"`
	Icon        string `json:"icon" dc:"预设图标标识"`
	Star        int    `json:"star" dc:"默认星星数量"`
	NeedPhoto   bool   `json:"needPhoto" dc:"是否建议照片凭证"`
	Description string `json:"description" dc:"预设描述"`
}

type TaskAssigneeStatus struct {
	KidId       uint64 `json:"kidId" dc:"分配儿童ID"`
	Order       int    `json:"order" dc:"轮流模式分配顺序"`
	Active      bool   `json:"active" dc:"该儿童今天是否需要完成任务"`
	Completed   bool   `json:"completed" dc:"该儿童是否已完成"`
	PhotoUrl    string `json:"photoUrl" dc:"照片凭证地址"`
	CompletedAt int64  `json:"completedAt" dc:"完成时间戳"`
}

type TaskItem struct {
	Id                    uint64               `json:"id" dc:"任务ID"`
	Title                 string               `json:"title" dc:"任务标题"`
	Icon                  string               `json:"icon" dc:"图标标识"`
	Note                  string               `json:"note" dc:"任务说明"`
	Star                  int                  `json:"star" dc:"星星奖励数量"`
	Date                  string               `json:"date" dc:"计划日期"`
	KidIds                []uint64             `json:"kidIds" dc:"分配儿童ID列表"`
	Assignees             []TaskAssigneeStatus `json:"assignees" dc:"分配儿童完成状态列表"`
	CompletionMode        string               `json:"completionMode" dc:"完成模式"`
	RepeatRule            string               `json:"repeatRule" dc:"重复规则"`
	RepeatEndType         string               `json:"repeatEndType" dc:"重复结束类型：never/date/count"`
	RepeatEndDate         string               `json:"repeatEndDate" dc:"重复结束日期"`
	RepeatEndCount        int                  `json:"repeatEndCount" dc:"重复结束次数"`
	TimeLimitType         string               `json:"timeLimitType" dc:"时间限制类型：all_day/range"`
	TimeLimitStart        string               `json:"timeLimitStart" dc:"开始时间，格式HH:mm"`
	TimeLimitEnd          string               `json:"timeLimitEnd" dc:"结束时间，格式HH:mm"`
	ReminderType          string               `json:"reminderType" dc:"提醒类型：none/at_time/before_start"`
	ReminderAt            string               `json:"reminderAt" dc:"提醒时间，格式HH:mm"`
	ReminderOffsetMinutes int                  `json:"reminderOffsetMinutes" dc:"提前提醒分钟数"`
	NeedPhotoProof        bool                 `json:"needPhotoProof" dc:"是否需要照片凭证"`
	TagId                 uint64               `json:"tagId" dc:"标签ID"`
	TagName               string               `json:"tagName" dc:"标签名称"`
	TagColor              string               `json:"tagColor" dc:"标签颜色"`
	CanComplete           bool                 `json:"canComplete" dc:"当前日期是否允许完成"`
	Completed             bool                 `json:"completed" dc:"任务是否按完成模式完成"`
	CompletedBy           uint64               `json:"completedBy" dc:"最近完成的儿童ID"`
	PhotoUrl              string               `json:"photoUrl" dc:"最近照片凭证地址"`
	CompletedAt           int64                `json:"completedAt" dc:"最近完成时间戳"`
}
type TaskPresetListReq struct {
	g.Meta  `path:"/task-presets" method:"get" tags:"儿童端任务" summary:"查询任务预设列表"`
	Keyword string `p:"keyword" dc:"搜索关键字"`
}

type TaskPresetListInput struct {
	Keyword string
}

type TaskPresetListOutput struct {
	List []TaskPreset
}

type TaskPresetListRes struct {
	List []TaskPreset `json:"list" dc:"任务预设列表"`
}

type TaskListReq struct {
	g.Meta `path:"/tasks" method:"get" tags:"儿童端任务" summary:"按日期、儿童和标签查询任务"`
	Date   string `p:"date" dc:"任务日期，格式YYYY-MM-DD"`
	KidId  uint64 `p:"kidId" dc:"儿童成员ID"`
	TagId  uint64 `p:"tagId" dc:"标签ID"`
	Status string `p:"status" dc:"任务状态筛选：pending待完成、completed已完成"`
}

type TaskListInput struct {
	Date   string
	KidId  uint64
	TagId  uint64
	Status string
}

type TaskListOutput struct {
	Date string
	List []TaskItem
}

type TaskListRes struct {
	Date string     `json:"date" dc:"任务日期"`
	List []TaskItem `json:"list" dc:"任务列表"`
}

type TaskCreateReq struct {
	g.Meta                `path:"/tasks" method:"post" tags:"儿童端任务" summary:"创建任务"`
	Title                 string   `json:"title" v:"required|max-length:128" dc:"任务标题"`
	Icon                  string   `json:"icon" dc:"图标标识"`
	Note                  string   `json:"note" dc:"任务说明"`
	Star                  int      `json:"star" v:"required|min:1|max:999" dc:"星星奖励数量"`
	KidIds                []uint64 `json:"kidIds" v:"required" dc:"分配儿童ID列表"`
	CompletionMode        string   `json:"completionMode" v:"required|in:single,rotation,anyone,everyone" dc:"完成模式"`
	RepeatRule            string   `json:"repeatRule" v:"in:none,daily,weekly,monthly,yearly,custom" dc:"重复规则：none、daily、weekly、monthly、yearly、custom"`
	RepeatEndType         string   `json:"repeatEndType" v:"in:never,date,count" dc:"重复结束类型：never、date、count"`
	RepeatEndDate         string   `json:"repeatEndDate" dc:"重复结束日期，格式YYYY-MM-DD"`
	RepeatEndCount        int      `json:"repeatEndCount" dc:"重复结束次数"`
	StartDate             string   `json:"startDate" dc:"开始日期，格式YYYY-MM-DD"`
	EndDate               string   `json:"endDate" dc:"兼容字段：结束日期，格式YYYY-MM-DD"`
	TimeLimitType         string   `json:"timeLimitType" v:"in:all_day,range" dc:"时间限制类型：all_day、range"`
	TimeLimitStart        string   `json:"timeLimitStart" dc:"开始时间，格式HH:mm"`
	TimeLimitEnd          string   `json:"timeLimitEnd" dc:"结束时间，格式HH:mm"`
	ReminderType          string   `json:"reminderType" v:"in:none,at_time,before_start" dc:"提醒类型：none、at_time、before_start"`
	ReminderAt            string   `json:"reminderAt" dc:"提醒时间，格式HH:mm"`
	ReminderOffsetMinutes int      `json:"reminderOffsetMinutes" dc:"提前提醒分钟数"`
	NeedPhotoProof        bool     `json:"needPhotoProof" dc:"是否需要照片凭证"`
	TagId                 uint64   `json:"tagId" dc:"标签ID"`
}

type TaskCreateInput struct {
	Title                 string
	Icon                  string
	Note                  string
	Star                  int
	KidIds                []uint64
	CompletionMode        string
	RepeatRule            string
	RepeatEndType         string
	RepeatEndDate         string
	RepeatEndCount        int
	StartDate             string
	EndDate               string
	TimeLimitType         string
	TimeLimitStart        string
	TimeLimitEnd          string
	ReminderType          string
	ReminderAt            string
	ReminderOffsetMinutes int
	NeedPhotoProof        bool
	TagId                 uint64
}
type TaskCreateOutput struct {
	Task TaskItem
}

type TaskCreateRes struct {
	Task TaskItem `json:"task" dc:"创建的任务"`
}

type TaskCompleteReq struct {
	g.Meta   `path:"/tasks/{id}/complete" method:"post" tags:"儿童端任务" summary:"完成任务并可提交照片凭证"`
	Id       uint64 `p:"id" v:"required|min:1" dc:"任务ID"`
	KidId    uint64 `json:"kidId" v:"required|min:1" dc:"完成任务的儿童成员ID"`
	PhotoUrl string `json:"photoUrl" dc:"照片凭证地址"`
}

type TaskCompleteInput struct {
	Id       uint64
	KidId    uint64
	PhotoUrl string
}

type TaskCompleteOutput struct {
	Task          TaskItem
	StarBalance   int
	TaskCompleted bool
}

type TaskCompleteRes struct {
	Task          TaskItem `json:"task" dc:"已完成任务"`
	StarBalance   int      `json:"starBalance" dc:"完成后的儿童星星余额"`
	TaskCompleted bool     `json:"taskCompleted" dc:"整条任务是否按完成模式完成"`
}

type TaskDetailReq struct {
	g.Meta `path:"/tasks/{id}" method:"get" tags:"儿童端任务" summary:"查询任务详情"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"任务ID"`
}

type TaskDetailInput struct {
	Id uint64
}

type TaskDetailOutput struct {
	Task TaskItem
}

type TaskDetailRes struct {
	Task TaskItem `json:"task" dc:"任务详情"`
}

type TaskUpdateReq struct {
	g.Meta                `path:"/tasks/{id}" method:"put" tags:"儿童端任务" summary:"编辑任务"`
	Id                    uint64   `p:"id" v:"required|min:1" dc:"任务ID"`
	Title                 string   `json:"title" v:"required|max-length:128" dc:"任务标题"`
	Icon                  string   `json:"icon" dc:"图标标识"`
	Note                  string   `json:"note" dc:"任务说明"`
	Star                  int      `json:"star" v:"required|min:1|max:999" dc:"星星奖励数量"`
	KidIds                []uint64 `json:"kidIds" v:"required" dc:"分配儿童ID列表"`
	CompletionMode        string   `json:"completionMode" v:"required|in:single,rotation,anyone,everyone" dc:"完成模式"`
	RepeatRule            string   `json:"repeatRule" v:"in:none,daily,weekly,monthly,yearly,custom" dc:"重复规则：none、daily、weekly、monthly、yearly、custom"`
	RepeatEndType         string   `json:"repeatEndType" v:"in:never,date,count" dc:"重复结束类型：never、date、count"`
	RepeatEndDate         string   `json:"repeatEndDate" dc:"重复结束日期"`
	RepeatEndCount        int      `json:"repeatEndCount" dc:"重复结束次数"`
	TimeLimitType         string   `json:"timeLimitType" v:"in:all_day,range" dc:"时间限制类型：all_day、range"`
	TimeLimitStart        string   `json:"timeLimitStart" dc:"开始时间，格式HH:mm"`
	TimeLimitEnd          string   `json:"timeLimitEnd" dc:"结束时间，格式HH:mm"`
	ReminderType          string   `json:"reminderType" v:"in:none,at_time,before_start" dc:"提醒类型：none、at_time、before_start"`
	ReminderAt            string   `json:"reminderAt" dc:"提醒时间，格式HH:mm"`
	ReminderOffsetMinutes int      `json:"reminderOffsetMinutes" dc:"提前提醒分钟数"`
	NeedPhotoProof        bool     `json:"needPhotoProof" dc:"是否需要照片凭证"`
	TagId                 uint64   `json:"tagId" dc:"标签ID"`
}

type TaskUpdateInput struct {
	Id                    uint64
	Title                 string
	Icon                  string
	Note                  string
	Star                  int
	KidIds                []uint64
	CompletionMode        string
	RepeatRule            string
	RepeatEndType         string
	RepeatEndDate         string
	RepeatEndCount        int
	TimeLimitType         string
	TimeLimitStart        string
	TimeLimitEnd          string
	ReminderType          string
	ReminderAt            string
	ReminderOffsetMinutes int
	NeedPhotoProof        bool
	TagId                 uint64
}

type TaskUpdateOutput struct {
	Task TaskItem
}

type TaskUpdateRes struct {
	Task TaskItem `json:"task" dc:"更新后的任务"`
}

type TaskDeleteReq struct {
	g.Meta `path:"/tasks/{id}" method:"delete" tags:"儿童端任务" summary:"删除任务"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"任务ID"`
}

type TaskDeleteInput struct {
	Id uint64
}

type TaskDeleteOutput struct {
	Id uint64
}

type TaskDeleteRes struct {
	Id uint64 `json:"id" dc:"已删除任务ID"`
}

type TaskCancelReq struct {
	g.Meta `path:"/tasks/{id}/cancel" method:"post" tags:"儿童端任务" summary:"取消任务完成状态"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"任务ID"`
	KidId  uint64 `json:"kidId" v:"required|min:1" dc:"儿童成员ID"`
	Reason string `json:"reason" dc:"取消原因"`
}

type TaskCancelInput struct {
	Id     uint64
	KidId  uint64
	Reason string
}

type TaskCancelOutput struct {
	Task        TaskItem
	StarBalance int
}

type TaskCancelRes struct {
	Task        TaskItem `json:"task" dc:"取消后的任务"`
	StarBalance int      `json:"starBalance" dc:"取消后的星星余额"`
}

type TaskTag struct {
	Id        uint64 `json:"id" dc:"标签ID"`
	Name      string `json:"name" dc:"标签名称"`
	Color     string `json:"color" dc:"标签颜色"`
	SortOrder int    `json:"sortOrder" dc:"排序值"`
}

type TaskTagListReq struct {
	g.Meta `path:"/task-tags" method:"get" tags:"儿童端任务标签" summary:"查询任务标签列表"`
}

type TaskTagListInput struct{}

type TaskTagListOutput struct {
	List []TaskTag
}

type TaskTagListRes struct {
	List []TaskTag `json:"list" dc:"标签列表"`
}

type TaskTagCreateReq struct {
	g.Meta    `path:"/task-tags" method:"post" tags:"儿童端任务标签" summary:"创建任务标签"`
	Name      string `json:"name" v:"required|max-length:32" dc:"标签名称"`
	Color     string `json:"color" dc:"标签颜色"`
	SortOrder int    `json:"sortOrder" dc:"排序值"`
}

type TaskTagCreateInput struct {
	Name      string
	Color     string
	SortOrder int
}

type TaskTagCreateOutput struct {
	Tag TaskTag
}

type TaskTagCreateRes struct {
	Tag TaskTag `json:"tag" dc:"创建的标签"`
}

type TaskTagUpdateReq struct {
	g.Meta    `path:"/task-tags/{id}" method:"put" tags:"儿童端任务标签" summary:"更新任务标签"`
	Id        uint64 `p:"id" v:"required|min:1" dc:"标签ID"`
	Name      string `json:"name" v:"required|max-length:32" dc:"标签名称"`
	Color     string `json:"color" dc:"标签颜色"`
	SortOrder int    `json:"sortOrder" dc:"排序值"`
}

type TaskTagUpdateInput struct {
	Id        uint64
	Name      string
	Color     string
	SortOrder int
}

type TaskTagUpdateOutput struct {
	Tag TaskTag
}

type TaskTagUpdateRes struct {
	Tag TaskTag `json:"tag" dc:"更新后的标签"`
}

type TaskTagDeleteReq struct {
	g.Meta `path:"/task-tags/{id}" method:"delete" tags:"儿童端任务标签" summary:"删除任务标签"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"标签ID"`
}

type TaskTagDeleteInput struct {
	Id uint64
}

type TaskTagDeleteOutput struct {
	Id uint64
}

type TaskTagDeleteRes struct {
	Id uint64 `json:"id" dc:"已删除标签ID"`
}
