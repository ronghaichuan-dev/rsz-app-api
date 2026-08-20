package v1

import "github.com/gogf/gf/v2/frame/g"

type ProfileGetReq struct {
	g.Meta `path:"/profile" method:"get" tags:"儿童端资料" summary:"获取当前儿童端用户资料"`
	UserId uint64 `p:"userId" dc:"用户标识，默认使用令牌上下文用户标识"`
}

type ProfileGetInput struct {
	UserId uint64
}

type ProfileGetOutput struct {
	Id       uint64
	Nickname string
	Avatar   string
	Provider string
	IsGuest  bool
	Stars    int
}

type ProfileGetRes struct {
	Id       uint64 `json:"id" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像地址"`
	Provider string `json:"provider" dc:"当前登录方式"`
	IsGuest  bool   `json:"isGuest" dc:"当前用户是否游客"`
	Stars    int    `json:"stars" dc:"当前星星余额"`
}

type ProfileUpdateReq struct {
	g.Meta   `path:"/profile" method:"put" tags:"儿童端资料" summary:"更新当前用户资料"`
	Nickname string `json:"nickname" v:"required|max-length:128" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像地址"`
}

type ProfileUpdateInput struct {
	UserId   uint64
	Nickname string
	Avatar   string
}

type ProfileUpdateOutput struct {
	Profile ProfileGetOutput
}

type ProfileUpdateRes struct {
	Profile ProfileGetRes `json:"profile" dc:"更新后的用户资料"`
}
