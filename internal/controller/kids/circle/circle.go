package circle

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

// CircleController 复用圈子领域的接口路由声明，避免建立重复的圈子控制器目录。
type CircleController struct{}

// New 创建根路径下的圈子路由控制器。
func New() *CircleController { return &CircleController{} }

// List 查询当前用户管理和加入的群组。
func (c *CircleController) List(ctx context.Context, req *v1.CircleListReq) (res *v1.CircleListRes, err error) {
	out, err := service.Kids().ListCircles(ctx, v1.CircleListInput{UserId: authUserId(ctx)})
	if err != nil {
		return nil, err
	}
	return &v1.CircleListRes{Managed: out.Managed, Joined: out.Joined}, nil
}

// Create 创建圈子并把当前用户设置为管理员。
func (c *CircleController) Create(ctx context.Context, req *v1.CircleCreateReq) (res *v1.CircleCreateRes, err error) {
	out, err := service.Kids().CreateCircle(ctx, v1.CircleCreateInput{UserId: authUserId(ctx), Name: req.Name, Icon: req.Icon})
	if err != nil {
		return nil, err
	}
	return &v1.CircleCreateRes{Circle: out.Circle}, nil
}

// Update 更新群组名称和图标。
func (c *CircleController) Update(ctx context.Context, req *v1.CircleUpdateReq) (res *v1.CircleUpdateRes, err error) {
	out, err := service.Kids().UpdateCircle(ctx, v1.CircleUpdateInput{UserId: authUserId(ctx), Id: req.Id, Name: req.Name, Icon: req.Icon})
	if err != nil {
		return nil, err
	}
	return &v1.CircleUpdateRes{Circle: out.Circle}, nil
}

// Delete 软删除当前用户拥有的群组。
func (c *CircleController) Delete(ctx context.Context, req *v1.CircleDeleteReq) (res *v1.CircleDeleteRes, err error) {
	out, err := service.Kids().DeleteCircle(ctx, v1.CircleDeleteInput{UserId: authUserId(ctx), Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.CircleDeleteRes{Id: out.Id}, nil
}

// Leave 退出非自己创建的群组。
func (c *CircleController) Leave(ctx context.Context, req *v1.CircleLeaveReq) (res *v1.CircleLeaveRes, err error) {
	out, err := service.Kids().LeaveCircle(ctx, v1.CircleLeaveInput{UserId: authUserId(ctx), Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.CircleLeaveRes{Id: out.Id}, nil
}

// Members 查询群组管理员和成员绑定状态。
func (c *CircleController) Members(ctx context.Context, req *v1.CircleMemberListReq) (res *v1.CircleMemberListRes, err error) {
	out, err := service.Kids().ListCircleMembers(ctx, v1.CircleMemberListInput{UserId: authUserId(ctx), CircleId: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.CircleMemberListRes{Owner: out.Owner, Managers: out.Managers, Members: out.Members}, nil
}

// RemoveAdmin 移除群组管理员。
func (c *CircleController) RemoveAdmin(ctx context.Context, req *v1.CircleAdminRemoveReq) (res *v1.CircleAdminRemoveRes, err error) {
	out, err := service.Kids().RemoveCircleAdmin(ctx, v1.CircleAdminRemoveInput{OperatorUserId: authUserId(ctx), CircleId: req.Id, AdminUserId: req.UserId})
	if err != nil {
		return nil, err
	}
	return &v1.CircleAdminRemoveRes{CircleId: out.CircleId, AdminUserId: out.AdminUserId}, nil
}

// CreateInviteCode 为指定圈子创建六位邀请码。
func (c *CircleController) CreateInviteCode(ctx context.Context, req *v1.InviteCodeCreateReq) (res *v1.InviteCodeCreateRes, err error) {
	out, err := service.Kids().CreateInviteCode(ctx, v1.InviteCodeCreateInput{UserId: authUserId(ctx), CircleId: req.CircleId, InviteRole: req.InviteRole, TargetMemberId: req.TargetMemberId, ExpireHours: req.ExpireHours})
	if err != nil {
		return nil, err
	}
	return &v1.InviteCodeCreateRes{Code: out.Code, ExpiredAt: out.ExpiredAt}, nil
}

// PreviewInviteCode 查询邀请码展示信息。
func (c *CircleController) PreviewInviteCode(ctx context.Context, req *v1.InviteCodePreviewReq) (res *v1.InviteCodePreviewRes, err error) {
	out, err := service.Kids().PreviewInviteCode(ctx, v1.InviteCodePreviewInput{Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &v1.InviteCodePreviewRes{Code: out.Code, CircleId: out.CircleId, CircleName: out.CircleName, InviteRole: out.InviteRole, ExpiredAt: out.ExpiredAt}, nil
}

// Join 使用邀请码加入圈子。
func (c *CircleController) Join(ctx context.Context, req *v1.CircleJoinReq) (res *v1.CircleJoinRes, err error) {
	out, err := service.Kids().JoinCircle(ctx, v1.CircleJoinInput{UserId: authUserId(ctx), Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &v1.CircleJoinRes{Circle: out.Circle}, nil
}

// authUserId 从 JWT 鉴权上下文中读取用户ID。
func authUserId(ctx context.Context) uint64 {
	return toUint64(ctx.Value(consts.CtxUserIdKey))
}

// toUint64 将上下文值转换为 uint64。
func toUint64(value any) uint64 {
	switch v := value.(type) {
	case uint64:
		return v
	case int:
		return uint64(v)
	case int64:
		return uint64(v)
	default:
		return 0
	}
}
