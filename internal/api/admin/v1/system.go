package v1

import "github.com/gogf/gf/v2/frame/g"

type SystemUserLoginReq struct {
	g.Meta   `path:"/users/login" method:"post" tags:"管理端系统" summary:"管理端系统用户登录"`
	Username string `json:"username" v:"required" dc:"用户名"`
	Password string `json:"password" v:"required" dc:"密码"`
}

type SystemUserLoginInput struct {
	Username string
	Password string
}

type SystemUserLoginOutput struct {
	UserId uint64
	Token  string
	Name   string
}

type SystemUserLoginRes struct {
	UserId uint64 `json:"userId" dc:"管理端用户标识"`
	Token  string `json:"token" dc:"登录访问令牌"`
	Name   string `json:"name" dc:"显示名称"`
}
