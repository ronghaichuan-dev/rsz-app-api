package service

import (
	"context"

	commonlogin "rslytics-app-api/internal/common/login"
)

type (
	// ILogin 定义公共登录模块的服务接口。
	ILogin interface {
		// Login 根据登录方式完成公共登录，并持久化用户、授权身份和 token。
		Login(ctx context.Context, cfg commonlogin.Config, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error)
		// LoginGuest 创建或复用同设备游客账号，并允许在同一事务内追加业务处理。
		LoginGuest(ctx context.Context, cfg commonlogin.Config, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error)
	}
)

var (
	localLogin ILogin
)

// Login 返回公共登录模块的服务实现。
func Login() ILogin {
	if localLogin == nil {
		panic("implement not found for interface ILogin, forgot register?")
	}
	return localLogin
}

// RegisterLogin 注册公共登录模块的服务实现。
func RegisterLogin(i ILogin) {
	localLogin = i
}
