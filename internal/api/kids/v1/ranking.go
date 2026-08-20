package v1

import "github.com/gogf/gf/v2/frame/g"

type StarRankingItem struct {
	Rank        int    `json:"rank" dc:"排名"`
	KidId       uint64 `json:"kidId" dc:"儿童成员ID"`
	Name        string `json:"name" dc:"显示名称"`
	Avatar      string `json:"avatar" dc:"头像"`
	AvatarStyle string `json:"avatarStyle" dc:"虚拟形象风格"`
	StarCount   int    `json:"starCount" dc:"星星余额"`
}

type StarRankingReq struct {
	g.Meta   `path:"/rankings/stars" method:"get" tags:"儿童端排行榜" summary:"查询群组星星排行榜"`
	CircleId uint64 `p:"circleId" v:"required|min:1" dc:"圈子ID"`
}

type StarRankingInput struct {
	UserId   uint64
	CircleId uint64
}

type StarRankingOutput struct {
	List []StarRankingItem
}

type StarRankingRes struct {
	List []StarRankingItem `json:"list" dc:"排行榜列表"`
}
