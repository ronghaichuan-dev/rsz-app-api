package star

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/service"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Balance(ctx context.Context, req *v1.StarBalanceReq) (res *v1.StarBalanceRes, err error) {
	out, err := service.Kids().GetStarBalance(ctx, v1.StarBalanceInput{KidId: req.KidId})
	if err != nil {
		return nil, err
	}
	return &v1.StarBalanceRes{KidId: out.KidId, Balance: out.Balance}, nil
}

func (c *Controller) Adjust(ctx context.Context, req *v1.StarAdjustReq) (res *v1.StarAdjustRes, err error) {
	out, err := service.Kids().AdjustStars(ctx, v1.StarAdjustInput{KidId: req.KidId, Amount: req.Amount, Reason: req.Reason})
	if err != nil {
		return nil, err
	}
	return &v1.StarAdjustRes{KidId: out.KidId, Balance: out.Balance, Record: out.Record}, nil
}

func (c *Controller) Records(ctx context.Context, req *v1.StarRecordListReq) (res *v1.StarRecordListRes, err error) {
	out, err := service.Kids().ListStarRecords(ctx, v1.StarRecordListInput{KidId: req.KidId, Type: req.Type, From: req.From, To: req.To})
	if err != nil {
		return nil, err
	}
	return &v1.StarRecordListRes{List: out.List}, nil
}
