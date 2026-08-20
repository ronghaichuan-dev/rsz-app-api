package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	CircleRoleAdmin  = "admin"
	CircleRoleMember = "member"
)

type CircleInfo struct {
	Id     uint64 `json:"id" dc:"圈子ID"`
	Name   string `json:"name" dc:"圈子名称"`
	Icon   string `json:"icon" dc:"圈子图标"`
	Role   string `json:"role" dc:"当前用户在圈子中的角色"`
	Joined bool   `json:"joined" dc:"是否已加入圈子"`
}

type CircleCreateReq struct {
	g.Meta `path:"/circles/create" method:"post" tags:"儿童端圈子" summary:"创建管理员圈子" description:"登录用户可以创建圈子，创建者会自动成为圈子管理员。"`
	Name   string `json:"name" v:"required|max-length:64" dc:"圈子名称"`
	Icon   string `json:"icon" dc:"圈子图标"`
}

type CircleCreateInput struct {
	UserId uint64
	Name   string
	Icon   string
}

type CircleCreateOutput struct {
	Circle CircleInfo
}

type CircleCreateRes struct {
	Circle CircleInfo `json:"circle" dc:"创建的圈子"`
}

type InviteCodeCreateReq struct {
	g.Meta         `path:"/invite-codes" method:"post" tags:"儿童端圈子" summary:"创建邀请码"`
	CircleId       uint64 `json:"circleId" v:"required|min:1" dc:"圈子ID"`
	InviteRole     string `json:"inviteRole" v:"required|in:admin,member" dc:"邀请加入角色：admin管理员、member成员"`
	TargetMemberId uint64 `json:"targetMemberId" dc:"目标家庭成员ID，成员邀请码可指定"`
	ExpireHours    int    `json:"expireHours" dc:"有效小时数，默认72小时"`
}

type InviteCodeCreateInput struct {
	UserId         uint64
	CircleId       uint64
	InviteRole     string
	TargetMemberId uint64
	ExpireHours    int
}

type InviteCodeCreateOutput struct {
	Code      string
	ExpiredAt int64
}

type InviteCodeCreateRes struct {
	Code      string `json:"code" dc:"六位邀请码"`
	ExpiredAt int64  `json:"expiredAt" dc:"过期时间戳"`
}

type InviteCodePreviewReq struct {
	g.Meta `path:"/invite-codes/preview" method:"get" tags:"儿童端圈子" summary:"预览邀请码"`
	Code   string `p:"code" v:"required|length:6" dc:"六位邀请码"`
}

type InviteCodePreviewInput struct {
	Code string
}

type InviteCodePreviewOutput struct {
	Code       string
	CircleId   uint64
	CircleName string
	InviteRole string
	ExpiredAt  int64
}

type InviteCodePreviewRes struct {
	Code       string `json:"code" dc:"六位邀请码"`
	CircleId   uint64 `json:"circleId" dc:"圈子ID"`
	CircleName string `json:"circleName" dc:"圈子名称"`
	InviteRole string `json:"inviteRole" dc:"邀请加入角色"`
	ExpiredAt  int64  `json:"expiredAt" dc:"过期时间戳"`
}

type CircleJoinReq struct {
	g.Meta `path:"/circles/join" method:"post" tags:"儿童端圈子" summary:"通过邀请码加入圈子"`
	Code   string `json:"code" v:"required|length:6" dc:"六位邀请码"`
}

type CircleJoinInput struct {
	UserId uint64
	Code   string
}

type CircleJoinOutput struct {
	Circle CircleInfo
}

type CircleJoinRes struct {
	Circle CircleInfo `json:"circle" dc:"加入的圈子"`
}

type CircleListReq struct {
	g.Meta `path:"/circles" method:"get" tags:"儿童端圈子" summary:"查询我的群组列表"`
}

type CircleListInput struct {
	UserId uint64
}

type CircleListOutput struct {
	Managed []CircleInfo
	Joined  []CircleInfo
}

type CircleListRes struct {
	Managed []CircleInfo `json:"managed" dc:"我管理的群组"`
	Joined  []CircleInfo `json:"joined" dc:"我加入的群组"`
}

type CircleUpdateReq struct {
	g.Meta `path:"/circles/{id}" method:"put" tags:"儿童端圈子" summary:"更新群组资料"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"圈子ID"`
	Name   string `json:"name" v:"required|max-length:64" dc:"圈子名称"`
	Icon   string `json:"icon" dc:"圈子图标"`
}

type CircleUpdateInput struct {
	UserId uint64
	Id     uint64
	Name   string
	Icon   string
}

type CircleUpdateOutput struct {
	Circle CircleInfo
}

type CircleUpdateRes struct {
	Circle CircleInfo `json:"circle" dc:"更新后的群组"`
}

type CircleDeleteReq struct {
	g.Meta `path:"/circles/{id}" method:"delete" tags:"儿童端圈子" summary:"删除群组"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"圈子ID"`
}

type CircleDeleteInput struct {
	UserId uint64
	Id     uint64
}

type CircleDeleteOutput struct {
	Id uint64
}

type CircleDeleteRes struct {
	Id uint64 `json:"id" dc:"已删除圈子ID"`
}

type CircleLeaveReq struct {
	g.Meta `path:"/circles/{id}/leave" method:"post" tags:"儿童端圈子" summary:"退出群组"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"圈子ID"`
}

type CircleLeaveInput struct {
	UserId uint64
	Id     uint64
}

type CircleLeaveOutput struct {
	Id uint64
}

type CircleLeaveRes struct {
	Id uint64 `json:"id" dc:"已退出圈子ID"`
}

type CircleMemberItem struct {
	UserId    uint64 `json:"userId" dc:"用户ID"`
	MemberId  uint64 `json:"memberId" dc:"家庭成员ID"`
	Name      string `json:"name" dc:"显示名称"`
	Avatar    string `json:"avatar" dc:"头像"`
	Role      string `json:"role" dc:"圈子角色"`
	Bound     bool   `json:"bound" dc:"是否已绑定登录用户"`
	IsOwner   bool   `json:"isOwner" dc:"是否群组所有者"`
	CreatedAt int64  `json:"createdAt" dc:"创建时间戳"`
}

type CircleMemberListReq struct {
	g.Meta `path:"/circles/{id}/members" method:"get" tags:"儿童端圈子" summary:"查询群组成员和管理员"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"圈子ID"`
}

type CircleMemberListInput struct {
	UserId   uint64
	CircleId uint64
}

type CircleMemberListOutput struct {
	Owner    CircleMemberItem
	Managers []CircleMemberItem
	Members  []CircleMemberItem
}

type CircleMemberListRes struct {
	Owner    CircleMemberItem   `json:"owner" dc:"群组所有者"`
	Managers []CircleMemberItem `json:"managers" dc:"管理员列表"`
	Members  []CircleMemberItem `json:"members" dc:"成员列表"`
}

type CircleAdminRemoveReq struct {
	g.Meta `path:"/circles/{id}/admins/{userId}" method:"delete" tags:"儿童端圈子" summary:"移除管理员"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"圈子ID"`
	UserId uint64 `p:"userId" v:"required|min:1" dc:"被移除管理员用户ID"`
}

type CircleAdminRemoveInput struct {
	OperatorUserId uint64
	CircleId       uint64
	AdminUserId    uint64
}

type CircleAdminRemoveOutput struct {
	CircleId    uint64
	AdminUserId uint64
}

type CircleAdminRemoveRes struct {
	CircleId    uint64 `json:"circleId" dc:"圈子ID"`
	AdminUserId uint64 `json:"adminUserId" dc:"被移除管理员用户ID"`
}
