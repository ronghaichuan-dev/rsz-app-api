package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	FamilyMemberGenderMale   = "male"
	FamilyMemberGenderFemale = "female"
)

type FamilyMember struct {
	Id          uint64 `json:"id" dc:"家庭成员ID"`
	CircleId    uint64 `json:"circleId" dc:"圈子ID"`
	Name        string `json:"name" dc:"显示名称"`
	Gender      string `json:"gender" v:"in:male,female" dc:"性别：male男、female女"`
	Avatar      string `json:"avatar" dc:"头像地址或预设标识"`
	AvatarStyle string `json:"avatarStyle" dc:"虚拟形象风格标识"`
	Relation    string `json:"relation" dc:"家庭关系，例如爸爸/妈妈/孩子"`
	Owner       bool   `json:"owner" dc:"该成人是否家庭拥有者"`
	BindUserId  uint64 `json:"bindUserId" dc:"绑定用户ID"`
	Bound       bool   `json:"bound" dc:"是否已绑定登录用户"`
	BoundAt     int64  `json:"boundAt" dc:"绑定时间戳"`
	StarCount   int    `json:"starCount" dc:"儿童成员当前星星余额"`
}

type FamilyMemberListReq struct {
	g.Meta   `path:"/family/members" method:"get" tags:"儿童端家庭" summary:"查询家庭成员列表"`
	CircleId uint64 `p:"circleId" dc:"圈子ID"`
}

type FamilyMemberListInput struct {
	CircleId uint64
}

type FamilyMemberListOutput struct {
	Members []FamilyMember
}

type FamilyMemberListRes struct {
	Members []FamilyMember `json:"members" dc:"家庭成员列表"`
}

type FamilyMemberCreateReq struct {
	g.Meta      `path:"/family/members" method:"post" tags:"儿童端家庭" summary:"创建家庭成员"`
	CircleId    uint64 `json:"circleId" v:"required|min:1" dc:"圈子ID"`
	Name        string `json:"name" v:"required|max-length:64" dc:"显示名称"`
	Gender      string `json:"gender" v:"in:male,female" dc:"性别：male男、female女"`
	Avatar      string `json:"avatar" dc:"头像地址或预设标识"`
	AvatarStyle string `json:"avatarStyle" dc:"虚拟形象风格标识"`
	Relation    string `json:"relation" dc:"家庭关系"`
	Owner       bool   `json:"owner" dc:"该成人是否家庭拥有者"`
}

type FamilyMemberCreateInput struct {
	CircleId    uint64
	Name        string
	Gender      string
	Avatar      string
	AvatarStyle string
	Relation    string
	Owner       bool
}

type FamilyMemberCreateOutput struct {
	Member FamilyMember
}

type FamilyMemberCreateRes struct {
	Member FamilyMember `json:"member" dc:"创建的家庭成员"`
}

type FamilyMemberDetailReq struct {
	g.Meta `path:"/family/members/{id}" method:"get" tags:"儿童端家庭" summary:"查询家庭成员详情"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"家庭成员ID"`
}

type FamilyMemberDetailInput struct {
	UserId uint64
	Id     uint64
}

type FamilyMemberDetailOutput struct {
	Member FamilyMember
}

type FamilyMemberDetailRes struct {
	Member FamilyMember `json:"member" dc:"家庭成员详情"`
}

type FamilyMemberUpdateReq struct {
	g.Meta      `path:"/family/members/{id}" method:"put" tags:"儿童端家庭" summary:"编辑家庭成员"`
	Id          uint64 `p:"id" v:"required|min:1" dc:"家庭成员ID"`
	Name        string `json:"name" v:"required|max-length:64" dc:"显示名称"`
	Gender      string `json:"gender" dc:"性别：male男、female女"`
	Avatar      string `json:"avatar" dc:"头像地址或预设标识"`
	AvatarStyle string `json:"avatarStyle" dc:"虚拟形象风格标识"`
	Relation    string `json:"relation" dc:"家庭关系"`
	Owner       bool   `json:"owner" dc:"该成人是否家庭拥有者"`
}

type FamilyMemberUpdateInput struct {
	UserId      uint64
	Id          uint64
	Name        string
	Gender      string
	Avatar      string
	AvatarStyle string
	Relation    string
	Owner       bool
}

type FamilyMemberUpdateOutput struct {
	Member FamilyMember
}

type FamilyMemberUpdateRes struct {
	Member FamilyMember `json:"member" dc:"更新后的家庭成员"`
}

type FamilyMemberDeleteReq struct {
	g.Meta `path:"/family/members/{id}" method:"delete" tags:"儿童端家庭" summary:"删除家庭成员"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"家庭成员ID"`
}

type FamilyMemberDeleteInput struct {
	UserId uint64
	Id     uint64
}

type FamilyMemberDeleteOutput struct {
	Id uint64
}

type FamilyMemberDeleteRes struct {
	Id uint64 `json:"id" dc:"已删除家庭成员ID"`
}
