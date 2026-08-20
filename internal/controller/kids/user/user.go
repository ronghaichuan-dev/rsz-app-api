package user

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	out, err := service.Kids().Login(ctx, v1.UserLoginInput{
		Provider:      req.Provider,
		DeviceId:      req.DeviceId,
		AuthCode:      req.AuthCode,
		IdentityToken: req.IdentityToken,
		OpenId:        req.OpenId,
		Email:         req.Email,
		Nickname:      req.Nickname,
		Avatar:        req.Avatar,
		InviteCode:    req.InviteCode,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserLoginRes{
		UserId:         out.UserId,
		Token:          out.Token,
		Provider:       out.Provider,
		IsGuest:        out.IsGuest,
		BoundGuest:     out.BoundGuest,
		IsNewUser:      out.IsNewUser,
		DeviceId:       out.DeviceId,
		Nickname:       out.Nickname,
		Avatar:         out.Avatar,
		AccessExpire:   out.AccessExpire,
		HasCircle:      out.HasCircle,
		CircleId:       out.CircleId,
		CircleRole:     out.CircleRole,
		JoinedByInvite: out.JoinedByInvite,
	}, nil
}
