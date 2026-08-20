package kids

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// ListTaskTags 从数据库读取未删除的任务标签列表。
func (s *sKids) ListTaskTags(ctx context.Context, in v1.TaskTagListInput) (*v1.TaskTagListOutput, error) {
	records, err := utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx).
		Where("deleted_at IS NULL").
		OrderAsc("sort_order,id").
		All()
	if err != nil {
		return nil, err
	}
	out := &v1.TaskTagListOutput{}
	for _, record := range records {
		out.List = append(out.List, taskTagFromDB(record))
	}
	return out, nil
}

// CreateTaskTag 校验标签名称并持久化创建任务标签。
func (s *sKids) CreateTaskTag(ctx context.Context, in v1.TaskTagCreateInput) (*v1.TaskTagCreateOutput, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "tag name is required")
	}
	id, err := utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx).Data(map[string]any{
		"name":       strings.TrimSpace(in.Name),
		"color":      strings.TrimSpace(in.Color),
		"sort_order": in.SortOrder,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	record, err := utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx).Where("id", id).One()
	if err != nil {
		return nil, err
	}
	return &v1.TaskTagCreateOutput{Tag: taskTagFromDB(record)}, nil
}

// UpdateTaskTag 校验标签名称并更新任务标签。
func (s *sKids) UpdateTaskTag(ctx context.Context, in v1.TaskTagUpdateInput) (*v1.TaskTagUpdateOutput, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "tag name is required")
	}
	result, err := utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx).
		Where("id", in.Id).
		Where("deleted_at IS NULL").
		Data(map[string]any{"name": strings.TrimSpace(in.Name), "color": strings.TrimSpace(in.Color), "sort_order": in.SortOrder}).
		Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "tag not found")
	}
	record, err := utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx).Where("id", in.Id).One()
	if err != nil {
		return nil, err
	}
	return &v1.TaskTagUpdateOutput{Tag: taskTagFromDB(record)}, nil
}

// DeleteTaskTag 软删除任务标签，历史任务仍保留标签 ID。
func (s *sKids) DeleteTaskTag(ctx context.Context, in v1.TaskTagDeleteInput) (*v1.TaskTagDeleteOutput, error) {
	result, err := utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx).
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
		return nil, gerror.NewCode(gcode.CodeNotFound, "tag not found")
	}
	return &v1.TaskTagDeleteOutput{Id: in.Id}, nil
}

// ensureTaskTagExists 校验任务标签存在且未被删除。
func ensureTaskTagExists(ctx context.Context, tx gdb.TX, tagId uint64) error {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsTaskTagTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsTaskTagTable).Ctx(ctx)
	}
	count, err := model.Where("id", tagId).Where("deleted_at IS NULL").Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "tag not found")
	}
	return nil
}

// taskTagFromDB 将数据库标签记录转换为接口响应结构。
func taskTagFromDB(record gdb.Record) v1.TaskTag {
	return v1.TaskTag{
		Id:        record["id"].Uint64(),
		Name:      record["name"].String(),
		Color:     record["color"].String(),
		SortOrder: record["sort_order"].Int(),
	}
}
