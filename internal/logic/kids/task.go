package kids

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/google/uuid"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// ListTaskPresets 从数据库读取启用的任务预设，并支持按关键字搜索。
func (s *sKids) ListTaskPresets(ctx context.Context, in v1.TaskPresetListInput) (*v1.TaskPresetListOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsTaskPresetTable).Ctx(ctx).Where("enabled", 1)
	if strings.TrimSpace(in.Keyword) != "" {
		model = model.WhereLike("title", "%"+strings.TrimSpace(in.Keyword)+"%")
	}
	records, err := model.OrderAsc("id").All()
	if err != nil {
		return nil, err
	}
	out := &v1.TaskPresetListOutput{}
	for _, record := range records {
		out.List = append(out.List, taskPresetFromDB(record))
	}
	return out, nil
}

// GetTask 查询未删除任务详情，并返回分配状态和标签信息。
func (s *sKids) GetTask(ctx context.Context, in v1.TaskDetailInput) (*v1.TaskDetailOutput, error) {
	record, err := taskRecordById(ctx, nil, in.Id)
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "task not found")
	}
	item, err := taskItemFromDB(ctx, nil, record)
	if err != nil {
		return nil, err
	}
	return &v1.TaskDetailOutput{Task: item}, nil
}

// taskPresetFromDB 将数据库任务预设记录转换为接口响应结构。
func taskPresetFromDB(record gdb.Record) v1.TaskPreset {
	return v1.TaskPreset{
		Id:          record["id"].Uint64(),
		Title:       record["title"].String(),
		Icon:        record["icon"].String(),
		Star:        record["star"].Int(),
		NeedPhoto:   record["need_photo"].Int() == 1,
		Description: record["description"].String(),
	}
}

// ListTasks 按日期、孩子、标签和状态从数据库查询任务列表。
func (s *sKids) ListTasks(ctx context.Context, in v1.TaskListInput) (*v1.TaskListOutput, error) {
	date := utils.NormalizeDate(in.Date)
	model := taskListModel(ctx, nil).
		Where("t.task_date", date).
		Where("t.deleted_at IS NULL").
		OrderAsc("t.id")
	if in.TagId > 0 {
		model = model.Where("t.tag_id", in.TagId)
	}
	if in.KidId > 0 {
		model = model.LeftJoin(consts.KidsTaskAssigneeTable+" ta", "ta.task_id = t.id").Where("ta.kid_id", in.KidId)
	}
	records, err := model.All()
	if err != nil {
		return nil, err
	}
	out := &v1.TaskListOutput{Date: date}
	for _, record := range records {
		item, err := taskItemFromDB(ctx, nil, record)
		if err != nil {
			return nil, err
		}
		if in.KidId > 0 && !taskVisibleForKid(item, in.KidId) {
			continue
		}
		if !matchTaskStatus(item, strings.TrimSpace(in.Status)) {
			continue
		}
		out.List = append(out.List, item)
	}
	return out, nil
}

// CreateTask 校验任务入参，并在事务中按重复规则持久化任务和任务分配关系。
func (s *sKids) CreateTask(ctx context.Context, in v1.TaskCreateInput) (*v1.TaskCreateOutput, error) {
	if err := validateTaskCreatePayload(ctx, nil, in); err != nil {
		return nil, err
	}
	dates, err := expandTaskDates(in)
	if err != nil {
		return nil, err
	}

	var task v1.TaskItem
	err = utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		firstId := int64(0)
		for _, date := range dates {
			id, err := tx.Model(consts.KidsTaskTable).Ctx(ctx).Data(taskDataFromCreateInput(in, date)).InsertAndGetId()
			if err != nil {
				return err
			}
			if firstId == 0 {
				firstId = id
			}
			if err = replaceTaskAssignees(ctx, tx, uint64(id), in.KidIds); err != nil {
				return err
			}
			for _, kidId := range in.KidIds {
				if _, err = createNotification(ctx, tx, kidId, "task", commoni18n.T(ctx, "You have a new task!"), strings.TrimSpace(in.Title)); err != nil {
					return err
				}
			}
		}
		record, err := taskRecordById(ctx, tx, uint64(firstId))
		if err != nil {
			return err
		}
		task, err = taskItemFromDB(ctx, tx, record)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCreateOutput{Task: task}, nil
}

// UpdateTask 校验任务入参并更新未删除任务，同时重建任务分配关系。
func (s *sKids) UpdateTask(ctx context.Context, in v1.TaskUpdateInput) (*v1.TaskUpdateOutput, error) {
	if err := validateTaskUpdatePayload(ctx, nil, in); err != nil {
		return nil, err
	}
	var task v1.TaskItem
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		record, err := taskRecordById(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		if record.IsEmpty() {
			return gerror.NewCode(gcode.CodeNotFound, "task not found")
		}
		if record["completed"].Int() == 1 {
			return gerror.NewCode(gcode.CodeInvalidOperation, "completed task cannot be edited")
		}
		data := taskDataFromUpdateInput(in)
		if _, err = tx.Model(consts.KidsTaskTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").Data(data).Update(); err != nil {
			return err
		}
		if err = replaceTaskAssignees(ctx, tx, in.Id, in.KidIds); err != nil {
			return err
		}
		record, err = taskRecordById(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		task, err = taskItemFromDB(ctx, tx, record)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskUpdateOutput{Task: task}, nil
}

// DeleteTask 软删除任务，保留完成历史和星星流水审计记录。
func (s *sKids) DeleteTask(ctx context.Context, in v1.TaskDeleteInput) (*v1.TaskDeleteOutput, error) {
	result, err := utils.KidsDB(ctx).Model(consts.KidsTaskTable).Ctx(ctx).
		Where("id", in.Id).
		Where("deleted_at IS NULL").
		Data(map[string]any{"deleted_at": time.Now().Format(consts.MySQLTimeLayout)}).
		Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "task not found")
	}
	return &v1.TaskDeleteOutput{Id: in.Id}, nil
}

// CompleteTask 在事务中按完成模式完成任务、保存照片凭证，并写入星星流水。
func (s *sKids) CompleteTask(ctx context.Context, in v1.TaskCompleteInput) (*v1.TaskCompleteOutput, error) {
	var task v1.TaskItem
	var balance int
	var taskCompleted bool
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		record, err := taskRecordById(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		if record.IsEmpty() {
			return gerror.NewCode(gcode.CodeNotFound, "task not found")
		}
		currentTask, err := taskItemFromDB(ctx, tx, record)
		if err != nil {
			return err
		}
		if !currentTask.CanComplete {
			return gerror.NewCode(gcode.CodeInvalidOperation, "future task cannot be completed")
		}
		if currentTask.Completed {
			return gerror.NewCode(gcode.CodeInvalidOperation, "task already completed")
		}
		assignee, ok := findTaskAssigneeStatus(currentTask, in.KidId)
		if !ok {
			return gerror.NewCode(gcode.CodeInvalidParameter, "kid is not assigned to this task")
		}
		if !assignee.Active {
			return gerror.NewCode(gcode.CodeInvalidParameter, "kid is not the active assignee for this task")
		}
		if assignee.Completed {
			return gerror.NewCode(gcode.CodeInvalidOperation, "kid already completed this task")
		}
		if currentTask.NeedPhotoProof && strings.TrimSpace(in.PhotoUrl) == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, "photoUrl is required for this task")
		}
		completedAt := time.Now().Format(consts.MySQLTimeLayout)
		if _, err = tx.Model(consts.KidsTaskAssigneeTable).Ctx(ctx).
			Where("task_id", in.Id).
			Where("kid_id", in.KidId).
			Data(map[string]any{"completed": 1, "photo_url": strings.TrimSpace(in.PhotoUrl), "completed_at": completedAt}).
			Update(); err != nil {
			return err
		}
		balance, err = addStarRecord(ctx, tx, in.KidId, currentTask.Star, v1.StarRecordTypeTask, commoni18n.Tf(ctx, "star.task_completed", currentTask.Title), "")
		if err != nil {
			return err
		}
		taskCompleted, err = shouldMarkTaskCompleted(ctx, tx, in.Id, currentTask.CompletionMode)
		if err != nil {
			return err
		}
		if taskCompleted {
			if _, err = tx.Model(consts.KidsTaskTable).Ctx(ctx).Where("id", in.Id).Data(map[string]any{
				"completed":    1,
				"completed_by": in.KidId,
				"photo_url":    strings.TrimSpace(in.PhotoUrl),
				"completed_at": completedAt,
			}).Update(); err != nil {
				return err
			}
		}
		if _, err = createNotification(ctx, tx, in.KidId, "progress", commoni18n.T(ctx, "Task completed"), commoni18n.Tf(ctx, "notice.earned_stars", currentTask.Title)); err != nil {
			return err
		}
		record, err = taskRecordById(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		task, err = taskItemFromDB(ctx, tx, record)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCompleteOutput{Task: task, StarBalance: balance, TaskCompleted: taskCompleted}, nil
}

// CancelTask 取消指定儿童的任务完成状态，并写入负向星星流水回滚奖励。
func (s *sKids) CancelTask(ctx context.Context, in v1.TaskCancelInput) (*v1.TaskCancelOutput, error) {
	var task v1.TaskItem
	var balance int
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		record, err := taskRecordById(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		if record.IsEmpty() {
			return gerror.NewCode(gcode.CodeNotFound, "task not found")
		}
		currentTask, err := taskItemFromDB(ctx, tx, record)
		if err != nil {
			return err
		}
		assignee, ok := findTaskAssigneeStatus(currentTask, in.KidId)
		if !ok {
			return gerror.NewCode(gcode.CodeInvalidParameter, "kid is not assigned to this task")
		}
		if !assignee.Completed {
			return gerror.NewCode(gcode.CodeInvalidOperation, "task is not completed")
		}
		if _, err = tx.Model(consts.KidsTaskAssigneeTable).Ctx(ctx).
			Where("task_id", in.Id).
			Where("kid_id", in.KidId).
			Data(map[string]any{"completed": 0, "photo_url": "", "completed_at": nil}).
			Update(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsTaskTable).Ctx(ctx).Where("id", in.Id).Data(map[string]any{
			"completed":    0,
			"completed_by": 0,
			"photo_url":    "",
			"completed_at": nil,
		}).Update(); err != nil {
			return err
		}
		reason := strings.TrimSpace(in.Reason)
		if reason == "" {
			reason = commoni18n.T(ctx, "Task completion canceled")
		}
		balance, err = addStarRecord(ctx, tx, in.KidId, -currentTask.Star, v1.StarRecordTypeAdjustment, commoni18n.Tf(ctx, "star.task_canceled", currentTask.Title), reason)
		if err != nil {
			return err
		}
		record, err = taskRecordById(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		task, err = taskItemFromDB(ctx, tx, record)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCancelOutput{Task: task, StarBalance: balance}, nil
}

// taskRecordById 查询任务及其标签信息，传入事务时使用事务连接。
func taskRecordById(ctx context.Context, tx gdb.TX, id uint64) (gdb.Record, error) {
	return taskListModel(ctx, tx).Where("t.id", id).Where("t.deleted_at IS NULL").One()
}

// taskListModel 创建带标签字段的任务查询模型。
func taskListModel(ctx context.Context, tx gdb.TX) *gdb.Model {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsTaskTable + " t").Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsTaskTable + " t").Ctx(ctx)
	}
	return model.Fields("t.*, tg.name AS tag_name, tg.color AS tag_color").LeftJoin(consts.KidsTaskTagTable+" tg", "tg.id = t.tag_id AND tg.deleted_at IS NULL")
}

// taskDataFromCreateInput 将创建入参转换为任务表写入字段。
func taskDataFromCreateInput(in v1.TaskCreateInput, date string) map[string]any {
	return map[string]any{
		"title":                   strings.TrimSpace(in.Title),
		"icon":                    strings.TrimSpace(in.Icon),
		"note":                    strings.TrimSpace(in.Note),
		"star":                    in.Star,
		"task_date":               date,
		"completion_mode":         strings.TrimSpace(in.CompletionMode),
		"repeat_rule":             normalizedRepeatRule(in.RepeatRule),
		"repeat_end_type":         normalizedRepeatEndType(in.RepeatEndType),
		"repeat_end_date":         nullDate(in.RepeatEndDate),
		"repeat_end_count":        in.RepeatEndCount,
		"time_limit_type":         normalizedTimeLimitType(in.TimeLimitType),
		"time_limit_start":        strings.TrimSpace(in.TimeLimitStart),
		"time_limit_end":          strings.TrimSpace(in.TimeLimitEnd),
		"reminder_type":           normalizedReminderType(in.ReminderType),
		"reminder_at":             strings.TrimSpace(in.ReminderAt),
		"reminder_offset_minutes": in.ReminderOffsetMinutes,
		"need_photo_proof":        utils.BoolToInt(in.NeedPhotoProof),
		"tag_id":                  in.TagId,
	}
}

// taskDataFromUpdateInput 将编辑入参转换为任务表更新字段。
func taskDataFromUpdateInput(in v1.TaskUpdateInput) map[string]any {
	return map[string]any{
		"title":                   strings.TrimSpace(in.Title),
		"icon":                    strings.TrimSpace(in.Icon),
		"note":                    strings.TrimSpace(in.Note),
		"star":                    in.Star,
		"completion_mode":         strings.TrimSpace(in.CompletionMode),
		"repeat_rule":             normalizedRepeatRule(in.RepeatRule),
		"repeat_end_type":         normalizedRepeatEndType(in.RepeatEndType),
		"repeat_end_date":         nullDate(in.RepeatEndDate),
		"repeat_end_count":        in.RepeatEndCount,
		"time_limit_type":         normalizedTimeLimitType(in.TimeLimitType),
		"time_limit_start":        strings.TrimSpace(in.TimeLimitStart),
		"time_limit_end":          strings.TrimSpace(in.TimeLimitEnd),
		"reminder_type":           normalizedReminderType(in.ReminderType),
		"reminder_at":             strings.TrimSpace(in.ReminderAt),
		"reminder_offset_minutes": in.ReminderOffsetMinutes,
		"need_photo_proof":        utils.BoolToInt(in.NeedPhotoProof),
		"tag_id":                  in.TagId,
	}
}

// validateTaskCreatePayload 校验创建任务全部前置参数，避免非法枚举被默认值吞掉。
func validateTaskCreatePayload(ctx context.Context, tx gdb.TX, in v1.TaskCreateInput) error {
	if err := validateTaskPayload(ctx, tx, in.Title, in.Star, in.KidIds, in.CompletionMode, in.RepeatRule, in.RepeatEndType, in.RepeatEndDate, in.RepeatEndCount, in.TimeLimitType, in.TimeLimitStart, in.TimeLimitEnd, in.ReminderType, in.ReminderAt, in.ReminderOffsetMinutes, in.TagId); err != nil {
		return err
	}
	if _, err := parseTaskDate(in.StartDate, true, "invalid startDate"); err != nil {
		return err
	}
	if strings.TrimSpace(in.EndDate) != "" {
		if _, err := parseTaskDate(in.EndDate, false, "invalid endDate"); err != nil {
			return err
		}
	}
	return nil
}

// validateTaskUpdatePayload 校验编辑任务全部前置参数，保持与创建任务一致的枚举约束。
func validateTaskUpdatePayload(ctx context.Context, tx gdb.TX, in v1.TaskUpdateInput) error {
	return validateTaskPayload(ctx, tx, in.Title, in.Star, in.KidIds, in.CompletionMode, in.RepeatRule, in.RepeatEndType, in.RepeatEndDate, in.RepeatEndCount, in.TimeLimitType, in.TimeLimitStart, in.TimeLimitEnd, in.ReminderType, in.ReminderAt, in.ReminderOffsetMinutes, in.TagId)
}

// validateTaskPayload 校验任务核心字段、枚举字段、儿童成员和标签是否有效。
func validateTaskPayload(ctx context.Context, tx gdb.TX, title string, star int, kidIds []uint64, mode, repeatRule, repeatEndType, repeatEndDate string, repeatEndCount int, timeLimitType, timeLimitStart, timeLimitEnd, reminderType, reminderAt string, reminderOffsetMinutes int, tagId uint64) error {
	if strings.TrimSpace(title) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "title is required")
	}
	if star <= 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "star must be greater than 0")
	}
	if len(kidIds) == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "kidIds is required")
	}
	if len(uniqueUint64s(kidIds)) == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "kidIds is required")
	}
	if !validCompletionMode(mode) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported completion mode")
	}
	if !validTaskRepeatRule(repeatRule) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported repeatRule")
	}
	if !validRepeatEndType(repeatEndType) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported repeatEndType")
	}
	normalizedEndType := normalizedRepeatEndType(repeatEndType)
	if normalizedRepeatRule(repeatRule) != consts.TaskRepeatNone {
		switch normalizedEndType {
		case consts.TaskRepeatEndDate:
			if strings.TrimSpace(repeatEndDate) == "" {
				return gerror.NewCode(gcode.CodeInvalidParameter, "repeatEndDate is required")
			}
		case consts.TaskRepeatEndCount:
			if repeatEndCount <= 0 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "repeatEndCount must be greater than 0")
			}
		}
	}
	if strings.TrimSpace(repeatEndDate) != "" {
		if _, err := parseTaskDate(repeatEndDate, false, "invalid repeatEndDate"); err != nil {
			return err
		}
	}
	if repeatEndCount < 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "repeatEndCount cannot be negative")
	}
	if !validTimeLimitType(timeLimitType) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported timeLimitType")
	}
	if err := validateTaskTimeLimit(timeLimitType, timeLimitStart, timeLimitEnd); err != nil {
		return err
	}
	if !validReminderType(reminderType) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "unsupported reminderType")
	}
	if err := validateTaskReminder(reminderType, reminderAt, reminderOffsetMinutes); err != nil {
		return err
	}
	for _, kidId := range uniqueUint64s(kidIds) {
		member, err := findFamilyMember(ctx, tx, kidId)
		if err != nil || member == nil {
			return utils.NewCodef(gcode.CodeInvalidParameter, "error.kid_not_found", kidId)
		}
	}
	if tagId > 0 {
		if err := ensureTaskTagExists(ctx, tx, tagId); err != nil {
			return err
		}
	}
	return nil
}

// replaceTaskAssignees 重建任务分配儿童列表。
func replaceTaskAssignees(ctx context.Context, tx gdb.TX, taskId uint64, kidIds []uint64) error {
	if _, err := tx.Model(consts.KidsTaskAssigneeTable).Ctx(ctx).Where("task_id", taskId).Delete(); err != nil {
		return err
	}
	for index, kidId := range uniqueUint64s(kidIds) {
		data := map[string]any{"task_id": taskId, "kid_id": kidId, "assignee_order": index}
		if _, err := tx.Model(consts.KidsTaskAssigneeTable).Ctx(ctx).Data(data).Insert(); err != nil {
			return err
		}
	}
	return nil
}

// expandTaskDates 根据重复规则计算需要创建的任务日期列表。
func expandTaskDates(in v1.TaskCreateInput) ([]string, error) {
	start, err := parseTaskDate(in.StartDate, true, "invalid startDate")
	if err != nil {
		return nil, err
	}
	startText := start.Format(consts.DateLayout)
	repeatRule := normalizedRepeatRule(in.RepeatRule)
	if repeatRule == consts.TaskRepeatNone {
		return []string{startText}, nil
	}
	endType := normalizedRepeatEndType(in.RepeatEndType)
	endText := strings.TrimSpace(in.RepeatEndDate)
	if endText == "" {
		endText = strings.TrimSpace(in.EndDate)
	}
	maxEnd := start.AddDate(0, 0, consts.MaxTaskRepeatGenerateDays-1)
	if endType == consts.TaskRepeatEndDate && endText != "" {
		parsedEnd, err := parseTaskDate(endText, false, "invalid repeatEndDate")
		if err != nil {
			return nil, err
		}
		if parsedEnd.Before(maxEnd) {
			maxEnd = parsedEnd
		}
	}
	limit := consts.MaxTaskRepeatGenerateDays
	if endType == consts.TaskRepeatEndCount && in.RepeatEndCount > 0 && in.RepeatEndCount < limit {
		limit = in.RepeatEndCount
	}
	dates := make([]string, 0, limit)
	for current := start; !current.After(maxEnd) && len(dates) < limit; current = nextRepeatDate(current, repeatRule) {
		dates = append(dates, current.Format(consts.DateLayout))
		if next := nextRepeatDate(current, repeatRule); !next.After(current) {
			break
		}
	}
	if len(dates) == 0 {
		dates = append(dates, startText)
	}
	return dates, nil
}

// nextRepeatDate 按重复规则返回下一次任务日期。
func nextRepeatDate(current time.Time, rule string) time.Time {
	switch rule {
	case consts.TaskRepeatWeekly:
		return current.AddDate(0, 0, 7)
	case consts.TaskRepeatMonthly:
		return current.AddDate(0, 1, 0)
	case consts.TaskRepeatYearly:
		return current.AddDate(1, 0, 0)
	case consts.TaskRepeatCustom:
		return current.AddDate(0, 0, 1)
	default:
		return current.AddDate(0, 0, 1)
	}
}

// taskItemFromDB 将任务表记录转换为接口任务结构，并补充任务分配儿童完成状态。
func taskItemFromDB(ctx context.Context, tx gdb.TX, record gdb.Record) (v1.TaskItem, error) {
	if record.IsEmpty() {
		return v1.TaskItem{}, nil
	}
	assignees, err := taskAssigneeStatuses(ctx, tx, record)
	if err != nil {
		return v1.TaskItem{}, err
	}
	kidIds := make([]uint64, 0, len(assignees))
	for _, assignee := range assignees {
		kidIds = append(kidIds, assignee.KidId)
	}
	date := record["task_date"].String()
	return v1.TaskItem{
		Id:                    record["id"].Uint64(),
		Title:                 record["title"].String(),
		Icon:                  record["icon"].String(),
		Note:                  record["note"].String(),
		Star:                  record["star"].Int(),
		Date:                  date,
		KidIds:                kidIds,
		Assignees:             assignees,
		CompletionMode:        record["completion_mode"].String(),
		RepeatRule:            record["repeat_rule"].String(),
		RepeatEndType:         record["repeat_end_type"].String(),
		RepeatEndDate:         record["repeat_end_date"].String(),
		RepeatEndCount:        record["repeat_end_count"].Int(),
		TimeLimitType:         record["time_limit_type"].String(),
		TimeLimitStart:        record["time_limit_start"].String(),
		TimeLimitEnd:          record["time_limit_end"].String(),
		ReminderType:          record["reminder_type"].String(),
		ReminderAt:            record["reminder_at"].String(),
		ReminderOffsetMinutes: record["reminder_offset_minutes"].Int(),
		NeedPhotoProof:        record["need_photo_proof"].Int() == 1,
		TagId:                 record["tag_id"].Uint64(),
		TagName:               record["tag_name"].String(),
		TagColor:              record["tag_color"].String(),
		CanComplete:           canCompleteTaskDate(date),
		Completed:             record["completed"].Int() == 1,
		CompletedBy:           record["completed_by"].Uint64(),
		PhotoUrl:              record["photo_url"].String(),
		CompletedAt:           utils.ParseDBTime(record["completed_at"].Val()),
	}, nil
}

// taskAssigneeStatuses 查询任务分配明细，并计算轮流模式下当天应完成的儿童。
func taskAssigneeStatuses(ctx context.Context, tx gdb.TX, taskRecord gdb.Record) ([]v1.TaskAssigneeStatus, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsTaskAssigneeTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsTaskAssigneeTable).Ctx(ctx)
	}
	records, err := model.Where("task_id", taskRecord["id"].Uint64()).OrderAsc("assignee_order,id").All()
	if err != nil {
		return nil, err
	}
	activeKidId := activeKidForTask(taskRecord, records)
	items := make([]v1.TaskAssigneeStatus, 0, len(records))
	for _, record := range records {
		kidId := record["kid_id"].Uint64()
		items = append(items, v1.TaskAssigneeStatus{
			KidId:       kidId,
			Order:       record["assignee_order"].Int(),
			Active:      activeKidId == 0 || activeKidId == kidId,
			Completed:   record["completed"].Int() == 1,
			PhotoUrl:    record["photo_url"].String(),
			CompletedAt: utils.ParseDBTime(record["completed_at"].Val()),
		})
	}
	return items, nil
}

// activeKidForTask 根据任务完成模式计算当前日期下可完成任务的儿童，非轮流模式返回 0 表示全部可完成。
func activeKidForTask(taskRecord gdb.Record, assignees gdb.Result) uint64 {
	if taskRecord["completion_mode"].String() != v1.TaskCompletionModeRotation || len(assignees) == 0 {
		return 0
	}
	taskDate, err := time.ParseInLocation(consts.DateLayout, taskRecord["task_date"].String(), time.Local)
	if err != nil {
		return assignees[0]["kid_id"].Uint64()
	}
	days := int(time.Since(taskDate).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return assignees[days%len(assignees)]["kid_id"].Uint64()
}

// findTaskAssigneeStatus 在任务响应中查找指定儿童的分配完成状态。
func findTaskAssigneeStatus(task v1.TaskItem, kidId uint64) (v1.TaskAssigneeStatus, bool) {
	for _, assignee := range task.Assignees {
		if assignee.KidId == kidId {
			return assignee, true
		}
	}
	return v1.TaskAssigneeStatus{}, false
}

// shouldMarkTaskCompleted 根据完成模式判断整条任务是否应该标记为完成。
func shouldMarkTaskCompleted(ctx context.Context, tx gdb.TX, taskId uint64, mode string) (bool, error) {
	if mode != v1.TaskCompletionModeEveryone {
		return true, nil
	}
	pending, err := tx.Model(consts.KidsTaskAssigneeTable).Ctx(ctx).Where("task_id", taskId).Where("completed", 0).Count()
	if err != nil {
		return false, err
	}
	return pending == 0, nil
}

// taskVisibleForKid 判断某个孩子查询任务时是否应该看到该任务，轮流模式只展示当天轮到的孩子。
func taskVisibleForKid(task v1.TaskItem, kidId uint64) bool {
	assignee, ok := findTaskAssigneeStatus(task, kidId)
	return ok && assignee.Active
}

// matchTaskStatus 判断任务是否符合 pending/completed 状态筛选。
func matchTaskStatus(task v1.TaskItem, status string) bool {
	switch status {
	case "", "all":
		return true
	case "pending":
		return !task.Completed
	case "completed":
		return task.Completed
	default:
		return true
	}
}

// validCompletionMode 校验任务完成模式是否在产品定义范围内。
func validCompletionMode(mode string) bool {
	switch mode {
	case v1.TaskCompletionModeSingle, v1.TaskCompletionModeRotation, v1.TaskCompletionModeAnyone, v1.TaskCompletionModeEveryone:
		return true
	default:
		return false
	}
}

// validTaskRepeatRule 校验任务重复规则，空值允许走默认不重复。
func validTaskRepeatRule(rule string) bool {
	switch strings.TrimSpace(rule) {
	case "", consts.TaskRepeatNone, consts.TaskRepeatDaily, consts.TaskRepeatWeekly, consts.TaskRepeatMonthly, consts.TaskRepeatYearly, consts.TaskRepeatCustom:
		return true
	default:
		return false
	}
}

// normalizedRepeatRule 标准化重复规则，空值默认为不重复。
func normalizedRepeatRule(rule string) string {
	switch strings.TrimSpace(rule) {
	case consts.TaskRepeatNone, consts.TaskRepeatDaily, consts.TaskRepeatWeekly, consts.TaskRepeatMonthly, consts.TaskRepeatYearly, consts.TaskRepeatCustom:
		return strings.TrimSpace(rule)
	default:
		return consts.TaskRepeatNone
	}
}

// validRepeatEndType 校验重复结束类型，空值允许走默认永不结束。
func validRepeatEndType(value string) bool {
	switch strings.TrimSpace(value) {
	case "", consts.TaskRepeatEndNever, consts.TaskRepeatEndDate, consts.TaskRepeatEndCount:
		return true
	default:
		return false
	}
}

// normalizedRepeatEndType 标准化重复结束类型，空值默认为永不结束。
func normalizedRepeatEndType(value string) string {
	switch strings.TrimSpace(value) {
	case consts.TaskRepeatEndNever, consts.TaskRepeatEndDate, consts.TaskRepeatEndCount:
		return strings.TrimSpace(value)
	default:
		return consts.TaskRepeatEndNever
	}
}

// validTimeLimitType 校验任务时间限制类型，空值允许走默认全天。
func validTimeLimitType(value string) bool {
	switch strings.TrimSpace(value) {
	case "", consts.TaskTimeLimitAllDay, consts.TaskTimeLimitRange:
		return true
	default:
		return false
	}
}

// normalizedTimeLimitType 标准化时间限制类型，空值默认为全天。
func normalizedTimeLimitType(value string) string {
	if strings.TrimSpace(value) == consts.TaskTimeLimitRange {
		return consts.TaskTimeLimitRange
	}
	return consts.TaskTimeLimitAllDay
}

// validReminderType 校验提醒类型，空值允许走默认不提醒。
func validReminderType(value string) bool {
	switch strings.TrimSpace(value) {
	case "", consts.TaskReminderNone, consts.TaskReminderAtTime, consts.TaskReminderBeforeStart:
		return true
	default:
		return false
	}
}

// normalizedReminderType 标准化提醒类型，空值默认为不提醒。
func normalizedReminderType(value string) string {
	switch strings.TrimSpace(value) {
	case consts.TaskReminderNone, consts.TaskReminderAtTime, consts.TaskReminderBeforeStart:
		return strings.TrimSpace(value)
	default:
		return consts.TaskReminderNone
	}
}

// validateTaskTimeLimit 校验任务时间段配置，range 模式必须提供合法开始和结束时间。
func validateTaskTimeLimit(timeLimitType, startText, endText string) error {
	if normalizedTimeLimitType(timeLimitType) != consts.TaskTimeLimitRange {
		return nil
	}
	start, err := parseClockTime(startText, "invalid timeLimitStart")
	if err != nil {
		return err
	}
	end, err := parseClockTime(endText, "invalid timeLimitEnd")
	if err != nil {
		return err
	}
	if !end.After(start) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "timeLimitEnd must be after timeLimitStart")
	}
	return nil
}

// validateTaskReminder 校验任务提醒配置，固定时间提醒必须提供时间，提前提醒分钟数不能为负数。
func validateTaskReminder(reminderType, reminderAt string, offsetMinutes int) error {
	switch normalizedReminderType(reminderType) {
	case consts.TaskReminderAtTime:
		if _, err := parseClockTime(reminderAt, "invalid reminderAt"); err != nil {
			return err
		}
	case consts.TaskReminderBeforeStart:
		if offsetMinutes <= 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "reminderOffsetMinutes must be greater than 0")
		}
		if strings.TrimSpace(reminderAt) != "" {
			if _, err := parseClockTime(reminderAt, "invalid reminderAt"); err != nil {
				return err
			}
		}
	default:
		if offsetMinutes < 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "reminderOffsetMinutes cannot be negative")
		}
		if strings.TrimSpace(reminderAt) != "" {
			if _, err := parseClockTime(reminderAt, "invalid reminderAt"); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseTaskDate 按业务日期格式解析任务日期，allowEmpty 为 true 时空值默认今天。
func parseTaskDate(value string, allowEmpty bool, message string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" && allowEmpty {
		value = time.Now().Format(consts.DateLayout)
	}
	parsed, err := time.ParseInLocation(consts.DateLayout, value, time.Local)
	if err != nil {
		return time.Time{}, gerror.NewCode(gcode.CodeInvalidParameter, message)
	}
	return parsed, nil
}

// parseClockTime 按 HH:mm 格式解析任务时间字段。
func parseClockTime(value string, message string) (time.Time, error) {
	parsed, err := time.ParseInLocation("15:04", strings.TrimSpace(value), time.Local)
	if err != nil {
		return time.Time{}, gerror.NewCode(gcode.CodeInvalidParameter, message)
	}
	return parsed, nil
}

// nullDate 将空日期转换为数据库 NULL 值。
func nullDate(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}

// canCompleteTaskDate 判断任务日期是否允许完成，未来任务不可提前完成。
func canCompleteTaskDate(date string) bool {
	taskDate, err := time.ParseInLocation(consts.DateLayout, date, time.Local)
	if err != nil {
		return true
	}
	today, _ := time.ParseInLocation(consts.DateLayout, time.Now().Format(consts.DateLayout), time.Local)
	return !taskDate.After(today)
}

// uniqueUint64s 对 ID 列表去重并保持原有顺序。
func uniqueUint64s(values []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(values))
	items := make([]uint64, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}
	return items
}

// v1UpsertTaskTag 创建或更新当前圈子的任务标签，并按 expected_version 执行乐观锁。
func (s *sKids) v1UpsertTaskTag(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	tagID := in.PathParameters["task_tag_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks); err != nil {
		return nil, "", err
	}
	expected, hasExpected := v1ExpectedVersion(in.Body["expected_version"])
	if !hasExpected && in.Body["expected_version"] != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "expected tag version is invalid")
	}
	now := time.Now()
	var tag map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks); e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1TaskTagTable).Ctx(ctx).Where("task_tag_id", tagID).Where("circle_id", circleID).LockUpdate().One()
		if e != nil {
			return e
		}
		if !hasExpected {
			if !row.IsEmpty() {
				return v1Error(409, "IDEMPOTENCY_CONFLICT", false, "task tag already exists")
			}
			if _, e = tx.Model(consts.KidsV1TaskTagTable).Ctx(ctx).Data(gdb.Map{"task_tag_id": tagID, "circle_id": circleID, "name": in.Body["name"], "status": "active", "version": 1, "created_at": now, "updated_at": now}).Insert(); e != nil {
				return e
			}
		} else {
			if row.IsEmpty() {
				return v1Error(404, "NOT_FOUND", false, "task tag is missing")
			}
			current := row["version"].Int64()
			if current != expected {
				return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "task tag version conflicts"}
			}
			if _, e = tx.Model(consts.KidsV1TaskTagTable).Ctx(ctx).Where("task_tag_id", tagID).Data(gdb.Map{"name": in.Body["name"], "version": current + 1, "updated_at": now}).Update(); e != nil {
				return e
			}
		}
		row, e = tx.Model(consts.KidsV1TaskTagTable).Ctx(ctx).Where("task_tag_id", tagID).One()
		if e != nil {
			return e
		}
		tag = v1TaskTagRecordProjection(row)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"task_tag": tag})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "task_tag": tag}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1DeleteTaskTag 软删除任务标签并生成 tombstone；未来任务引用清理由任务领域事务负责。
func (s *sKids) v1DeleteTaskTag(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	tagID := in.PathParameters["task_tag_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "task tag version is invalid")
	}
	now := time.Now()
	var tombstone map[string]any
	var receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, e := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks)
		if e != nil {
			return e
		}
		row, e := tx.Model(consts.KidsV1TaskTagTable).Ctx(ctx).Where("task_tag_id", tagID).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if e != nil {
			return e
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "task tag is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "task tag version conflicts"}
		}
		version := current + 1
		if _, e = tx.Model(consts.KidsV1TaskTagTable).Ctx(ctx).Where("task_tag_id", tagID).Data(gdb.Map{"status": "deleted", "version": version, "deleted_at": now, "updated_at": now}).Update(); e != nil {
			return e
		}
		actor, e := v1ActorSnapshotTx(ctx, tx, membership)
		if e != nil {
			return e
		}
		tombstone = v1EntityTombstone("task_tag", tagID, version, now, actor)
		var ce error
		receipt, ce = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"tombstone": tombstone})
		return ce
	})
	if err != nil {
		return nil, "", err
	}
	return map[string]any{"receipt": receipt, "tombstone": tombstone, "active_task_references_cleared": int64(0)}, v1CommitCursor(receipt["commit_sequence"].(int64)), nil
}

// v1TaskTagRecordProjection 将任务标签事实还原为合同实体。
func v1TaskTagRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"task_tag_id": row["task_tag_id"].String(), "circle_id": row["circle_id"].String(), "name": row["name"].String(), "status": row["status"].String(), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "updated_at_ms": row["updated_at"].Time().UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(row["deleted_at"].Time())}
}

// v1UpsertTask 创建或更新任务定义，并同步未来尚未解决的首个 occurrence 快照。
func (s *sKids) v1UpsertTask(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, taskID := in.PathParameters["circle_id"], in.PathParameters["task_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks); err != nil {
		return nil, "", err
	}
	_, hasExpected := v1ExpectedVersion(in.Body["expected_version"])
	if in.Body["expected_version"] == nil && hasExpected {
		hasExpected = false
	}
	if in.Body["expected_version"] != nil && !hasExpected {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "task version is invalid")
	}
	assigned, ok := in.Body["assigned_member_ids"].([]any)
	if !ok || len(assigned) == 0 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "assigned members are invalid")
	}
	now := time.Now()
	effective := fmt.Sprint(in.Body["future_effective_from_date"])
	startDate := fmt.Sprint(in.Body["start_date"])
	var task map[string]any
	var receipt map[string]any
	var replaced int64
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks); err != nil {
			return err
		}
		row, err := tx.Model(consts.KidsV1TaskTable).Ctx(ctx).Where("task_id", taskID).Where("circle_id", circleID).LockUpdate().One()
		if err != nil {
			return err
		}
		if hasExpected && row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "task is missing")
		}
		if !hasExpected && !row.IsEmpty() {
			return v1Error(409, "IDEMPOTENCY_CONFLICT", false, "task already exists")
		}
		version, revision := int64(1), int64(1)
		if !row.IsEmpty() {
			version, revision = row["version"].Int64()+1, row["series_revision"].Int64()+1
			if _, err = tx.Model(consts.KidsV1TaskAssignmentTable).Ctx(ctx).Where("task_id", taskID).Delete(); err != nil {
				return err
			}
		}
		repeatJSON, _ := json.Marshal(in.Body["repeat_rule"])
		endJSON, _ := json.Marshal(in.Body["end_rule"])
		reminderJSON, _ := json.Marshal(in.Body["reminder_config"])
		values := gdb.Map{"task_id": taskID, "circle_id": circleID, "title": in.Body["title"], "notes": nullableV1String(in.Body["notes"]), "emoji": nullableV1String(in.Body["emoji"]), "stars": in.Body["stars"], "start_date": startDate, "zone_id": in.Body["zone_id"], "repeat_rule": string(repeatJSON), "end_rule": string(endJSON), "time_limit_minute_of_day": in.Body["time_limit_minute_of_day"], "reminder_config": string(reminderJSON), "photo_required": in.Body["photo_required"], "task_tag_id": nullableV1String(in.Body["task_tag_id"]), "series_revision": revision, "future_effective_from_date": effective, "status": "active", "version": version, "updated_at": now}
		if row.IsEmpty() {
			values["created_at"] = now
			if _, err = tx.Model(consts.KidsV1TaskTable).Ctx(ctx).Data(values).Insert(); err != nil {
				return err
			}
		} else if _, err = tx.Model(consts.KidsV1TaskTable).Ctx(ctx).Where("task_id", taskID).Data(values).Update(); err != nil {
			return err
		}
		for _, raw := range assigned {
			memberID, ok := raw.(string)
			if !ok {
				return v1Error(422, "VALIDATION_FAILED", false, "assigned member is invalid")
			}
			member, e := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").One()
			if e != nil {
				return e
			}
			if member.IsEmpty() {
				return v1Error(404, "NOT_FOUND", false, "assigned member is missing")
			}
			if _, e = tx.Model(consts.KidsV1TaskAssignmentTable).Ctx(ctx).Data(gdb.Map{"task_id": taskID, "member_id": memberID, "created_at": now, "updated_at": now}).Insert(); e != nil {
				return e
			}
		}
		if !row.IsEmpty() {
			result, e := tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("task_id", taskID).Where("scheduled_date >= ?", effective).Where("state", "pending").Delete()
			if e != nil {
				return e
			}
			replaced, e = result.RowsAffected()
			if e != nil {
				return e
			}
		}
		if _, err = tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("task_id", taskID).Where("scheduled_date", startDate).Delete(); err != nil {
			return err
		}
		for _, raw := range assigned {
			memberID := fmt.Sprint(raw)
			if _, err = tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Data(gdb.Map{"circle_id": circleID, "task_id": taskID, "member_id": memberID, "scheduled_date": startDate, "zone_id": in.Body["zone_id"], "definition_revision": revision, "title_snapshot": in.Body["title"], "notes_snapshot": nullableV1String(in.Body["notes"]), "emoji_snapshot": nullableV1String(in.Body["emoji"]), "stars_snapshot": in.Body["stars"], "photo_required_snapshot": in.Body["photo_required"], "task_tag_id_snapshot": nullableV1String(in.Body["task_tag_id"]), "state": "pending", "version": 1, "created_at": now, "updated_at": now}).Insert(); err != nil {
				return err
			}
		}
		row, err = tx.Model(consts.KidsV1TaskTable).Ctx(ctx).Where("task_id", taskID).One()
		if err != nil {
			return err
		}
		task, err = v1TaskProjectionTx(ctx, tx, row)
		if err != nil {
			return err
		}
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"task": task})
		return err
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "task": task, "occurrence_boundary": map[string]any{"effective_from_date": effective, "zone_id": in.Body["zone_id"], "pending_occurrences_replaced": replaced, "resolved_occurrences_preserved": int64(0)}}, cursor, nil
}

// v1DeleteTask 将任务标记为 deleted，并保留已解决 occurrence，删除未来 pending occurrence。
func (s *sKids) v1DeleteTask(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, taskID := in.PathParameters["circle_id"], in.PathParameters["task_id"]
	if err := v1RequireMembershipPermission(ctx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks); err != nil {
		return nil, "", err
	}
	expected, ok := v1ExpectedVersion(in.Body["expected_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "task version is invalid")
	}
	effective := fmt.Sprint(in.Body["future_effective_from_date"])
	now := time.Now()
	var tombstone, receipt map[string]any
	var removed int64
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks)
		if err != nil {
			return err
		}
		row, err := tx.Model(consts.KidsV1TaskTable).Ctx(ctx).Where("task_id", taskID).Where("circle_id", circleID).Where("status", "active").LockUpdate().One()
		if err != nil {
			return err
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "task is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "task version conflicts"}
		}
		if _, err = tx.Model(consts.KidsV1TaskTable).Ctx(ctx).Where("task_id", taskID).Data(gdb.Map{"status": "deleted", "version": current + 1, "deleted_at": now, "updated_at": now}).Update(); err != nil {
			return err
		}
		result, err := tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("task_id", taskID).Where("scheduled_date >= ?", effective).Where("state", "pending").Delete()
		if err != nil {
			return err
		}
		removed, err = result.RowsAffected()
		if err != nil {
			return err
		}
		actor, err := v1ActorSnapshotTx(ctx, tx, membership)
		if err != nil {
			return err
		}
		tombstone = v1EntityTombstone("task", taskID, current+1, now, actor)
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"tombstone": tombstone})
		return err
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "tombstone": tombstone, "occurrence_boundary": map[string]any{"effective_from_date": effective, "zone_id": in.Body["zone_id"], "pending_occurrences_replaced": removed, "resolved_occurrences_preserved": int64(0)}}, cursor, nil
}

// v1TaskProjectionTx 将任务定义和分配表投影为合同 TaskDefinition。
func v1TaskProjectionTx(ctx context.Context, tx gdb.TX, row gdb.Record) (map[string]any, error) {
	assignments, err := tx.Model(consts.KidsV1TaskAssignmentTable).Ctx(ctx).Where("task_id", row["task_id"].String()).Order("id ASC").All()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(assignments))
	for _, item := range assignments {
		ids = append(ids, item["member_id"].String())
	}
	return map[string]any{"task_id": row["task_id"].String(), "circle_id": row["circle_id"].String(), "title": row["title"].String(), "notes": nullableString(row["notes"].String()), "emoji": nullableString(row["emoji"].String()), "stars": row["stars"].Int64(), "start_date": row["start_date"].Time().Format(consts.DateLayout), "zone_id": row["zone_id"].String(), "repeat_rule": v1JSONValue(row["repeat_rule"].String()), "end_rule": v1JSONValue(row["end_rule"].String()), "time_limit_minute_of_day": nullableInt64(row["time_limit_minute_of_day"].Int64()), "reminder_config": v1JSONValue(row["reminder_config"].String()), "photo_required": row["photo_required"].Bool(), "task_tag_id": nullableString(row["task_tag_id"].String()), "assigned_member_ids": ids, "series_revision": row["series_revision"].Int64(), "future_effective_from_date": row["future_effective_from_date"].Time().Format(consts.DateLayout), "status": row["status"].String(), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "updated_at_ms": row["updated_at"].Time().UnixMilli(), "deleted_at_ms": v1NullableTimeMillis(row["deleted_at"].Time())}, nil
}

// nullableInt64 将数据库零值转换为合同 nullable 整数。
func nullableInt64(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

// v1TaskOccurrenceProjection 将 occurrence 数据库记录投影为合同实体。
func v1TaskOccurrenceProjection(row gdb.Record) map[string]any {
	return map[string]any{"circle_id": row["circle_id"].String(), "task_id": row["task_id"].String(), "member_id": row["member_id"].String(), "scheduled_date": row["scheduled_date"].Time().Format(consts.DateLayout), "zone_id": row["zone_id"].String(), "definition_revision": row["definition_revision"].Int64(), "title_snapshot": row["title_snapshot"].String(), "notes_snapshot": nullableString(row["notes_snapshot"].String()), "emoji_snapshot": nullableString(row["emoji_snapshot"].String()), "stars_snapshot": row["stars_snapshot"].Int64(), "photo_required_snapshot": row["photo_required_snapshot"].Bool(), "task_tag_id_snapshot": nullableString(row["task_tag_id_snapshot"].String()), "state": row["state"].String(), "completion_id": nullableString(row["completion_id"].String()), "version": row["version"].Int64(), "created_at_ms": row["created_at"].Time().UnixMilli(), "updated_at_ms": row["updated_at"].Time().UnixMilli()}
}

// v1CompleteTask 原子完成 occurrence，并追加正向流水和余额投影。
func (s *sKids) v1CompleteTask(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	expected, ok := v1ExpectedVersion(in.Body["expected_occurrence_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "occurrence version is invalid")
	}
	now := time.Now()
	var occurrence, completion, ledger, balance, receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, "")
		if err != nil {
			return err
		}
		memberID, taskID, date := fmt.Sprint(in.Body["member_id"]), fmt.Sprint(in.Body["task_id"]), fmt.Sprint(in.Body["scheduled_date"])
		if membership["actor_type"].String() == "member" && membership["actor_id"].String() != memberID {
			return v1Error(422, "NOT_ASSIGNED", false, "member is not assigned to this completion")
		}
		row, err := tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("circle_id", circleID).Where("task_id", taskID).Where("member_id", memberID).Where("scheduled_date", date).Where("zone_id", in.Body["zone_id"]).LockUpdate().One()
		if err != nil {
			return err
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "task occurrence is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "occurrence version conflicts"}
		}
		if row["state"].String() != "pending" {
			return v1Error(409, "VERSION_CONFLICT", false, "task occurrence is already resolved")
		}
		proofID := nullableV1String(in.Body["proof_asset_id"])
		if row["photo_required_snapshot"].Bool() && proofID == "" {
			return v1Error(422, "PHOTO_REQUIRED", false, "task proof is required")
		}
		actor, err := v1ActorSnapshotTx(ctx, tx, membership)
		if err != nil {
			return err
		}
		completionID := fmt.Sprint(in.Body["completion_id"])
		if _, err = tx.Model(consts.KidsV1TaskCompletionTable).Ctx(ctx).Data(gdb.Map{"completion_id": completionID, "circle_id": circleID, "task_id": taskID, "member_id": memberID, "scheduled_date": date, "zone_id": in.Body["zone_id"], "proof_asset_id": proofID, "title_snapshot": row["title_snapshot"], "stars_snapshot": row["stars_snapshot"], "completed_by": mustV1JSON(actor), "completed_at": now, "commit_sequence": 0, "version": 1}).Insert(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("id", row["id"].Int64()).Data(gdb.Map{"state": "completed", "completion_id": completionID, "version": current + 1, "updated_at": now}).Update(); err != nil {
			return err
		}
		ledgerID := v1ID("ledger", uuid.NewString())
		source := map[string]any{"source_type": "task", "source_id": taskID, "title_snapshot": row["title_snapshot"].String(), "stars_snapshot": row["stars_snapshot"].Int64(), "asset_id_snapshot": nullableString(proofID), "scheduled_date_snapshot": date}
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Data(gdb.Map{"ledger_id": ledgerID, "circle_id": circleID, "member_id": memberID, "source": mustV1JSON(source), "delta": row["stars_snapshot"].Int64(), "actor": mustV1JSON(actor), "commit_sequence": 0, "created_at": now}).Insert(); err != nil {
			return err
		}
		balanceRow, err := tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).LockUpdate().One()
		if err != nil {
			return err
		}
		nextBalance, balanceVersion := row["stars_snapshot"].Int64(), int64(1)
		if !balanceRow.IsEmpty() {
			nextBalance += balanceRow["balance"].Int64()
			balanceVersion = balanceRow["version"].Int64() + 1
		}
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"completion_id": completionID, "ledger_id": ledgerID})
		if err != nil {
			return err
		}
		sequence := receipt["commit_sequence"].(int64)
		commitID := receipt["commit_id"].(string)
		if _, err = tx.Model(consts.KidsV1TaskCompletionTable).Ctx(ctx).Where("completion_id", completionID).Data(gdb.Map{"commit_sequence": sequence}).Update(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("ledger_id", ledgerID).Data(gdb.Map{"commit_sequence": sequence}).Update(); err != nil {
			return err
		}
		data := gdb.Map{"balance": nextBalance, "version": balanceVersion, "source_commit_id": commitID, "source_commit_sequence": sequence, "updated_at": now}
		if balanceRow.IsEmpty() {
			data["circle_id"], data["member_id"] = circleID, memberID
			if _, err = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Data(data).Insert(); err != nil {
				return err
			}
		} else if _, err = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("id", balanceRow["id"].Int64()).Data(data).Update(); err != nil {
			return err
		}
		row, err = tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("id", row["id"].Int64()).One()
		if err != nil {
			return err
		}
		occurrence = v1TaskOccurrenceProjection(row)
		completion = v1CompletionProjection(completionID, circleID, taskID, memberID, date, fmt.Sprint(in.Body["zone_id"]), proofID, row["title_snapshot"].String(), row["stars_snapshot"].Int64(), actor, now, sequence, 1)
		ledger = v1LedgerProjection(ledgerID, circleID, memberID, source, row["stars_snapshot"].Int64(), nil, actor, nil, now, sequence)
		balance = v1BalanceProjection(circleID, memberID, nextBalance, balanceVersion, commitID, sequence, now)
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "occurrence": occurrence, "completion": completion, "ledger_entry": ledger, "balance": balance, "change_cursor": cursor}, cursor, nil
}

// v1CompletionProjection 构造 append-only completion 的合同投影。
func v1CompletionProjection(completionID, circleID, taskID, memberID, date, zoneID, proofID, title string, stars int64, actor map[string]any, at time.Time, sequence, version int64) map[string]any {
	return map[string]any{"completion_id": completionID, "circle_id": circleID, "task_id": taskID, "member_id": memberID, "scheduled_date": date, "zone_id": zoneID, "completed_at_ms": at.UnixMilli(), "proof_asset_id": nullableString(proofID), "title_snapshot": title, "stars_snapshot": stars, "completed_by": actor, "version": version, "commit_sequence": sequence}
}

// v1LedgerProjection 构造 append-only 星星流水的合同投影。
func v1LedgerProjection(ledgerID, circleID, memberID string, source map[string]any, delta int64, reason any, actor map[string]any, reversal any, at time.Time, sequence int64) map[string]any {
	return map[string]any{"ledger_id": ledgerID, "circle_id": circleID, "member_id": memberID, "source": source, "delta": delta, "reason": reason, "actor": actor, "reversal_of_ledger_id": reversal, "created_at_ms": at.UnixMilli(), "commit_sequence": sequence}
}

// v1BalanceProjection 构造成员当前余额的合同投影。
func v1BalanceProjection(circleID, memberID string, amount, version int64, commitID string, sequence int64, at time.Time) map[string]any {
	return map[string]any{"circle_id": circleID, "member_id": memberID, "balance": amount, "version": version, "source_commit_id": commitID, "source_commit_sequence": sequence, "updated_at_ms": at.UnixMilli()}
}

// v1CancelTaskCompletion 保留原 completion 和正向流水，并原子追加 cancellation 与冲销流水。
func (s *sKids) v1CancelTaskCompletion(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, completionID := in.PathParameters["circle_id"], in.PathParameters["completion_id"]
	expected, ok := v1ExpectedVersion(in.Body["expected_completion_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "completion version is invalid")
	}
	now := time.Now()
	var occurrence, completion, cancellation, original, reversal, balance, receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionManageTasks)
		if err != nil {
			return err
		}
		row, err := tx.Model(consts.KidsV1TaskCompletionTable).Ctx(ctx).Where("completion_id", completionID).Where("circle_id", circleID).LockUpdate().One()
		if err != nil {
			return err
		}
		if row.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "completion is missing")
		}
		current := row["version"].Int64()
		if current != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &current, Message: "completion version conflicts"}
		}
		exists, err := tx.Model(consts.KidsV1TaskCancellationTable).Ctx(ctx).Where("completion_id", completionID).One()
		if err != nil {
			return err
		}
		if !exists.IsEmpty() {
			return v1Error(409, "COMPLETION_ALREADY_CANCELLED", false, "completion is already cancelled")
		}
		occurrenceRow, err := tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("task_id", row["task_id"].String()).Where("member_id", row["member_id"].String()).Where("scheduled_date", row["scheduled_date"].Time().Format(consts.DateLayout)).LockUpdate().One()
		if err != nil {
			return err
		}
		if occurrenceRow.IsEmpty() {
			return v1Error(409, "AUDIT_INCONSISTENT", false, "completion occurrence is unavailable")
		}
		originalRow, err := tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", row["member_id"].String()).Where("delta", row["stars_snapshot"].Int64()).Order("id ASC").One()
		if err != nil {
			return err
		}
		if originalRow.IsEmpty() {
			return v1Error(409, "AUDIT_INCONSISTENT", false, "completion ledger is unavailable")
		}
		balanceRow, err := tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", row["member_id"].String()).LockUpdate().One()
		if err != nil {
			return err
		}
		if balanceRow.IsEmpty() || balanceRow["balance"].Int64() < row["stars_snapshot"].Int64() {
			return v1Error(409, "AUDIT_INCONSISTENT", false, "balance cannot be reversed")
		}
		actor, err := v1ActorSnapshotTx(ctx, tx, membership)
		if err != nil {
			return err
		}
		cancellationID := fmt.Sprint(in.Body["cancellation_id"])
		reversalID := v1ID("ledger", uuid.NewString())
		source := map[string]any{"source_type": "task", "source_id": row["task_id"].String(), "title_snapshot": row["title_snapshot"].String(), "stars_snapshot": row["stars_snapshot"].Int64(), "asset_id_snapshot": nullableString(row["proof_asset_id"].String()), "scheduled_date_snapshot": row["scheduled_date"].Time().Format(consts.DateLayout)}
		if _, err = tx.Model(consts.KidsV1TaskCancellationTable).Ctx(ctx).Data(gdb.Map{"cancellation_id": cancellationID, "completion_id": completionID, "reason_code": in.Body["reason_code"], "cancelled_by": mustV1JSON(actor), "cancelled_at": now, "commit_sequence": 0}).Insert(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Data(gdb.Map{"ledger_id": reversalID, "circle_id": circleID, "member_id": row["member_id"], "source": mustV1JSON(source), "delta": -row["stars_snapshot"].Int64(), "reason": in.Body["reason_code"], "actor": mustV1JSON(actor), "reversal_of_ledger_id": originalRow["ledger_id"], "commit_sequence": 0, "created_at": now}).Insert(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("id", occurrenceRow["id"].Int64()).Data(gdb.Map{"state": "cancelled", "version": occurrenceRow["version"].Int64() + 1, "updated_at": now}).Update(); err != nil {
			return err
		}
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"completion_id": completionID, "cancellation_id": cancellationID, "ledger_id": reversalID})
		if err != nil {
			return err
		}
		sequence, commitID := receipt["commit_sequence"].(int64), receipt["commit_id"].(string)
		if _, err = tx.Model(consts.KidsV1TaskCancellationTable).Ctx(ctx).Where("cancellation_id", cancellationID).Data(gdb.Map{"commit_sequence": sequence}).Update(); err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("ledger_id", reversalID).Data(gdb.Map{"commit_sequence": sequence}).Update(); err != nil {
			return err
		}
		amount, version := balanceRow["balance"].Int64()-row["stars_snapshot"].Int64(), balanceRow["version"].Int64()+1
		if _, err = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("id", balanceRow["id"].Int64()).Data(gdb.Map{"balance": amount, "version": version, "source_commit_id": commitID, "source_commit_sequence": sequence, "updated_at": now}).Update(); err != nil {
			return err
		}
		occurrenceRow, err = tx.Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).Where("id", occurrenceRow["id"].Int64()).One()
		if err != nil {
			return err
		}
		occurrence = v1TaskOccurrenceProjection(occurrenceRow)
		completedBy := v1JSONValue(row["completed_by"].String()).(map[string]any)
		completion = v1CompletionProjection(completionID, circleID, row["task_id"].String(), row["member_id"].String(), row["scheduled_date"].Time().Format(consts.DateLayout), row["zone_id"].String(), row["proof_asset_id"].String(), row["title_snapshot"].String(), row["stars_snapshot"].Int64(), completedBy, row["completed_at"].Time(), row["commit_sequence"].Int64(), current)
		cancellation = map[string]any{"cancellation_id": cancellationID, "completion_id": completionID, "cancelled_at_ms": now.UnixMilli(), "cancelled_by": actor, "reason_code": fmt.Sprint(in.Body["reason_code"]), "commit_sequence": sequence}
		original = v1LedgerRecordProjection(originalRow)
		reversal = v1LedgerProjection(reversalID, circleID, row["member_id"].String(), source, -row["stars_snapshot"].Int64(), fmt.Sprint(in.Body["reason_code"]), actor, originalRow["ledger_id"].String(), now, sequence)
		balance = v1BalanceProjection(circleID, row["member_id"].String(), amount, version, commitID, sequence, now)
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "occurrence": occurrence, "completion": completion, "cancellation": cancellation, "original_ledger_entry": original, "reversal_ledger_entry": reversal, "balance": balance, "change_cursor": cursor}, cursor, nil
}

// v1LedgerRecordProjection 将持久化流水还原为合同 LedgerEntry。
func v1LedgerRecordProjection(row gdb.Record) map[string]any {
	return map[string]any{"ledger_id": row["ledger_id"].String(), "circle_id": row["circle_id"].String(), "member_id": row["member_id"].String(), "source": v1JSONValue(row["source"].String()), "delta": row["delta"].Int64(), "reason": nullableString(row["reason"].String()), "actor": v1JSONValue(row["actor"].String()), "reversal_of_ledger_id": nullableString(row["reversal_of_ledger_id"].String()), "created_at_ms": row["created_at"].Time().UnixMilli(), "commit_sequence": row["commit_sequence"].Int64()}
}

// v1AdjustMemberStars 原子追加管理员调整流水并更新成员余额。
func (s *sKids) v1AdjustMemberStars(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID, memberID := in.PathParameters["circle_id"], fmt.Sprint(in.Body["member_id"])
	expected, ok := v1ExpectedVersion(in.Body["expected_balance_version"])
	if !ok {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "balance version is invalid")
	}
	delta, ok := v1Integer(in.Body["delta"])
	if !ok || delta == 0 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "adjustment delta is invalid")
	}
	now := time.Now()
	var ledger, balance, receipt map[string]any
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		membership, err := v1RequireMembershipTx(ctx, tx, in.PrincipalID, circleID, consts.KidsV1PermissionAdjustStars)
		if err != nil {
			return err
		}
		member, err := tx.Model(consts.KidsV1MemberTable).Ctx(ctx).Where("member_id", memberID).Where("circle_id", circleID).Where("status", "active").One()
		if err != nil {
			return err
		}
		if member.IsEmpty() {
			return v1Error(404, "NOT_FOUND", false, "member is missing")
		}
		row, err := tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", memberID).LockUpdate().One()
		if err != nil {
			return err
		}
		current, version := int64(0), int64(1)
		if !row.IsEmpty() {
			current, version = row["balance"].Int64(), row["version"].Int64()
		}
		if version != expected {
			return &v1.V1Error{Status: 409, Code: "VERSION_CONFLICT", Version: &version, Message: "balance version conflicts"}
		}
		next := current + delta
		if next < 0 || next > 2147483647 {
			return v1Error(422, "BALANCE_OVERFLOW", false, "adjustment exceeds balance range")
		}
		actor, err := v1ActorSnapshotTx(ctx, tx, membership)
		if err != nil {
			return err
		}
		ledgerID := v1ID("ledger", uuid.NewString())
		source := map[string]any{"source_type": "adjustment", "source_id": fmt.Sprint(in.Body["adjustment_id"]), "title_snapshot": nil, "stars_snapshot": nil, "asset_id_snapshot": nil, "scheduled_date_snapshot": nil}
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Data(gdb.Map{"ledger_id": ledgerID, "circle_id": circleID, "member_id": memberID, "source": mustV1JSON(source), "delta": delta, "reason": in.Body["reason"], "actor": mustV1JSON(actor), "commit_sequence": 0, "created_at": now}).Insert(); err != nil {
			return err
		}
		receipt, err = v1CreateCommitTx(ctx, tx, circleID, in.OperationID, map[string]any{"ledger_id": ledgerID})
		if err != nil {
			return err
		}
		sequence, commitID := receipt["commit_sequence"].(int64), receipt["commit_id"].(string)
		if _, err = tx.Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("ledger_id", ledgerID).Data(gdb.Map{"commit_sequence": sequence}).Update(); err != nil {
			return err
		}
		nextVersion := version + 1
		values := gdb.Map{"balance": next, "version": nextVersion, "source_commit_id": commitID, "source_commit_sequence": sequence, "updated_at": now}
		if row.IsEmpty() {
			values["circle_id"], values["member_id"] = circleID, memberID
			if _, err = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Data(values).Insert(); err != nil {
				return err
			}
		} else if _, err = tx.Model(consts.KidsV1BalanceTable).Ctx(ctx).Where("id", row["id"].Int64()).Data(values).Update(); err != nil {
			return err
		}
		ledger = v1LedgerProjection(ledgerID, circleID, memberID, source, delta, fmt.Sprint(in.Body["reason"]), actor, nil, now, sequence)
		balance = v1BalanceProjection(circleID, memberID, next, nextVersion, commitID, sequence, now)
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	cursor := v1CommitCursor(receipt["commit_sequence"].(int64))
	return map[string]any{"receipt": receipt, "ledger_entry": ledger, "balance": balance, "change_cursor": cursor}, cursor, nil
}

// v1Integer 将经过 schema 校验的 JSON 整数转为 int64。
func v1Integer(value any) (int64, bool) {
	switch typed := value.(type) {
	case json.Number:
		parsed, err := typed.Int64()
		return parsed, err == nil
	case float64:
		if typed == float64(int64(typed)) {
			return int64(typed), true
		}
	case int64:
		return typed, true
	case int:
		return int64(typed), true
	}
	return 0, false
}

// v1TaskOccurrences 按日期窗口返回任务 occurrence 的稳定分页快照。
func (s *sKids) v1TaskOccurrences(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	limit, err := strconv.Atoi(v1QueryFirst(in.Query, "limit"))
	if err != nil || limit < 1 || limit > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "occurrence limit is invalid")
	}
	offset, err := v1PageOffset(v1QueryFirst(in.Query, "cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "occurrence cursor is invalid")
	}
	model := utils.KidsDB(ctx).Model(consts.KidsV1TaskOccurrenceTable).Ctx(ctx).
		Where("circle_id", circleID).
		Where("scheduled_date >= ?", v1QueryFirst(in.Query, "start_date")).
		Where("scheduled_date < ?", v1QueryFirst(in.Query, "end_date_exclusive")).
		Where("zone_id", v1QueryFirst(in.Query, "zone_id"))
	if memberID := v1QueryFirst(in.Query, "member_id"); memberID != "" {
		model = model.Where("member_id", memberID)
	}
	rows, err := model.Order("scheduled_date ASC,id ASC").Limit(offset, limit+1).All()
	if err != nil {
		return nil, "", err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, v1TaskOccurrenceProjection(row))
	}
	cursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	next := any(nil)
	if hasMore {
		next = v1PageCursor(offset + limit)
	}
	return map[string]any{"items": items, "next_cursor": next, "has_more": hasMore, "snapshot_cursor": cursor}, cursor, nil
}

// v1CompletionDetails 返回完成和取消流水的关联审计明细。
func (s *sKids) v1CompletionDetails(ctx context.Context, in v1.V1OperationInput) (map[string]any, string, error) {
	circleID := in.PathParameters["circle_id"]
	if _, err := v1RequireMembership(ctx, in.PrincipalID, circleID, ""); err != nil {
		return nil, "", err
	}
	limit, err := strconv.Atoi(v1QueryFirst(in.Query, "limit"))
	if err != nil || limit < 1 || limit > 200 {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "completion limit is invalid")
	}
	offset, err := v1PageOffset(v1QueryFirst(in.Query, "cursor"))
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "completion cursor is invalid")
	}
	start, err := strconv.ParseInt(v1QueryFirst(in.Query, "start_at_ms"), 10, 64)
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "completion start is invalid")
	}
	end, err := strconv.ParseInt(v1QueryFirst(in.Query, "end_at_ms"), 10, 64)
	if err != nil {
		return nil, "", v1Error(422, "VALIDATION_FAILED", false, "completion end is invalid")
	}
	model := utils.KidsDB(ctx).Model(consts.KidsV1TaskCompletionTable).Ctx(ctx).
		Where("circle_id", circleID).
		Where("member_id", v1QueryFirst(in.Query, "member_id")).
		Where("zone_id", v1QueryFirst(in.Query, "zone_id")).
		Where("completed_at >= ?", time.UnixMilli(start)).
		Where("completed_at < ?", time.UnixMilli(end))
	if taskID := v1QueryFirst(in.Query, "task_id"); taskID != "" {
		model = model.Where("task_id", taskID)
	}
	rows, err := model.Order("completed_at DESC,id DESC").Limit(offset, limit+1).All()
	if err != nil {
		return nil, "", err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		completedBy, ok := v1JSONValue(row["completed_by"].String()).(map[string]any)
		if !ok {
			return nil, "", v1Error(409, "AUDIT_INCONSISTENT", false, "completion actor is unavailable")
		}
		completion := v1CompletionProjection(row["completion_id"].String(), circleID, row["task_id"].String(), row["member_id"].String(), row["scheduled_date"].Time().Format(consts.DateLayout), row["zone_id"].String(), row["proof_asset_id"].String(), row["title_snapshot"].String(), row["stars_snapshot"].Int64(), completedBy, row["completed_at"].Time(), row["commit_sequence"].Int64(), row["version"].Int64())
		positive, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("circle_id", circleID).Where("member_id", row["member_id"].String()).Where("commit_sequence", row["commit_sequence"].Int64()).Where("delta > 0").One()
		if queryErr != nil {
			return nil, "", queryErr
		}
		if positive.IsEmpty() {
			return nil, "", v1Error(409, "AUDIT_INCONSISTENT", false, "completion ledger is unavailable")
		}
		detail := map[string]any{"completion": completion, "cancellation": nil, "positive_ledger_entry": v1LedgerRecordProjection(positive), "reversal_ledger_entry": nil}
		cancelled, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1TaskCancellationTable).Ctx(ctx).Where("completion_id", row["completion_id"].String()).One()
		if queryErr != nil {
			return nil, "", queryErr
		}
		if !cancelled.IsEmpty() {
			actor, ok := v1JSONValue(cancelled["cancelled_by"].String()).(map[string]any)
			if !ok {
				return nil, "", v1Error(409, "AUDIT_INCONSISTENT", false, "cancellation actor is unavailable")
			}
			detail["cancellation"] = map[string]any{"cancellation_id": cancelled["cancellation_id"].String(), "completion_id": row["completion_id"].String(), "cancelled_at_ms": cancelled["cancelled_at"].Time().UnixMilli(), "cancelled_by": actor, "reason_code": cancelled["reason_code"].String(), "commit_sequence": cancelled["commit_sequence"].Int64()}
			reversal, queryErr := utils.KidsDB(ctx).Model(consts.KidsV1LedgerTable).Ctx(ctx).Where("reversal_of_ledger_id", positive["ledger_id"].String()).One()
			if queryErr != nil {
				return nil, "", queryErr
			}
			if reversal.IsEmpty() {
				return nil, "", v1Error(409, "AUDIT_INCONSISTENT", false, "reversal ledger is unavailable")
			}
			detail["reversal_ledger_entry"] = v1LedgerRecordProjection(reversal)
		}
		items = append(items, detail)
	}
	cursor, err := v1LatestCursor(ctx)
	if err != nil {
		return nil, "", err
	}
	next := any(nil)
	if hasMore {
		next = v1PageCursor(offset + limit)
	}
	return map[string]any{"items": items, "next_cursor": next, "has_more": hasMore, "snapshot_cursor": cursor}, cursor, nil
}
