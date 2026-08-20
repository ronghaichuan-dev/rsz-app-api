package profile

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Get(ctx context.Context, req *v1.ProfileGetReq) (res *v1.ProfileGetRes, err error) {
	out, err := service.Kids().GetProfile(ctx, v1.ProfileGetInput{UserId: profileUserId(ctx, req.UserId)})
	if err != nil {
		return nil, err
	}
	return &v1.ProfileGetRes{Id: out.Id, Nickname: out.Nickname, Avatar: out.Avatar, Provider: out.Provider, IsGuest: out.IsGuest, Stars: out.Stars}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.ProfileUpdateReq) (res *v1.ProfileUpdateRes, err error) {
	out, err := service.Kids().UpdateProfile(ctx, v1.ProfileUpdateInput{UserId: profileUserId(ctx, 0), Nickname: req.Nickname, Avatar: req.Avatar})
	if err != nil {
		return nil, err
	}
	profile := out.Profile
	return &v1.ProfileUpdateRes{Profile: v1.ProfileGetRes{Id: profile.Id, Nickname: profile.Nickname, Avatar: profile.Avatar, Provider: profile.Provider, IsGuest: profile.IsGuest, Stars: profile.Stars}}, nil
}

// profileUserId 优先使用请求参数，否则使用 JWT 上下文用户ID。
func profileUserId(ctx context.Context, fallback uint64) uint64 {
	if fallback > 0 {
		return fallback
	}
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}
