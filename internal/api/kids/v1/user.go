package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"rslytics-app-api/internal/consts"
)

const (
	LoginProviderGuest  = consts.LoginProviderGuest
	LoginProviderGoogle = consts.LoginProviderGoogle
	LoginProviderApple  = consts.LoginProviderApple
	LoginProviderInvite = consts.LoginProviderInvite
)

type UserLoginReq struct {
	g.Meta        `path:"/users/login" method:"post" tags:"登录" summary:"游客、谷歌、苹果或邀请码登录" description:"游客登录提交 deviceId；Google/Apple 授权登录必须提交 identityToken，后端会校验服务商签名、签发方、客户端ID和过期时间，并以校验后的 sub 作为唯一授权身份；邀请码登录提交 inviteCode，后端创建或复用同设备游客账号并加入邀请码对应圈子。"`
	Provider      string `json:"provider" v:"required|in:guest,google,apple,invite" dc:"登录方式：guest游客、google谷歌、apple苹果、invite邀请码"`
	DeviceId      string `json:"deviceId" v:"required" dc:"唯一设备ID"`
	AuthCode      string `json:"authCode" dc:"授权码，当前登录以identityToken服务端校验为准"`
	IdentityToken string `json:"identityToken" v:"required-if:provider,google,provider,apple#identityToken is required" dc:"谷歌或苹果身份令牌，授权登录必填"`
	OpenId        string `json:"openId" dc:"兼容字段，授权登录不信任客户端传入值"`
	Email         string `json:"email" dc:"兼容字段，授权登录以identityToken校验结果为准"`
	Nickname      string `json:"nickname" dc:"游客昵称，授权登录以identityToken校验结果优先"`
	Avatar        string `json:"avatar" dc:"游客头像，授权登录以identityToken校验结果优先"`
	InviteCode    string `json:"inviteCode" v:"required-if:provider,invite#inviteCode is required|length:6" dc:"六位邀请码，邀请码登录必填"`
}

type UserLoginInput struct {
	Provider      string
	DeviceId      string
	AuthCode      string
	IdentityToken string
	OpenId        string
	Email         string
	Nickname      string
	Avatar        string
	InviteCode    string
}

type UserLoginOutput struct {
	UserId         uint64
	Token          string
	Provider       string
	IsGuest        bool
	BoundGuest     bool
	IsNewUser      bool
	DeviceId       string
	Nickname       string
	Avatar         string
	AccessExpire   int64
	HasCircle      bool
	CircleId       uint64
	CircleRole     string
	JoinedByInvite bool
}

type UserLoginRes struct {
	UserId         uint64 `json:"userId" dc:"儿童端用户ID"`
	Token          string `json:"token" dc:"登录访问令牌"`
	Provider       string `json:"provider" dc:"当前登录方式"`
	IsGuest        bool   `json:"isGuest" dc:"当前用户是否游客"`
	BoundGuest     bool   `json:"boundGuest" dc:"本次授权登录是否绑定已有游客账号"`
	IsNewUser      bool   `json:"isNewUser" dc:"是否创建新用户"`
	DeviceId       string `json:"deviceId" dc:"设备ID"`
	Nickname       string `json:"nickname" dc:"昵称"`
	Avatar         string `json:"avatar" dc:"头像地址"`
	AccessExpire   int64  `json:"accessExpire" dc:"登录访问令牌过期时间戳"`
	HasCircle      bool   `json:"hasCircle" dc:"用户是否已加入圈子"`
	CircleId       uint64 `json:"circleId" dc:"当前圈子ID"`
	CircleRole     string `json:"circleRole" dc:"当前圈子角色"`
	JoinedByInvite bool   `json:"joinedByInvite" dc:"本次登录是否通过邀请码加入圈子"`
}
