// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractSession is the golang structure for table kids_contract_session.
type KidsContractSession struct {
	Id                    uint64      `json:"id"                    orm:"id"                       description:"主键"`       // 主键
	SessionId             string      `json:"sessionId"             orm:"session_id"               description:"合同会话标识"`   // 合同会话标识
	AccountId             string      `json:"accountId"             orm:"account_id"               description:"合同账号标识"`   // 合同账号标识
	PrincipalKind         string      `json:"principalKind"         orm:"principal_kind"           description:"主体类型"`     // 主体类型
	AccessTokenHash       string      `json:"accessTokenHash"       orm:"access_token_hash"        description:"访问令牌摘要"`   // 访问令牌摘要
	RefreshTokenHash      string      `json:"refreshTokenHash"      orm:"refresh_token_hash"       description:"刷新令牌摘要"`   // 刷新令牌摘要
	GuestUpgradeGrantHash string      `json:"guestUpgradeGrantHash" orm:"guest_upgrade_grant_hash" description:"游客升级凭据摘要"` // 游客升级凭据摘要
	ExpiresAt             *gtime.Time `json:"expiresAt"             orm:"expires_at"               description:"访问令牌到期时间"` // 访问令牌到期时间
	RefreshExpiresAt      *gtime.Time `json:"refreshExpiresAt"      orm:"refresh_expires_at"       description:"刷新令牌到期时间"` // 刷新令牌到期时间
	RevokedAt             *gtime.Time `json:"revokedAt"             orm:"revoked_at"               description:"撤销时间"`     // 撤销时间
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"               description:"创建时间"`     // 创建时间
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"               description:"更新时间"`     // 更新时间
}
