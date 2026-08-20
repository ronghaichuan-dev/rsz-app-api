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

// ListFamilyMembers 从数据库读取家庭成员，并按成人和儿童分组返回。
func (s *sKids) ListFamilyMembers(ctx context.Context, in v1.FamilyMemberListInput) (*v1.FamilyMemberListOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx).Where("deleted_at IS NULL").OrderAsc("id")
	if in.CircleId > 0 {
		model = model.Where("circle_id", in.CircleId)
	}
	records, err := model.All()
	if err != nil {
		return nil, err
	}
	out := &v1.FamilyMemberListOutput{}
	for _, record := range records {
		out.Members = append(out.Members, familyMemberFromDB(ctx, record))
	}
	return out, nil
}

// CreateFamilyMember 校验并持久化一个家庭成员。
func (s *sKids) CreateFamilyMember(ctx context.Context, in v1.FamilyMemberCreateInput) (*v1.FamilyMemberCreateOutput, error) {
	if in.CircleId == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "circleId is required")
	}
	if err := validateFamilyMemberPayload(in.Name, in.Gender); err != nil {
		return nil, err
	}
	id, err := utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx).Data(map[string]any{
		"circle_id":    in.CircleId,
		"name":         strings.TrimSpace(in.Name),
		"gender":       strings.TrimSpace(in.Gender),
		"avatar":       strings.TrimSpace(in.Avatar),
		"avatar_style": strings.TrimSpace(in.AvatarStyle),
		"relation":     strings.TrimSpace(in.Relation),
		"owner":        utils.BoolToInt(in.Owner),
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	member, err := findFamilyMember(ctx, nil, uint64(id))
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberCreateOutput{Member: familyMemberRecordToAPI(ctx, member)}, nil
}

// familyMemberRecord 承载 kids_family_member 表字段。
type familyMemberRecord struct {
	Id          uint64
	CircleId    uint64
	Name        string
	Gender      string
	Avatar      string
	AvatarStyle string
	Relation    string
	Owner       bool
	BindUserId  uint64
	BoundAt     int64
}

// findFamilyMember 根据成员 ID 查询家庭成员，传入事务时使用事务连接。
func findFamilyMember(ctx context.Context, tx gdb.TX, id uint64) (*familyMemberRecord, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsFamilyMemberTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx)
	}
	record, err := model.Where("id", id).Where("deleted_at IS NULL").One()
	if err != nil || record.IsEmpty() {
		return nil, err
	}
	return familyMemberRecordFromDB(record), nil
}

// familyMemberFromDB 将数据库记录转换为接口家庭成员结构。
func familyMemberFromDB(ctx context.Context, record gdb.Record) v1.FamilyMember {
	return familyMemberRecordToAPI(ctx, familyMemberRecordFromDB(record))
}

// familyMemberRecordFromDB 将数据库记录转换为内部家庭成员结构。
func familyMemberRecordFromDB(record gdb.Record) *familyMemberRecord {
	return &familyMemberRecord{
		Id:          record["id"].Uint64(),
		CircleId:    record["circle_id"].Uint64(),
		Name:        record["name"].String(),
		Gender:      record["gender"].String(),
		Avatar:      record["avatar"].String(),
		AvatarStyle: record["avatar_style"].String(),
		Relation:    record["relation"].String(),
		Owner:       record["owner"].Int() == 1,
		BindUserId:  record["bind_user_id"].Uint64(),
		BoundAt:     utils.ParseDBTime(record["bound_at"].Val()),
	}
}

// familyMemberRecordToAPI 将内部家庭成员结构转换为接口响应结构。
func familyMemberRecordToAPI(ctx context.Context, member *familyMemberRecord) v1.FamilyMember {
	if member == nil {
		return v1.FamilyMember{}
	}
	return v1.FamilyMember{
		Id:          member.Id,
		CircleId:    member.CircleId,
		Name:        member.Name,
		Gender:      member.Gender,
		Avatar:      member.Avatar,
		AvatarStyle: member.AvatarStyle,
		Relation:    member.Relation,
		Owner:       member.Owner,
		BindUserId:  member.BindUserId,
		Bound:       member.BindUserId > 0,
		BoundAt:     member.BoundAt,
		StarCount:   getStarBalanceValue(ctx, member.Id),
	}
}

// GetFamilyMember 查询家庭成员详情。
func (s *sKids) GetFamilyMember(ctx context.Context, in v1.FamilyMemberDetailInput) (*v1.FamilyMemberDetailOutput, error) {
	member, err := findFamilyMember(ctx, nil, in.Id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "family member not found")
	}
	if ok, err := userCanAccessCircle(ctx, nil, in.UserId, member.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "circle permission denied")
	}
	return &v1.FamilyMemberDetailOutput{Member: familyMemberRecordToAPI(ctx, member)}, nil
}

// UpdateFamilyMember 校验管理权限后更新家庭成员资料。
func (s *sKids) UpdateFamilyMember(ctx context.Context, in v1.FamilyMemberUpdateInput) (*v1.FamilyMemberUpdateOutput, error) {
	member, err := findFamilyMember(ctx, nil, in.Id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "family member not found")
	}
	if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, member.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can manage member")
	}
	if err := validateFamilyMemberPayload(in.Name, in.Gender); err != nil {
		return nil, err
	}
	_, err = utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx).Where("id", in.Id).Where("deleted_at IS NULL").Data(map[string]any{
		"name":         strings.TrimSpace(in.Name),
		"gender":       strings.TrimSpace(in.Gender),
		"avatar":       strings.TrimSpace(in.Avatar),
		"avatar_style": strings.TrimSpace(in.AvatarStyle),
		"relation":     strings.TrimSpace(in.Relation),
		"owner":        utils.BoolToInt(in.Owner),
	}).Update()
	if err != nil {
		return nil, err
	}
	member, err = findFamilyMember(ctx, nil, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberUpdateOutput{Member: familyMemberRecordToAPI(ctx, member)}, nil
}

// DeleteFamilyMember 校验管理权限后软删除家庭成员。
func (s *sKids) DeleteFamilyMember(ctx context.Context, in v1.FamilyMemberDeleteInput) (*v1.FamilyMemberDeleteOutput, error) {
	member, err := findFamilyMember(ctx, nil, in.Id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "family member not found")
	}
	if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, member.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can manage member")
	}
	_, err = utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx).
		Where("id", in.Id).
		Where("deleted_at IS NULL").
		Data(map[string]any{"deleted_at": time.Now().Format(consts.MySQLTimeLayout)}).
		Update()
	if err != nil {
		return nil, err
	}
	return &v1.FamilyMemberDeleteOutput{Id: in.Id}, nil
}

// validateFamilyMemberPayload 校验家庭成员名称和性别枚举。
func validateFamilyMemberPayload(name, gender string) error {
	if strings.TrimSpace(name) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "name is required")
	}
	switch strings.ToLower(strings.TrimSpace(gender)) {
	case "", v1.FamilyMemberGenderMale, v1.FamilyMemberGenderFemale:
		return nil
	default:
		return gerror.NewCode(gcode.CodeInvalidParameter, "gender must be male or female")
	}
}
