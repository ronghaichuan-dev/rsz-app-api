package notification

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) List(ctx context.Context, req *v1.NotificationListReq) (res *v1.NotificationListRes, err error) {
	out, err := service.Kids().ListNotifications(ctx, v1.NotificationListInput{MemberId: req.MemberId, Unread: req.Unread})
	if err != nil {
		return nil, err
	}
	return &v1.NotificationListRes{List: out.List}, nil
}

func (c *Controller) Read(ctx context.Context, req *v1.NotificationReadReq) (res *v1.NotificationReadRes, err error) {
	out, err := service.Kids().ReadNotification(ctx, v1.NotificationReadInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.NotificationReadRes{Notification: out.Notification}, nil
}
