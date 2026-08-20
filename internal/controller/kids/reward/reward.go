package reward

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

// Controller 处理奖励商城、预设和兑换历史接口。
type Controller struct{}

// New 创建奖励控制器实例。
func New() *Controller {
	return &Controller{}
}

// Presets 查询奖励预设列表。
func (c *Controller) Presets(ctx context.Context, req *v1.RewardPresetListReq) (res *v1.RewardPresetListRes, err error) {
	out, err := service.Kids().ListRewardPresets(ctx, v1.RewardPresetListInput{Keyword: req.Keyword})
	if err != nil {
		return nil, err
	}
	return &v1.RewardPresetListRes{List: out.List}, nil
}

// List 查询奖励列表，支持群组和儿童兑换进度筛选。
func (c *Controller) List(ctx context.Context, req *v1.RewardListReq) (res *v1.RewardListRes, err error) {
	out, err := service.Kids().ListRewards(ctx, v1.RewardListInput{CircleId: req.CircleId, KidId: req.KidId})
	if err != nil {
		return nil, err
	}
	return &v1.RewardListRes{List: out.List}, nil
}

// Detail 查询奖励详情。
func (c *Controller) Detail(ctx context.Context, req *v1.RewardDetailReq) (res *v1.RewardDetailRes, err error) {
	out, err := service.Kids().GetReward(ctx, v1.RewardDetailInput{Id: req.Id, KidId: req.KidId})
	if err != nil {
		return nil, err
	}
	return &v1.RewardDetailRes{Reward: out.Reward}, nil
}

// Create 创建奖励并保存指派儿童和重复兑换规则。
func (c *Controller) Create(ctx context.Context, req *v1.RewardCreateReq) (res *v1.RewardCreateRes, err error) {
	in := v1.RewardCreateInput{UserId: authUserId(ctx), CircleId: req.CircleId, PresetId: req.PresetId, Title: req.Title, Icon: req.Icon, ImageUrl: req.ImageUrl, StarCost: req.StarCost, Stock: req.Stock, Description: req.Description, RepeatRule: req.RepeatRule, RepeatIntervalDays: req.RepeatIntervalDays, KidIds: req.KidIds}
	out, err := service.Kids().CreateReward(ctx, in)
	if err != nil {
		return nil, err
	}
	return &v1.RewardCreateRes{Reward: out.Reward}, nil
}

// Update 编辑奖励配置。
func (c *Controller) Update(ctx context.Context, req *v1.RewardUpdateReq) (res *v1.RewardUpdateRes, err error) {
	in := v1.RewardUpdateInput{UserId: authUserId(ctx), Id: req.Id, CircleId: req.CircleId, Title: req.Title, Icon: req.Icon, ImageUrl: req.ImageUrl, StarCost: req.StarCost, Stock: req.Stock, Description: req.Description, RepeatRule: req.RepeatRule, RepeatIntervalDays: req.RepeatIntervalDays, KidIds: req.KidIds}
	out, err := service.Kids().UpdateReward(ctx, in)
	if err != nil {
		return nil, err
	}
	return &v1.RewardUpdateRes{Reward: out.Reward}, nil
}

// Delete 软删除奖励。
func (c *Controller) Delete(ctx context.Context, req *v1.RewardDeleteReq) (res *v1.RewardDeleteRes, err error) {
	out, err := service.Kids().DeleteReward(ctx, v1.RewardDeleteInput{UserId: authUserId(ctx), Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.RewardDeleteRes{Id: out.Id}, nil
}

// Redeem 兑换奖励。
func (c *Controller) Redeem(ctx context.Context, req *v1.RewardRedeemReq) (res *v1.RewardRedeemRes, err error) {
	out, err := service.Kids().RedeemReward(ctx, v1.RewardRedeemInput{Id: req.Id, KidId: req.KidId, Remark: req.Remark})
	if err != nil {
		return nil, err
	}
	return &v1.RewardRedeemRes{Reward: out.Reward, StarBalance: out.StarBalance}, nil
}

// Records 查询奖励兑换历史。
func (c *Controller) Records(ctx context.Context, req *v1.RewardRecordListReq) (res *v1.RewardRecordListRes, err error) {
	out, err := service.Kids().ListRewardRedeemRecords(ctx, v1.RewardRecordListInput{CircleId: req.CircleId, KidId: req.KidId, Month: req.Month, From: req.From, To: req.To})
	if err != nil {
		return nil, err
	}
	return &v1.RewardRecordListRes{Total: out.Total, List: out.List}, nil
}

// authUserId 从 JWT 鉴权上下文中读取用户ID。
func authUserId(ctx context.Context) uint64 {
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}
