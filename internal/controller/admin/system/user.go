package system

import (
	"context"

	v1 "rslytics-app-api/internal/api/admin/v1"
	"rslytics-app-api/internal/service"
)

type UserController struct{}

func New() *UserController {
	return &UserController{}
}

func (c *UserController) Login(ctx context.Context, req *v1.SystemUserLoginReq) (res *v1.SystemUserLoginRes, err error) {
	out, err := service.Admin().Login(ctx, v1.SystemUserLoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SystemUserLoginRes{
		UserId: out.UserId,
		Token:  out.Token,
		Name:   out.Name,
	}, nil
}
