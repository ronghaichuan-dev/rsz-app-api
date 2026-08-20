package v1

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

// ListTaskCompletionDetails 列出任务完成明细。
func (c *Controller) ListTaskCompletionDetailsV1(ctx context.Context, req *v1.ListTaskCompletionDetailsReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "listTaskCompletionDetails", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ListTaskCompletionDetailsV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// CompleteTask 完成任务。
func (c *Controller) CompleteTask(ctx context.Context, req *v1.CompleteTaskReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "completeTask", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CompleteTaskWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// CancelTaskCompletion 取消任务完成。
func (c *Controller) CancelTaskCompletion(ctx context.Context, req *v1.CancelTaskCompletionReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "cancelTaskCompletion", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CancelTaskCompletionWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// ListTaskOccurrences 列出任务 occurrence。
func (c *Controller) ListTaskOccurrencesV1(ctx context.Context, req *v1.ListTaskOccurrencesReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "listTaskOccurrences", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ListTaskOccurrencesV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// DeleteTaskTag 删除任务标签。
func (c *Controller) DeleteTaskTag(ctx context.Context, req *v1.DeleteTaskTagReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "deleteTaskTag", "DELETE")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().DeleteTaskTagWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// UpsertTaskTag 更新任务标签。
func (c *Controller) UpsertTaskTag(ctx context.Context, req *v1.UpsertTaskTagReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "upsertTaskTag", "PUT")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().UpsertTaskTagWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// DeleteTask 删除任务。
func (c *Controller) DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "deleteTask", "DELETE")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().DeleteTaskWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// UpsertTask 更新任务。
func (c *Controller) UpsertTask(ctx context.Context, req *v1.UpsertTaskReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "upsertTask", "PUT")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().UpsertTaskWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// AdjustMemberStars 调整成员星星。
func (c *Controller) AdjustMemberStars(ctx context.Context, req *v1.AdjustMemberStarsReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "adjustMemberStars", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().AdjustMemberStarsWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}
