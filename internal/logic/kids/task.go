package kids

import (
	"context"
	"strings"
	"time"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
