package index

import (
	"context"
	v1 "rslytics-app-api/internal/api/kids/v1"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	return &v1.IndexRes{}, nil
}
