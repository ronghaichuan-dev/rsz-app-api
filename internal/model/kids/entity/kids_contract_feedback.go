// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractFeedback is the golang structure for table kids_contract_feedback.
type KidsContractFeedback struct {
	Id                    uint64      `json:"id"                    orm:"id"                      description:"主键"`     // 主键
	FeedbackId            string      `json:"feedbackId"            orm:"feedback_id"             description:"反馈标识"`   // 反馈标识
	AccountId             string      `json:"accountId"             orm:"account_id"              description:"账号标识"`   // 账号标识
	Category              string      `json:"category"              orm:"category"                description:"反馈分类"`   // 反馈分类
	Content               string      `json:"content"               orm:"content"                 description:"反馈内容"`   // 反馈内容
	ContactType           string      `json:"contactType"           orm:"contact_type"            description:"联系类型"`   // 联系类型
	Contact               string      `json:"contact"               orm:"contact"                 description:"联系方式"`   // 联系方式
	PrivacyConsentVersion string      `json:"privacyConsentVersion" orm:"privacy_consent_version" description:"隐私同意版本"` // 隐私同意版本
	AttachmentAssetIds    string      `json:"attachmentAssetIds"    orm:"attachment_asset_ids"    description:"附件资产标识"` // 附件资产标识
	Version               uint64      `json:"version"               orm:"version"                 description:"版本号"`    // 版本号
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"              description:"创建时间"`   // 创建时间
}
