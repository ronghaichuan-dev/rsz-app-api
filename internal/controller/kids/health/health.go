package health

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	return &v1.HealthRes{App: "kids", Status: "ok"}, nil
}
