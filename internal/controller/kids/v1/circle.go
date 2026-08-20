package v1

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

// SelectCircle 选择当前合同圈子。
func (c *Controller) SelectCircle(ctx context.Context, req *v1.SelectCurrentCircleReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "selectCurrentCircle", "PUT")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().SelectCircle(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// ListV1Circles 查询当前账号的合同圈子列表。
func (c *Controller) ListV1Circles(ctx context.Context, req *v1.ListMyCirclesReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "listMyCircles", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ListV1Circles(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// DeleteCircle 删除合同圈子。
func (c *Controller) DeleteCircle(ctx context.Context, req *v1.DeleteCircleReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "deleteCircle", "DELETE")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().DeleteCircleWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// UpdateCircle 更新合同圈子。
func (c *Controller) UpdateCircle(ctx context.Context, req *v1.UpdateCircleReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "updateCircle", "PATCH")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().UpdateCircleWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// DeleteAdministrator 删除合同管理员。
func (c *Controller) DeleteAdministrator(ctx context.Context, req *v1.DeleteAdministratorReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "deleteAdministrator", "DELETE")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().DeleteCircleAdministrator(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// UpsertAdministrator 新增或更新合同管理员。
func (c *Controller) UpsertAdministrator(ctx context.Context, req *v1.UpsertAdministratorReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "upsertAdministrator", "PUT")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().UpsertCircleAdministrator(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// GetCircleBootstrap 获取合同圈子 bootstrap。
func (c *Controller) GetCircleBootstrap(ctx context.Context, req *v1.GetCircleBootstrapReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "getCircleBootstrap", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().GetV1CircleBootstrap(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// CreateCircleMember 创建合同成员。
func (c *Controller) CreateCircleMember(ctx context.Context, req *v1.CreateCircleMemberReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "createCircleMember", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CreateCircleMember(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// DeleteMember 删除合同成员。
func (c *Controller) DeleteMember(ctx context.Context, req *v1.DeleteMemberReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "deleteMember", "DELETE")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().DeleteCircleMember(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// UpsertMember 更新合同成员。
func (c *Controller) UpsertMember(ctx context.Context, req *v1.UpsertMemberReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "upsertMember", "PUT")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().UpsertCircleMember(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// GetMemberBalances 获取成员的合同星星余额。
func (c *Controller) GetMemberBalances(ctx context.Context, req *v1.GetMemberBalancesReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "getMemberBalances", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().GetMemberBalancesV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// LeaveCircle 退出合同圈子。
func (c *Controller) LeaveCircle(ctx context.Context, req *v1.LeaveCircleReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "leaveCircle", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().LeaveCircleWithVersion(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// CompleteOnboarding 完成合同 onboarding。
func (c *Controller) CompleteOnboarding(ctx context.Context, req *v1.CompleteOnboardingReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "completeOnboarding", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CompleteOnboardingV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}
