package kids

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// GetStatisticsV1 以合同任务完成和星星账本事实返回成员统计序列。
func (s *sKids) GetStatisticsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, func(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
		return s.v1StatisticsSeries(ctx, in, v1QueryFirst(in.Query, "member_id"))
	})
}

// CompareStatisticsV1 在相同快照、时区和 bucket 规则下比较两名成员。
func (s *sKids) CompareStatisticsV1(ctx context.Context, in v1.V1OperationInput) (*v1.V1OperationOutput, error) {
	return s.runV1(ctx, in, func(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
		base, cursor, err := s.v1StatisticsSeries(ctx, in, v1QueryFirst(in.Query, "base_member_id"))
		if err != nil {
			return nil, "", err
		}
		compare, _, err := s.v1StatisticsSeries(ctx, in, v1QueryFirst(in.Query, "compare_member_id"))
		if err != nil {
			return nil, "", err
		}
		return map[string]any{"base_member": base, "compare_member": compare, "snapshot_cursor": cursor}, cursor, nil
	})
}

// v1StatisticsSeries 根据请求的日、周或月 bucket 生成零值完整的合同统计序列。
func (s *sKids) v1StatisticsSeries(ctx context.Context, in v1.V1OperationInput, memberID string) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	member, err := utils.KidsDB(ctx).Model(consts.KidsV1MemberTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).Where("status", "active").One()
	if err != nil {
		return nil, "", err
	}
	if member.IsEmpty() {
		return nil, "", v1Error(404, "NOT_FOUND", false, "member is missing")
	}
	startAt, err := strconv.ParseInt(v1QueryFirst(in.Query, "start_at_ms"), 10, 64)
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "statistics start time is invalid")
	}
	endAt, err := strconv.ParseInt(v1QueryFirst(in.Query, "end_at_ms"), 10, 64)
	if err != nil || endAt <= startAt {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "statistics end time is invalid")
	}
	zoneID := v1QueryFirst(in.Query, "zone_id")
	location, err := time.LoadLocation(zoneID)
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "statistics zone is invalid")
	}
	start, end := time.UnixMilli(startAt), time.UnixMilli(endAt)
	metric, unit, weekStart := v1QueryFirst(in.Query, "metric"), v1QueryFirst(in.Query, "bucket_unit"), v1QueryFirst(in.Query, "week_start")
	values, err := s.v1StatisticsValues(ctx, circleID, memberID, metric, start, end, location, unit, weekStart)
	if err != nil {
		return nil, "", err
	}
	buckets := make([]map[string]any, 0)
	var total, peak, nonZero int64
	localEnd := end.In(location)
	for current, index := v1BucketStart(start.In(location), unit, weekStart), 0; current.Before(localEnd); current, index = v1NextBucketStart(current, unit), index+1 {
		next := v1NextBucketStart(current, unit)
		bucketStart, bucketEnd := current.UTC().UnixMilli(), next.UTC().UnixMilli()
		if bucketStart < startAt {
			bucketStart = startAt
		}
		if bucketEnd > endAt {
			bucketEnd = endAt
		}
		value := values[current.Format(time.RFC3339Nano)]
		total += value
		if index == 0 || value > peak {
			peak = value
		}
		if value != 0 {
			nonZero++
		}
		buckets = append(buckets, map[string]any{"index": index, "start_at_ms": bucketStart, "end_at_ms": bucketEnd, "local_start_date": current.Format(consts.DateLayout), "local_end_date_exclusive": next.Format(consts.DateLayout), "value": value})
	}
	cursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"member_id": memberID, "metric": metric, "period_type": v1QueryFirst(in.Query, "period_type"), "bucket_unit": unit, "zone_id": zoneID, "start_at_ms": startAt, "end_at_ms": endAt, "buckets": buckets, "summary": map[string]any{"total": total, "peak_value": peak, "non_zero_bucket_count": nonZero}, "as_of_cursor": cursor}, cursor, nil
}

// v1StatisticsValues 把合同任务完成或星星流水按请求时区归入对应统计 bucket。
func (s *sKids) v1StatisticsValues(ctx context.Context, circleID, memberID, metric string, start, end time.Time, location *time.Location, unit, weekStart string) (map[string]int64, error) {
	values := make(map[string]int64)
	if metric == "tasks" {
		rows, err := utils.KidsDB(ctx).Model(consts.KidsV1TaskCompletionTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).Where("completed_at >= ?", start).Where("completed_at < ?", end).All()
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			key := v1BucketStart(row["completed_at"].Time().In(location), unit, weekStart).Format(time.RFC3339Nano)
			values[key]++
		}
		return values, nil
	}
	rows, err := utils.KidsDB(ctx).Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).Where("created_at >= ?", start).Where("created_at < ?", end).All()
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		key := v1BucketStart(row["created_at"].Time().In(location), unit, weekStart).Format(time.RFC3339Nano)
		values[key] += row["delta"].Int64()
	}
	return values, nil
}

// v1BucketStart 返回本地时间点所属统计 bucket 的起点。
func v1BucketStart(value time.Time, unit, weekStart string) time.Time {
	day := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
	switch unit {
	case "month":
		return time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, day.Location())
	case "week":
		first := time.Monday
		if weekStart == "sunday" {
			first = time.Sunday
		}
		return day.AddDate(0, 0, -(int(day.Weekday())-int(first)+7)%7)
	default:
		return day
	}
}

// v1NextBucketStart 计算本地日、周或月 bucket 的下一个起点。
func v1NextBucketStart(value time.Time, unit string) time.Time {
	switch unit {
	case "month":
		return value.AddDate(0, 1, 0)
	case "week":
		return value.AddDate(0, 0, 7)
	default:
		return value.AddDate(0, 0, 1)
	}
}

// GetAnalyticsSummary 查询单个儿童在指定时间范围内的任务或星星统计。
func (s *sKids) GetAnalyticsSummary(ctx context.Context, in v1.AnalyticsSummaryInput) (*v1.AnalyticsSummaryOutput, error) {
	if err := validateAnalyticsPayload(in.Metric, in.Range); err != nil {
		return nil, err
	}
	from, to, rangeType, err := analyticsDateRange(in.Range, in.From, in.To, in.BaseDate)
	if err != nil {
		return nil, err
	}
	daily, hourly, total, err := collectAnalyticsPoints(ctx, in.KidId, in.Metric, from, to)
	if err != nil {
		return nil, err
	}
	return &v1.AnalyticsSummaryOutput{KidId: in.KidId, Metric: normalizedAnalyticsMetric(in.Metric), Range: rangeType, From: from.Format(consts.DateLayout), To: to.Format(consts.DateLayout), Total: total, Daily: daily, Hourly: hourly}, nil
}

// ListCompletedTaskDetails 查询指定儿童在时间范围内的已完成任务明细。
func (s *sKids) ListCompletedTaskDetails(ctx context.Context, in v1.CompletedTaskListInput) (*v1.CompletedTaskListOutput, error) {
	if !validAnalyticsRange(in.Range) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "unsupported analytics range")
	}
	from, to, _, err := analyticsDateRange(in.Range, in.From, in.To, in.BaseDate)
	if err != nil {
		return nil, err
	}
	records, err := completedTaskModel(ctx, in.KidId, from, to).OrderDesc("ta.completed_at").All()
	if err != nil {
		return nil, err
	}
	out := &v1.CompletedTaskListOutput{From: from.Format(consts.DateLayout), To: to.Format(consts.DateLayout), Total: len(records)}
	for _, record := range records {
		out.List = append(out.List, completedTaskDetailFromDB(record))
	}
	return out, nil
}

// CompareAnalytics 查询两个儿童在同一时间范围内的任务或星星统计对比。
func (s *sKids) CompareAnalytics(ctx context.Context, in v1.AnalyticsCompareInput) (*v1.AnalyticsCompareOutput, error) {
	if err := validateAnalyticsPayload(in.Metric, in.Range); err != nil {
		return nil, err
	}
	from, to, rangeType, err := analyticsDateRange(in.Range, in.From, in.To, in.BaseDate)
	if err != nil {
		return nil, err
	}
	left, err := analyticsCompareMember(ctx, in.KidId, in.Metric, from, to)
	if err != nil {
		return nil, err
	}
	right, err := analyticsCompareMember(ctx, in.CompareKidId, in.Metric, from, to)
	if err != nil {
		return nil, err
	}
	return &v1.AnalyticsCompareOutput{Metric: normalizedAnalyticsMetric(in.Metric), Range: rangeType, From: from.Format(consts.DateLayout), To: to.Format(consts.DateLayout), Left: left, Right: right}, nil
}

// analyticsCompareMember 汇总单个成员在对比接口中的展示资料和统计数据。
func analyticsCompareMember(ctx context.Context, kidId uint64, metric string, from, to time.Time) (v1.AnalyticsCompareMember, error) {
	member, err := findFamilyMember(ctx, nil, kidId)
	if err != nil {
		return v1.AnalyticsCompareMember{}, err
	}
	if member == nil {
		return v1.AnalyticsCompareMember{}, utils.NewCodef(gcode.CodeInvalidParameter, "error.kid_not_found", kidId)
	}
	daily, hourly, total, err := collectAnalyticsPoints(ctx, kidId, metric, from, to)
	if err != nil {
		return v1.AnalyticsCompareMember{}, err
	}
	return v1.AnalyticsCompareMember{KidId: kidId, Name: member.Name, Avatar: member.Avatar, Total: total, Daily: daily, Hourly: hourly}, nil
}

// collectAnalyticsPoints 按天和小时聚合任务或星星指标。
func collectAnalyticsPoints(ctx context.Context, kidId uint64, metric string, from, to time.Time) ([]v1.AnalyticsPoint, []v1.AnalyticsPoint, int, error) {
	metric = normalizedAnalyticsMetric(metric)
	daily := initDailyPoints(from, to)
	hourly := initHourlyPoints()
	total := 0
	if metric == v1.AnalyticsMetricStars {
		records, err := starAnalyticsRecords(ctx, kidId, from, to)
		if err != nil {
			return nil, nil, 0, err
		}
		for _, record := range records {
			value := record["change_amount"].Int()
			createdAt := utils.ParseDBTime(record["created_at"].Val())
			addAnalyticsValue(daily, hourly, createdAt, value)
			total += value
		}
		return daily, hourly, total, nil
	}
	records, err := completedTaskModel(ctx, kidId, from, to).All()
	if err != nil {
		return nil, nil, 0, err
	}
	for _, record := range records {
		completedAt := utils.ParseDBTime(record["completed_at"].Val())
		addAnalyticsValue(daily, hourly, completedAt, 1)
		total++
	}
	return daily, hourly, total, nil
}

// completedTaskModel 构造指定儿童和时间范围的完成任务查询。
func completedTaskModel(ctx context.Context, kidId uint64, from, to time.Time) *gdb.Model {
	return utils.KidsDB(ctx).Model(consts.KidsTaskAssigneeTable+" ta").Ctx(ctx).
		LeftJoin(consts.KidsTaskTable+" t", "t.id = ta.task_id").
		Fields("ta.task_id,ta.kid_id,ta.photo_url,ta.completed_at,t.title,t.icon,t.star").
		Where("ta.kid_id", kidId).
		Where("ta.completed", 1).
		WhereGTE("ta.completed_at", from.Format(consts.MySQLTimeLayout)).
		WhereLTE("ta.completed_at", endOfDay(to).Format(consts.MySQLTimeLayout)).
		Where("t.deleted_at IS NULL")
}

// starAnalyticsRecords 查询指定儿童和时间范围内的星星流水。
func starAnalyticsRecords(ctx context.Context, kidId uint64, from, to time.Time) (gdb.Result, error) {
	return utils.KidsDB(ctx).Model(consts.KidsStarRecordTable).Ctx(ctx).
		Fields("change_amount,created_at").
		Where("kid_id", kidId).
		WhereGTE("created_at", from.Format(consts.MySQLTimeLayout)).
		WhereLTE("created_at", endOfDay(to).Format(consts.MySQLTimeLayout)).
		All()
}

// completedTaskDetailFromDB 将数据库完成任务记录转换为接口明细结构。
func completedTaskDetailFromDB(record gdb.Record) v1.CompletedTaskDetail {
	return v1.CompletedTaskDetail{TaskId: record["task_id"].Uint64(), KidId: record["kid_id"].Uint64(), Title: record["title"].String(), Icon: record["icon"].String(), Star: record["star"].Int(), PhotoUrl: record["photo_url"].String(), CompletedAt: utils.ParseDBTime(record["completed_at"].Val())}
}

// analyticsDateRange 按周、月或自定义日期解析统计时间范围。
func analyticsDateRange(rangeType, fromText, toText, baseDateText string) (time.Time, time.Time, string, error) {
	rangeType = normalizedAnalyticsRange(rangeType)
	baseDate, err := time.ParseInLocation(consts.DateLayout, utils.NormalizeDate(baseDateText), time.Local)
	if err != nil {
		baseDate = time.Now()
	}
	switch rangeType {
	case v1.AnalyticsRangeCustom:
		from, err := time.ParseInLocation(consts.DateLayout, strings.TrimSpace(fromText), time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, rangeType, gerror.NewCode(gcode.CodeInvalidParameter, "invalid analytics date range")
		}
		to, err := time.ParseInLocation(consts.DateLayout, strings.TrimSpace(toText), time.Local)
		if err != nil || to.Before(from) {
			return time.Time{}, time.Time{}, rangeType, gerror.NewCode(gcode.CodeInvalidParameter, "invalid analytics date range")
		}
		return startOfDay(from), startOfDay(to), rangeType, nil
	case v1.AnalyticsRangeMonthly:
		from := time.Date(baseDate.Year(), baseDate.Month(), 1, 0, 0, 0, 0, time.Local)
		to := from.AddDate(0, 1, -1)
		return from, to, rangeType, nil
	default:
		weekday := int(baseDate.Weekday())
		from := startOfDay(baseDate.AddDate(0, 0, -weekday))
		to := from.AddDate(0, 0, 6)
		return from, to, v1.AnalyticsRangeWeekly, nil
	}
}

// initDailyPoints 初始化日期范围内每天的统计点。
func initDailyPoints(from, to time.Time) []v1.AnalyticsPoint {
	points := make([]v1.AnalyticsPoint, 0)
	for current := startOfDay(from); !current.After(to); current = current.AddDate(0, 0, 1) {
		points = append(points, v1.AnalyticsPoint{Label: current.Format("01/02"), Date: current.Format(consts.DateLayout)})
	}
	return points
}

// initHourlyPoints 初始化一天 24 小时的统计点。
func initHourlyPoints() []v1.AnalyticsPoint {
	points := make([]v1.AnalyticsPoint, 24)
	for hour := 0; hour < 24; hour++ {
		points[hour] = v1.AnalyticsPoint{Label: time.Date(2000, 1, 1, hour, 0, 0, 0, time.Local).Format("15:00"), Hour: hour}
	}
	return points
}

// addAnalyticsValue 将一条记录累加到按天和按小时的统计点中。
func addAnalyticsValue(daily []v1.AnalyticsPoint, hourly []v1.AnalyticsPoint, ts int64, value int) {
	if ts == 0 {
		return
	}
	t := time.Unix(ts, 0)
	date := t.Format(consts.DateLayout)
	for i := range daily {
		if daily[i].Date == date {
			daily[i].Value += value
			break
		}
	}
	hour := t.Hour()
	if hour >= 0 && hour < len(hourly) {
		hourly[hour].Value += value
	}
}

// normalizedAnalyticsMetric 标准化统计指标。
func normalizedAnalyticsMetric(metric string) string {
	if strings.TrimSpace(metric) == v1.AnalyticsMetricStars {
		return v1.AnalyticsMetricStars
	}
	return v1.AnalyticsMetricTasks
}

// validateAnalyticsPayload 校验统计指标和时间范围枚举。
func validateAnalyticsPayload(metric string, rangeType string) error {
	if !validAnalyticsMetric(metric) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported analytics metric")
	}
	if !validAnalyticsRange(rangeType) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported analytics range")
	}
	return nil
}

// validAnalyticsMetric 校验统计指标枚举。
func validAnalyticsMetric(metric string) bool {
	switch strings.TrimSpace(metric) {
	case v1.AnalyticsMetricTasks, v1.AnalyticsMetricStars:
		return true
	default:
		return false
	}
}

// validAnalyticsRange 校验统计时间范围，空值允许走默认周视图。
func validAnalyticsRange(rangeType string) bool {
	switch strings.TrimSpace(rangeType) {
	case "", v1.AnalyticsRangeWeekly, v1.AnalyticsRangeMonthly, v1.AnalyticsRangeCustom:
		return true
	default:
		return false
	}
}

// normalizedAnalyticsRange 标准化统计时间范围。
func normalizedAnalyticsRange(rangeType string) string {
	switch strings.TrimSpace(rangeType) {
	case v1.AnalyticsRangeWeekly, v1.AnalyticsRangeMonthly, v1.AnalyticsRangeCustom:
		return strings.TrimSpace(rangeType)
	default:
		return v1.AnalyticsRangeWeekly
	}
}

// startOfDay 返回日期当天零点。
func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.Local)
}

// endOfDay 返回日期当天最后一秒。
func endOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, 0, time.Local)
}
