package v1

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

// CommitAssetUpload 提交资产上传。
func (c *Controller) CommitAssetUpload(ctx context.Context, req *v1.CommitAssetUploadReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "commitAssetUpload", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// PrepareAssetUpload 准备资产上传。
func (c *Controller) PrepareAssetUpload(ctx context.Context, req *v1.PrepareAssetUploadReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "prepareAssetUpload", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// ListExchangeHistory 列出兑换历史。
func (c *Controller) ListExchangeHistory(ctx context.Context, req *v1.ListExchangeHistoryReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "listExchangeHistory", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RedeemReward 兑换奖励。
func (c *Controller) RedeemReward(ctx context.Context, req *v1.RedeemRewardReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "redeemReward", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// DeleteReward 删除奖励。
func (c *Controller) DeleteReward(ctx context.Context, req *v1.DeleteRewardReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "deleteReward", "DELETE")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// UpsertReward 更新奖励。
func (c *Controller) UpsertReward(ctx context.Context, req *v1.UpsertRewardReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "upsertReward", "PUT")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// GetRewardEligibility 获取奖励资格。
func (c *Controller) GetRewardEligibility(ctx context.Context, req *v1.GetRewardEligibilityReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "getRewardEligibility", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// ListStarTransactions 列出星星流水。
func (c *Controller) ListStarTransactions(ctx context.Context, req *v1.ListStarTransactionsReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "listStarTransactions", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ListStarTransactionsV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// PullCircleBootstrapDelta 拉取圈子增量。
func (c *Controller) PullCircleBootstrapDelta(ctx context.Context, req *v1.PullCircleBootstrapDeltaReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "pullCircleBootstrapDelta", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// SubmitFeedback 提交反馈。
func (c *Controller) SubmitFeedbackV1(ctx context.Context, req *v1.SubmitFeedbackReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "submitFeedback", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().SubmitFeedbackV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}
