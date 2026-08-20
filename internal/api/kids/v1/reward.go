package v1

import "github.com/gogf/gf/v2/frame/g"

type RewardPreset struct {
	Id          uint64 `json:"id" dc:"预设ID"`
	Title       string `json:"title" dc:"预设奖励标题"`
	Icon        string `json:"icon" dc:"图标标识"`
	ImageUrl    string `json:"imageUrl" dc:"奖励图片地址"`
	StarCost    int    `json:"starCost" dc:"默认所需星星数量"`
	Description string `json:"description" dc:"预设描述"`
}

type RewardItem struct {
	Id                 uint64   `json:"id" dc:"奖励ID"`
	CircleId           uint64   `json:"circleId" dc:"圈子ID"`
	Title              string   `json:"title" dc:"奖励标题"`
	Icon               string   `json:"icon" dc:"图标标识"`
	ImageUrl           string   `json:"imageUrl" dc:"奖励图片地址"`
	StarCost           int      `json:"starCost" dc:"所需星星数量"`
	Stock              int      `json:"stock" dc:"可用库存，-1表示不限量"`
	Description        string   `json:"description" dc:"奖励描述"`
	RepeatRule         string   `json:"repeatRule" dc:"重复兑换规则"`
	RepeatIntervalDays int      `json:"repeatIntervalDays" dc:"自定义重复兑换间隔天数"`
	KidIds             []uint64 `json:"kidIds" dc:"指派儿童ID列表，空表示全部儿童可见"`
	Progress           float64  `json:"progress" dc:"当前儿童兑换进度，0到1"`
	CanRedeem          bool     `json:"canRedeem" dc:"当前儿童是否可兑换"`
	NextRedeemAt       int64    `json:"nextRedeemAt" dc:"下一次可兑换时间戳"`
}

type RewardPresetListReq struct {
	g.Meta  `path:"/reward-presets" method:"get" tags:"儿童端奖励" summary:"查询奖励预设列表"`
	Keyword string `p:"keyword" dc:"搜索关键字"`
}

type RewardPresetListInput struct {
	Keyword string
}

type RewardPresetListOutput struct {
	List []RewardPreset
}

type RewardPresetListRes struct {
	List []RewardPreset `json:"list" dc:"奖励预设列表"`
}

type RewardListReq struct {
	g.Meta   `path:"/rewards" method:"get" tags:"儿童端奖励" summary:"查询奖励列表"`
	CircleId uint64 `p:"circleId" dc:"圈子ID"`
	KidId    uint64 `p:"kidId" dc:"儿童成员ID"`
}

type RewardListInput struct {
	CircleId uint64
	KidId    uint64
}

type RewardListOutput struct {
	List []RewardItem
}

type RewardListRes struct {
	List []RewardItem `json:"list" dc:"奖励列表"`
}

type RewardCreateReq struct {
	g.Meta             `path:"/rewards" method:"post" tags:"儿童端奖励" summary:"创建奖励"`
	CircleId           uint64   `json:"circleId" dc:"圈子ID"`
	PresetId           uint64   `json:"presetId" dc:"预设奖励ID"`
	Title              string   `json:"title" v:"required|max-length:128" dc:"奖励标题"`
	Icon               string   `json:"icon" dc:"图标标识"`
	ImageUrl           string   `json:"imageUrl" dc:"奖励图片地址"`
	StarCost           int      `json:"starCost" v:"required|min:1|max:99999" dc:"所需星星数量"`
	Stock              int      `json:"stock" dc:"可用库存，-1表示不限量"`
	Description        string   `json:"description" dc:"奖励描述"`
	RepeatRule         string   `json:"repeatRule" v:"in:none,daily,weekly,monthly,custom" dc:"重复兑换规则：none/daily/weekly/monthly/custom"`
	RepeatIntervalDays int      `json:"repeatIntervalDays" dc:"自定义重复兑换间隔天数"`
	KidIds             []uint64 `json:"kidIds" dc:"指派儿童ID列表，空表示全部儿童可见"`
}

type RewardCreateInput struct {
	UserId             uint64
	CircleId           uint64
	PresetId           uint64
	Title              string
	Icon               string
	ImageUrl           string
	StarCost           int
	Stock              int
	Description        string
	RepeatRule         string
	RepeatIntervalDays int
	KidIds             []uint64
}

type RewardCreateOutput struct {
	Reward RewardItem
}

type RewardCreateRes struct {
	Reward RewardItem `json:"reward" dc:"创建的奖励"`
}

type RewardDetailReq struct {
	g.Meta `path:"/rewards/{id}" method:"get" tags:"儿童端奖励" summary:"查询奖励详情"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"奖励ID"`
	KidId  uint64 `p:"kidId" dc:"儿童成员ID"`
}

type RewardDetailInput struct {
	Id    uint64
	KidId uint64
}

type RewardDetailOutput struct {
	Reward RewardItem
}

type RewardDetailRes struct {
	Reward RewardItem `json:"reward" dc:"奖励详情"`
}

type RewardUpdateReq struct {
	g.Meta             `path:"/rewards/{id}" method:"put" tags:"儿童端奖励" summary:"编辑奖励"`
	Id                 uint64   `p:"id" v:"required|min:1" dc:"奖励ID"`
	CircleId           uint64   `json:"circleId" dc:"圈子ID"`
	Title              string   `json:"title" v:"required|max-length:128" dc:"奖励标题"`
	Icon               string   `json:"icon" dc:"图标标识"`
	ImageUrl           string   `json:"imageUrl" dc:"奖励图片地址"`
	StarCost           int      `json:"starCost" v:"required|min:1|max:99999" dc:"所需星星数量"`
	Stock              int      `json:"stock" dc:"可用库存，-1表示不限量"`
	Description        string   `json:"description" dc:"奖励描述"`
	RepeatRule         string   `json:"repeatRule" v:"in:none,daily,weekly,monthly,custom" dc:"重复兑换规则：none/daily/weekly/monthly/custom"`
	RepeatIntervalDays int      `json:"repeatIntervalDays" dc:"自定义重复兑换间隔天数"`
	KidIds             []uint64 `json:"kidIds" dc:"指派儿童ID列表"`
}

type RewardUpdateInput struct {
	UserId             uint64
	Id                 uint64
	CircleId           uint64
	Title              string
	Icon               string
	ImageUrl           string
	StarCost           int
	Stock              int
	Description        string
	RepeatRule         string
	RepeatIntervalDays int
	KidIds             []uint64
}

type RewardUpdateOutput struct {
	Reward RewardItem
}

type RewardUpdateRes struct {
	Reward RewardItem `json:"reward" dc:"更新后的奖励"`
}

type RewardDeleteReq struct {
	g.Meta `path:"/rewards/{id}" method:"delete" tags:"儿童端奖励" summary:"删除奖励"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"奖励ID"`
}

type RewardDeleteInput struct {
	UserId uint64
	Id     uint64
}

type RewardDeleteOutput struct {
	Id uint64
}

type RewardDeleteRes struct {
	Id uint64 `json:"id" dc:"已删除奖励ID"`
}

type RewardRedeemReq struct {
	g.Meta `path:"/rewards/{id}/redeem" method:"post" tags:"儿童端奖励" summary:"兑换奖励"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"奖励ID"`
	KidId  uint64 `json:"kidId" v:"required|min:1" dc:"儿童成员ID"`
	Remark string `json:"remark" dc:"兑换备注"`
}

type RewardRedeemInput struct {
	Id     uint64
	KidId  uint64
	Remark string
}

type RewardRedeemOutput struct {
	Reward      RewardItem
	StarBalance int
}

type RewardRedeemRes struct {
	Reward      RewardItem `json:"reward" dc:"已兑换奖励"`
	StarBalance int        `json:"starBalance" dc:"兑换后的儿童星星余额"`
}

type RewardRedeemRecord struct {
	Id        uint64 `json:"id" dc:"兑换记录ID"`
	CircleId  uint64 `json:"circleId" dc:"圈子ID"`
	RewardId  uint64 `json:"rewardId" dc:"奖励ID"`
	KidId     uint64 `json:"kidId" dc:"儿童成员ID"`
	KidName   string `json:"kidName" dc:"儿童名称"`
	UserId    uint64 `json:"userId" dc:"兑换用户ID"`
	Title     string `json:"title" dc:"奖励标题"`
	Icon      string `json:"icon" dc:"图标标识"`
	ImageUrl  string `json:"imageUrl" dc:"奖励图片地址"`
	StarCost  int    `json:"starCost" dc:"消耗星星数量"`
	Remark    string `json:"remark" dc:"兑换备注"`
	CreatedAt int64  `json:"createdAt" dc:"创建时间戳"`
}

type RewardRecordListReq struct {
	g.Meta   `path:"/rewards/records" method:"get" tags:"儿童端奖励" summary:"查询奖励兑换历史"`
	CircleId uint64 `p:"circleId" dc:"圈子ID"`
	KidId    uint64 `p:"kidId" dc:"儿童成员ID"`
	Month    string `p:"month" dc:"月份筛选，格式YYYY-MM"`
	From     string `p:"from" dc:"开始日期，格式YYYY-MM-DD"`
	To       string `p:"to" dc:"结束日期，格式YYYY-MM-DD"`
}

type RewardRecordListInput struct {
	CircleId uint64
	KidId    uint64
	Month    string
	From     string
	To       string
}

type RewardRecordListOutput struct {
	Total int
	List  []RewardRedeemRecord
}

type RewardRecordListRes struct {
	Total int                  `json:"total" dc:"总兑换次数"`
	List  []RewardRedeemRecord `json:"list" dc:"奖励兑换历史"`
}
