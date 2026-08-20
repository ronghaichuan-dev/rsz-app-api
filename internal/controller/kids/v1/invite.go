package v1

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

// CreateInvite 创建邀请。
func (c *Controller) CreateInvite(ctx context.Context, req *v1.CreateInviteReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "createInvite", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CreateCircleInvite(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RefreshInvite 刷新邀请。
func (c *Controller) RefreshInvite(ctx context.Context, req *v1.RefreshInviteReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "refreshInvite", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().RefreshCircleInvite(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RevokeInvite 撤销邀请。
func (c *Controller) RevokeInvite(ctx context.Context, req *v1.RevokeInviteReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "revokeInvite", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().RevokeCircleInvite(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RedeemAdministratorInvite 兑换管理员邀请。
func (c *Controller) RedeemAdministratorInvite(ctx context.Context, req *v1.RedeemAdministratorInviteReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "redeemAdministratorInvite", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().RedeemAdministratorInvite(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RedeemMemberInvite 兑换成员邀请。
func (c *Controller) RedeemMemberInvite(ctx context.Context, req *v1.RedeemMemberInviteReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "redeemMemberInvite", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().RedeemMemberInvite(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}
