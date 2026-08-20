package kids

import (
	"context"
	v1 "rslytics-app-api/internal/api/admin/v1"
)

type UserController struct {
}

func New() *UserController {
	return &UserController{}
}

func (c *UserController) List(ctx context.Context, req *v1.UserGetReq) (res *v1.UserGetRes, err error) {
	return nil, nil
}
