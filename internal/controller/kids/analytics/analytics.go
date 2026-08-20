package analytics

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	v1controller "rslytics-app-api/internal/controller/kids/v1"
	"rslytics-app-api/internal/service"
)

// Controller 处理 kids 数据统计和对比接口。
type Controller struct{}

// New 创建数据统计控制器实例。
func New() *Controller {
	return &Controller{}
}

// V1Controller 处理 Clearwave v1 的统计领域路由。
type V1Controller struct{}

// NewV1 创建v1统计控制器实例。
func NewV1() *V1Controller {
	return &V1Controller{}
}

// Summary 查询单成员任务或星星统计汇总。
func (c *Controller) Summary(ctx context.Context, req *v1.AnalyticsSummaryReq) (res *v1.AnalyticsSummaryRes, err error) {
	out, err := service.Kids().GetAnalyticsSummary(ctx, v1.AnalyticsSummaryInput{KidId: req.KidId, Metric: req.Metric, Range: req.Range, From: req.From, To: req.To, BaseDate: req.BaseDate})
	if err != nil {
		return nil, err
	}
	return &v1.AnalyticsSummaryRes{KidId: out.KidId, Metric: out.Metric, Range: out.Range, From: out.From, To: out.To, Total: out.Total, Daily: out.Daily, Hourly: out.Hourly}, nil
}

// CompletedTasks 查询已完成任务明细。
func (c *Controller) CompletedTasks(ctx context.Context, req *v1.CompletedTaskListReq) (res *v1.CompletedTaskListRes, err error) {
	out, err := service.Kids().ListCompletedTaskDetails(ctx, v1.CompletedTaskListInput{KidId: req.KidId, From: req.From, To: req.To, Range: req.Range, BaseDate: req.BaseDate})
	if err != nil {
		return nil, err
	}
	return &v1.CompletedTaskListRes{From: out.From, To: out.To, Total: out.Total, List: out.List}, nil
}

// Compare 查询两个成员的数据对比。
func (c *Controller) Compare(ctx context.Context, req *v1.AnalyticsCompareReq) (res *v1.AnalyticsCompareRes, err error) {
	out, err := service.Kids().CompareAnalytics(ctx, v1.AnalyticsCompareInput{KidId: req.KidId, CompareKidId: req.CompareKidId, Metric: req.Metric, Range: req.Range, From: req.From, To: req.To, BaseDate: req.BaseDate})
	if err != nil {
		return nil, err
	}
	return &v1.AnalyticsCompareRes{Metric: out.Metric, Range: out.Range, From: out.From, To: out.To, Left: out.Left, Right: out.Right}, nil
}

// GetStatistics 返回 Clearwave v1 要求的成员统计序列。
func (c *V1Controller) GetStatistics(ctx context.Context, req *v1.GetStatisticsReq) (*v1.V1Response, error) {
	in, err := v1controller.RequestInput(ctx, "getStatistics", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().GetStatisticsV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return v1controller.SuccessResponse(in, out), nil
}

// CompareStatistics 返回 Clearwave v1 要求的成员统计对比序列。
func (c *V1Controller) CompareStatistics(ctx context.Context, req *v1.CompareStatisticsReq) (*v1.V1Response, error) {
	in, err := v1controller.RequestInput(ctx, "compareStatistics", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CompareStatisticsV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return v1controller.SuccessResponse(in, out), nil
}
