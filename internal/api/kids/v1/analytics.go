package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	AnalyticsMetricTasks  = "tasks"
	AnalyticsMetricStars  = "stars"
	AnalyticsRangeWeekly  = "weekly"
	AnalyticsRangeMonthly = "monthly"
	AnalyticsRangeCustom  = "custom"
)

type AnalyticsPoint struct {
	Label string `json:"label" dc:"展示标签"`
	Date  string `json:"date" dc:"日期，格式YYYY-MM-DD"`
	Hour  int    `json:"hour" dc:"小时，0-23"`
	Value int    `json:"value" dc:"统计值"`
}

type AnalyticsSummaryReq struct {
	g.Meta   `path:"/analytics/summary" method:"get" tags:"儿童端数据" summary:"查询任务或星星统计汇总"`
	KidId    uint64 `p:"kidId" v:"required|min:1" dc:"儿童成员ID"`
	Metric   string `p:"metric" v:"required|in:tasks,stars" dc:"指标：tasks任务完成数、stars星星数量"`
	Range    string `p:"range" v:"in:weekly,monthly,custom" dc:"范围：weekly、monthly、custom"`
	From     string `p:"from" dc:"开始日期，格式YYYY-MM-DD"`
	To       string `p:"to" dc:"结束日期，格式YYYY-MM-DD"`
	BaseDate string `p:"baseDate" dc:"周/月基准日期，格式YYYY-MM-DD"`
}

type AnalyticsSummaryInput struct {
	KidId    uint64
	Metric   string
	Range    string
	From     string
	To       string
	BaseDate string
}

type AnalyticsSummaryOutput struct {
	KidId  uint64
	Metric string
	Range  string
	From   string
	To     string
	Total  int
	Daily  []AnalyticsPoint
	Hourly []AnalyticsPoint
}

type AnalyticsSummaryRes struct {
	KidId  uint64           `json:"kidId" dc:"儿童成员ID"`
	Metric string           `json:"metric" dc:"指标"`
	Range  string           `json:"range" dc:"范围"`
	From   string           `json:"from" dc:"开始日期"`
	To     string           `json:"to" dc:"结束日期"`
	Total  int              `json:"total" dc:"总数"`
	Daily  []AnalyticsPoint `json:"daily" dc:"按天统计"`
	Hourly []AnalyticsPoint `json:"hourly" dc:"按小时统计"`
}

type CompletedTaskDetail struct {
	TaskId      uint64 `json:"taskId" dc:"任务ID"`
	KidId       uint64 `json:"kidId" dc:"儿童成员ID"`
	Title       string `json:"title" dc:"任务标题"`
	Icon        string `json:"icon" dc:"图标标识"`
	Star        int    `json:"star" dc:"星星数量"`
	PhotoUrl    string `json:"photoUrl" dc:"照片凭证地址"`
	CompletedAt int64  `json:"completedAt" dc:"完成时间戳"`
}

type CompletedTaskListReq struct {
	g.Meta   `path:"/analytics/tasks/completed" method:"get" tags:"儿童端数据" summary:"查询已完成任务明细"`
	KidId    uint64 `p:"kidId" v:"required|min:1" dc:"儿童成员ID"`
	From     string `p:"from" dc:"开始日期，格式YYYY-MM-DD"`
	To       string `p:"to" dc:"结束日期，格式YYYY-MM-DD"`
	Range    string `p:"range" v:"in:weekly,monthly,custom" dc:"范围：weekly、monthly、custom"`
	BaseDate string `p:"baseDate" dc:"周/月基准日期，格式YYYY-MM-DD"`
}

type CompletedTaskListInput struct {
	KidId    uint64
	From     string
	To       string
	Range    string
	BaseDate string
}

type CompletedTaskListOutput struct {
	From  string
	To    string
	Total int
	List  []CompletedTaskDetail
}

type CompletedTaskListRes struct {
	From  string                `json:"from" dc:"开始日期"`
	To    string                `json:"to" dc:"结束日期"`
	Total int                   `json:"total" dc:"完成任务总数"`
	List  []CompletedTaskDetail `json:"list" dc:"完成任务明细"`
}

type AnalyticsCompareMember struct {
	KidId  uint64           `json:"kidId" dc:"儿童成员ID"`
	Name   string           `json:"name" dc:"儿童名称"`
	Avatar string           `json:"avatar" dc:"头像"`
	Total  int              `json:"total" dc:"总数"`
	Daily  []AnalyticsPoint `json:"daily" dc:"按天统计"`
	Hourly []AnalyticsPoint `json:"hourly" dc:"按小时统计"`
}

type AnalyticsCompareReq struct {
	g.Meta       `path:"/analytics/compare" method:"get" tags:"儿童端数据" summary:"查询两个成员数据对比"`
	KidId        uint64 `p:"kidId" v:"required|min:1" dc:"当前儿童成员ID"`
	CompareKidId uint64 `p:"compareKidId" v:"required|min:1" dc:"对比儿童成员ID"`
	Metric       string `p:"metric" v:"required|in:tasks,stars" dc:"指标：tasks任务完成数、stars星星数量"`
	Range        string `p:"range" v:"in:weekly,monthly,custom" dc:"范围：weekly、monthly、custom"`
	From         string `p:"from" dc:"开始日期，格式YYYY-MM-DD"`
	To           string `p:"to" dc:"结束日期，格式YYYY-MM-DD"`
	BaseDate     string `p:"baseDate" dc:"周/月基准日期，格式YYYY-MM-DD"`
}

type AnalyticsCompareInput struct {
	KidId        uint64
	CompareKidId uint64
	Metric       string
	Range        string
	From         string
	To           string
	BaseDate     string
}

type AnalyticsCompareOutput struct {
	Metric string
	Range  string
	From   string
	To     string
	Left   AnalyticsCompareMember
	Right  AnalyticsCompareMember
}

type AnalyticsCompareRes struct {
	Metric string                 `json:"metric" dc:"指标"`
	Range  string                 `json:"range" dc:"范围"`
	From   string                 `json:"from" dc:"开始日期"`
	To     string                 `json:"to" dc:"结束日期"`
	Left   AnalyticsCompareMember `json:"left" dc:"当前成员数据"`
	Right  AnalyticsCompareMember `json:"right" dc:"对比成员数据"`
}
