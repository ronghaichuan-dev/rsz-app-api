package v1

import "github.com/gogf/gf/v2/frame/g"

type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"儿童端" summary:"儿童端健康检查"`
}

type HealthRes struct {
	App    string `json:"app" dc:"应用名称"`
	Status string `json:"status" dc:"服务状态"`
}
