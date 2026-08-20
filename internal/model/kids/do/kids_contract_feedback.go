// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractFeedback is the golang structure of table kids_contract_feedback for DAO operations like Where/Data.
type KidsContractFeedback struct {
	g.Meta                `orm:"table:kids_contract_feedback, do:true"`
	Id                    any         // 主键
	FeedbackId            any         // 反馈标识
	AccountId             any         // 账号标识
	Category              any         // 反馈分类
	Content               any         // 反馈内容
	ContactType           any         // 联系类型
	Contact               any         // 联系方式
	PrivacyConsentVersion any         // 隐私同意版本
	AttachmentAssetIds    any         // 附件资产标识
	Version               any         // 版本号
	CreatedAt             *gtime.Time // 创建时间
}
