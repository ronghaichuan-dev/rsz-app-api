package v1

import "github.com/gogf/gf/v2/frame/g"

type NotificationItem struct {
	Id        uint64 `json:"id" dc:"通知ID"`
	MemberId  uint64 `json:"memberId" dc:"目标成员ID，0表示全家庭"`
	Type      string `json:"type" dc:"通知类型"`
	Title     string `json:"title" dc:"通知标题"`
	Content   string `json:"content" dc:"通知内容"`
	Read      bool   `json:"read" dc:"是否已读"`
	CreatedAt int64  `json:"createdAt" dc:"创建时间戳"`
}

type NotificationListReq struct {
	g.Meta   `path:"/notifications" method:"get" tags:"儿童端通知" summary:"查询通知列表"`
	MemberId uint64 `p:"memberId" dc:"目标成员ID"`
	Unread   bool   `p:"unread" dc:"仅查询未读通知"`
}

type NotificationListInput struct {
	MemberId uint64
	Unread   bool
}

type NotificationListOutput struct {
	List []NotificationItem
}

type NotificationListRes struct {
	List []NotificationItem `json:"list" dc:"通知列表"`
}

type NotificationReadReq struct {
	g.Meta `path:"/notifications/{id}/read" method:"post" tags:"儿童端通知" summary:"标记通知已读"`
	Id     uint64 `p:"id" v:"required|min:1" dc:"通知ID"`
}

type NotificationReadInput struct {
	Id uint64
}

type NotificationReadOutput struct {
	Notification NotificationItem
}

type NotificationReadRes struct {
	Notification NotificationItem `json:"notification" dc:"已读通知"`
}
