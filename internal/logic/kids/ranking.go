package kids

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// GetStarRanking 查询群组内儿童星星排行榜。
func (s *sKids) GetStarRanking(ctx context.Context, in v1.StarRankingInput) (*v1.StarRankingOutput, error) {
	if ok, err := userCanAccessCircle(ctx, nil, in.UserId, in.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "circle permission denied")
	}
	records, err := utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx).
		Where("circle_id", in.CircleId).
		Where("deleted_at IS NULL").
		OrderAsc("id").All()
	if err != nil {
		return nil, err
	}
	items := make([]v1.StarRankingItem, 0, len(records))
	for _, record := range records {
		kidId := record["id"].Uint64()
		items = append(items, v1.StarRankingItem{KidId: kidId, Name: record["name"].String(), Avatar: record["avatar"].String(), AvatarStyle: record["avatar_style"].String(), StarCount: getStarBalanceValue(ctx, kidId)})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].StarCount == items[j].StarCount {
			return items[i].KidId < items[j].KidId
		}
		return items[i].StarCount > items[j].StarCount
	})
	for i := range items {
		items[i].Rank = i + 1
	}
	return &v1.StarRankingOutput{List: items}, nil
}
