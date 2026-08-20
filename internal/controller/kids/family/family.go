package family

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

// Controller 处理家庭成员相关接口。
type Controller struct{}

// New 创建家庭成员控制器实例。
func New() *Controller {
	return &Controller{}
}

// List 按圈子查询家庭成员列表。
func (c *Controller) List(ctx context.Context, req *v1.FamilyMemberListReq) (res *v1.FamilyMemberListRes, err error) {
	out, err := service.Kids().ListFamilyMembers(ctx, v1.FamilyMemberListInput{CircleId: req.CircleId})
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberListRes{Members: out.Members}, nil
}

// Create 创建家庭成员，并保存成员虚拟形象配置。
func (c *Controller) Create(ctx context.Context, req *v1.FamilyMemberCreateReq) (res *v1.FamilyMemberCreateRes, err error) {
	out, err := service.Kids().CreateFamilyMember(ctx, v1.FamilyMemberCreateInput{CircleId: req.CircleId, Name: req.Name, Gender: req.Gender, Avatar: req.Avatar, AvatarStyle: req.AvatarStyle, Relation: req.Relation, Owner: req.Owner})
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberCreateRes{Member: out.Member}, nil
}

// Detail 查询家庭成员详情。
func (c *Controller) Detail(ctx context.Context, req *v1.FamilyMemberDetailReq) (res *v1.FamilyMemberDetailRes, err error) {
	out, err := service.Kids().GetFamilyMember(ctx, v1.FamilyMemberDetailInput{UserId: authUserId(ctx), Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberDetailRes{Member: out.Member}, nil
}

// Update 更新家庭成员资料和虚拟形象。
func (c *Controller) Update(ctx context.Context, req *v1.FamilyMemberUpdateReq) (res *v1.FamilyMemberUpdateRes, err error) {
	out, err := service.Kids().UpdateFamilyMember(ctx, v1.FamilyMemberUpdateInput{UserId: authUserId(ctx), Id: req.Id, Name: req.Name, Gender: req.Gender, Avatar: req.Avatar, AvatarStyle: req.AvatarStyle, Relation: req.Relation, Owner: req.Owner})
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberUpdateRes{Member: out.Member}, nil
}

// Delete 软删除家庭成员。
func (c *Controller) Delete(ctx context.Context, req *v1.FamilyMemberDeleteReq) (res *v1.FamilyMemberDeleteRes, err error) {
	out, err := service.Kids().DeleteFamilyMember(ctx, v1.FamilyMemberDeleteInput{UserId: authUserId(ctx), Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberDeleteRes{Id: out.Id}, nil
}

// authUserId 从 JWT 鉴权上下文中读取用户ID。
func authUserId(ctx context.Context) uint64 {
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}
