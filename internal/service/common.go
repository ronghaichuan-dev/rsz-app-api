// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	commonlogin "rslytics-app-api/internal/common/login"
)

type (
	ICommon interface {
		// Login 根据登录方式完成通用登录，并把用户、授权身份和 token 持久化到调用方数据库。
		Login(ctx context.Context, cfg commonlogin.Config, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error)
		// LoginGuest 创建或复用同设备游客账号，并允许调用方在同一事务内追加业务处理。
		LoginGuest(ctx context.Context, cfg commonlogin.Config, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error)
	}
)

var (
	localCommon ICommon
)

func Common() ICommon {
	if localCommon == nil {
		panic("implement not found for interface ICommon, forgot register?")
	}
	return localCommon
}

func RegisterCommon(i ICommon) {
	localCommon = i
}
