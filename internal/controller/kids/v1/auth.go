package v1

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

// GetCurrentAccount 获取当前账号 bootstrap。
func (c *Controller) GetAccountBootstrap(ctx context.Context, req *v1.GetCurrentAccountReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "getCurrentAccount", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().GetAccountBootstrap(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// CreateInviteGuestSession 创建邀请码游客会话。
func (c *Controller) CreateGuestSession(ctx context.Context, req *v1.CreateInviteGuestSessionReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "createInviteGuestSession", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().CreateGuestSession(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// ExchangeGoogleProof 交换 Google proof。
func (c *Controller) ExchangeGoogleProof(ctx context.Context, req *v1.ExchangeGoogleProofReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "exchangeGoogleProof", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RefreshSession 刷新会话。
func (c *Controller) RefreshV1Session(ctx context.Context, req *v1.RefreshSessionReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "refreshSession", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().RefreshV1Session(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// RevokeSession 撤销会话。
func (c *Controller) RevokeSession(ctx context.Context, req *v1.RevokeSessionReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "revokeSession", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().RevokeSession(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// GetCurrentEntitlement 获取当前权益。
func (c *Controller) GetEntitlement(ctx context.Context, req *v1.GetCurrentEntitlementReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "getCurrentEntitlement", "GET")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().GetEntitlement(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}

// VerifyPlayPurchase 验证 Play purchase。
func (c *Controller) VerifyPlayPurchase(ctx context.Context, req *v1.VerifyPlayPurchaseReq) (*v1.V1Response, error) {
	in, err := RequestInput(ctx, "verifyPlayPurchase", "POST")
	if err != nil {
		return nil, err
	}
	out, err := service.Kids().ExecuteV1(ctx, in)
	if err != nil {
		return nil, err
	}
	return SuccessResponse(in, out), nil
}
