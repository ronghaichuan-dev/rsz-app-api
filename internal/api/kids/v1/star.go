package v1

import "github.com/gogf/gf/v2/frame/g"

const (
	StarRecordTypeTask       = "task"
	StarRecordTypeReward     = "reward"
	StarRecordTypeAdjustment = "adjustment"
)

type StarRecord struct {
	Id        uint64 `json:"id" dc:"流水ID"`
	KidId     uint64 `json:"kidId" dc:"儿童成员ID"`
	Change    int    `json:"change" dc:"星星变动数量"`
	Balance   int    `json:"balance" dc:"变动后余额"`
	Type      string `json:"type" dc:"流水类型"`
	Title     string `json:"title" dc:"流水标题"`
	Remark    string `json:"remark" dc:"流水备注"`
	CreatedAt int64  `json:"createdAt" dc:"创建时间戳"`
}

type StarBalanceReq struct {
	g.Meta `path:"/stars/balance" method:"get" tags:"儿童端星星" summary:"查询儿童星星余额"`
	KidId  uint64 `p:"kidId" v:"required|min:1" dc:"儿童成员ID"`
}

type StarBalanceInput struct {
	KidId uint64
}

type StarBalanceOutput struct {
	KidId   uint64
	Balance int
}

type StarBalanceRes struct {
	KidId   uint64 `json:"kidId" dc:"儿童成员ID"`
	Balance int    `json:"balance" dc:"星星余额"`
}

type StarAdjustReq struct {
	g.Meta `path:"/stars/adjust" method:"post" tags:"儿童端星星" summary:"调整儿童星星余额"`
	KidId  uint64 `json:"kidId" v:"required|min:1" dc:"儿童成员ID"`
	Amount int    `json:"amount" v:"required" dc:"调整数量，可正可负"`
	Reason string `json:"reason" v:"required|max-length:255" dc:"调整原因"`
}

type StarAdjustInput struct {
	KidId  uint64
	Amount int
	Reason string
}

type StarAdjustOutput struct {
	KidId   uint64
	Balance int
	Record  StarRecord
}

type StarAdjustRes struct {
	KidId   uint64     `json:"kidId" dc:"儿童成员ID"`
	Balance int        `json:"balance" dc:"调整后的星星余额"`
	Record  StarRecord `json:"record" dc:"创建的星星流水"`
}

type StarRecordListReq struct {
	g.Meta `path:"/stars/records" method:"get" tags:"儿童端星星" summary:"查询星星流水"`
	KidId  uint64 `p:"kidId" v:"required|min:1" dc:"儿童成员ID"`
	Type   string `p:"type" v:"in:task,reward,adjustment" dc:"流水类型筛选：task任务、reward奖励、adjustment调整"`
	From   string `p:"from" dc:"开始日期，格式YYYY-MM-DD"`
	To     string `p:"to" dc:"结束日期，格式YYYY-MM-DD"`
}

type StarRecordListInput struct {
	KidId uint64
	Type  string
	From  string
	To    string
}

type StarRecordListOutput struct {
	List []StarRecord
}

type StarRecordListRes struct {
	List []StarRecord `json:"list" dc:"星星流水列表"`
}
