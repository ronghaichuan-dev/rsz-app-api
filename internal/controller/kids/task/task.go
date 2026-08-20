package task

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

// Controller 处理 kids 任务、任务预设和任务标签相关接口。
type Controller struct{}

// New 创建 kids 任务控制器实例。
func New() *Controller {
	return &Controller{}
}

// Presets 返回可用于快速创建任务的预设列表。
func (c *Controller) Presets(ctx context.Context, req *v1.TaskPresetListReq) (res *v1.TaskPresetListRes, err error) {
	out, err := service.Kids().ListTaskPresets(ctx, v1.TaskPresetListInput{Keyword: req.Keyword})
	if err != nil {
		return nil, err
	}
	return &v1.TaskPresetListRes{List: out.List}, nil
}

// List 按日期、孩子、标签和状态查询任务列表。
func (c *Controller) List(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error) {
	out, err := service.Kids().ListTasks(ctx, v1.TaskListInput{Date: req.Date, KidId: req.KidId, TagId: req.TagId, Status: req.Status})
	if err != nil {
		return nil, err
	}
	return &v1.TaskListRes{Date: out.Date, List: out.List}, nil
}

// Detail 查询任务详情。
func (c *Controller) Detail(ctx context.Context, req *v1.TaskDetailReq) (res *v1.TaskDetailRes, err error) {
	out, err := service.Kids().GetTask(ctx, v1.TaskDetailInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.TaskDetailRes{Task: out.Task}, nil
}

// Create 创建任务并分配给一个或多个儿童。
func (c *Controller) Create(ctx context.Context, req *v1.TaskCreateReq) (res *v1.TaskCreateRes, err error) {
	in := v1.TaskCreateInput{Title: req.Title, Icon: req.Icon, Note: req.Note, Star: req.Star, KidIds: req.KidIds, CompletionMode: req.CompletionMode, RepeatRule: req.RepeatRule, RepeatEndType: req.RepeatEndType, RepeatEndDate: req.RepeatEndDate, RepeatEndCount: req.RepeatEndCount, StartDate: req.StartDate, EndDate: req.EndDate, TimeLimitType: req.TimeLimitType, TimeLimitStart: req.TimeLimitStart, TimeLimitEnd: req.TimeLimitEnd, ReminderType: req.ReminderType, ReminderAt: req.ReminderAt, ReminderOffsetMinutes: req.ReminderOffsetMinutes, NeedPhotoProof: req.NeedPhotoProof, TagId: req.TagId}
	out, err := service.Kids().CreateTask(ctx, in)
	if err != nil {
		return nil, err
	}
	return &v1.TaskCreateRes{Task: out.Task}, nil
}

// Update 编辑未完成任务的基础配置和分配儿童。
func (c *Controller) Update(ctx context.Context, req *v1.TaskUpdateReq) (res *v1.TaskUpdateRes, err error) {
	in := v1.TaskUpdateInput{Id: req.Id, Title: req.Title, Icon: req.Icon, Note: req.Note, Star: req.Star, KidIds: req.KidIds, CompletionMode: req.CompletionMode, RepeatRule: req.RepeatRule, RepeatEndType: req.RepeatEndType, RepeatEndDate: req.RepeatEndDate, RepeatEndCount: req.RepeatEndCount, TimeLimitType: req.TimeLimitType, TimeLimitStart: req.TimeLimitStart, TimeLimitEnd: req.TimeLimitEnd, ReminderType: req.ReminderType, ReminderAt: req.ReminderAt, ReminderOffsetMinutes: req.ReminderOffsetMinutes, NeedPhotoProof: req.NeedPhotoProof, TagId: req.TagId}
	out, err := service.Kids().UpdateTask(ctx, in)
	if err != nil {
		return nil, err
	}
	return &v1.TaskUpdateRes{Task: out.Task}, nil
}

// Delete 软删除任务。
func (c *Controller) Delete(ctx context.Context, req *v1.TaskDeleteReq) (res *v1.TaskDeleteRes, err error) {
	out, err := service.Kids().DeleteTask(ctx, v1.TaskDeleteInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.TaskDeleteRes{Id: out.Id}, nil
}

// Complete 提交任务完成结果，并按任务规则发放星星。
func (c *Controller) Complete(ctx context.Context, req *v1.TaskCompleteReq) (res *v1.TaskCompleteRes, err error) {
	out, err := service.Kids().CompleteTask(ctx, v1.TaskCompleteInput{Id: req.Id, KidId: req.KidId, PhotoUrl: req.PhotoUrl})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCompleteRes{Task: out.Task, StarBalance: out.StarBalance, TaskCompleted: out.TaskCompleted}, nil
}

// Cancel 取消任务完成状态并回滚对应星星。
func (c *Controller) Cancel(ctx context.Context, req *v1.TaskCancelReq) (res *v1.TaskCancelRes, err error) {
	out, err := service.Kids().CancelTask(ctx, v1.TaskCancelInput{Id: req.Id, KidId: req.KidId, Reason: req.Reason})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCancelRes{Task: out.Task, StarBalance: out.StarBalance}, nil
}

// Tags 查询任务标签列表。
func (c *Controller) Tags(ctx context.Context, req *v1.TaskTagListReq) (res *v1.TaskTagListRes, err error) {
	out, err := service.Kids().ListTaskTags(ctx, v1.TaskTagListInput{})
	if err != nil {
		return nil, err
	}
	return &v1.TaskTagListRes{List: out.List}, nil
}

// CreateTag 创建任务标签。
func (c *Controller) CreateTag(ctx context.Context, req *v1.TaskTagCreateReq) (res *v1.TaskTagCreateRes, err error) {
	out, err := service.Kids().CreateTaskTag(ctx, v1.TaskTagCreateInput{Name: req.Name, Color: req.Color, SortOrder: req.SortOrder})
	if err != nil {
		return nil, err
	}
	return &v1.TaskTagCreateRes{Tag: out.Tag}, nil
}

// UpdateTag 更新任务标签。
func (c *Controller) UpdateTag(ctx context.Context, req *v1.TaskTagUpdateReq) (res *v1.TaskTagUpdateRes, err error) {
	out, err := service.Kids().UpdateTaskTag(ctx, v1.TaskTagUpdateInput{Id: req.Id, Name: req.Name, Color: req.Color, SortOrder: req.SortOrder})
	if err != nil {
		return nil, err
	}
	return &v1.TaskTagUpdateRes{Tag: out.Tag}, nil
}

// DeleteTag 删除任务标签。
func (c *Controller) DeleteTag(ctx context.Context, req *v1.TaskTagDeleteReq) (res *v1.TaskTagDeleteRes, err error) {
	out, err := service.Kids().DeleteTaskTag(ctx, v1.TaskTagDeleteInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.TaskTagDeleteRes{Id: out.Id}, nil
}
