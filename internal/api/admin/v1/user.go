package v1

import "github.com/gogf/gf/v2/frame/g"

type UserGetReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"管理端用户" summary:"获取管理端用户"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"用户ID"`
}

type UserGetRes struct {
	Id       uint64 `json:"id" dc:"用户ID"`
	Username string `json:"username" dc:"用户名"`
	Role     string `json:"role" dc:"角色编码"`
}
