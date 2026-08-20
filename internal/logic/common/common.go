package common

import (
	"rslytics-app-api/internal/service"
)

// sLogin 是公共登录模块的 service 实现。
type sCommon struct{}

// init 在包加载时注册公共登录实现到 service 层。
func init() {
	service.RegisterCommon(NewCommon())
}

// NewCommon 创建公共登录服务实例。
func NewCommon() *sCommon {
	return &sCommon{}
}
