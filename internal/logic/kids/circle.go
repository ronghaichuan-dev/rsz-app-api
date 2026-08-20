package kids

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// CreateCircle 创建管理员圈子，并把当前用户加入为管理员。
func (s *sKids) CreateCircle(ctx context.Context, in v1.CircleCreateInput) (*v1.CircleCreateOutput, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "circle name is required")
	}
	var circle v1.CircleInfo
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model(consts.KidsCircleTable).Ctx(ctx).Data(map[string]any{
			"name":          strings.TrimSpace(in.Name),
			"icon":          strings.TrimSpace(in.Icon),
			"owner_user_id": in.UserId,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		if _, err = tx.Model(consts.KidsCircleUserTable).Ctx(ctx).Data(map[string]any{
			"circle_id": id,
			"user_id":   in.UserId,
			"role":      v1.CircleRoleAdmin,
		}).Insert(); err != nil {
			return err
		}
		circle = v1.CircleInfo{Id: uint64(id), Name: strings.TrimSpace(in.Name), Icon: strings.TrimSpace(in.Icon), Role: v1.CircleRoleAdmin, Joined: true}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.CircleCreateOutput{Circle: circle}, nil
}

// CreateInviteCode 创建管理员或成员邀请码，默认有效期为 72 小时。
func (s *sKids) CreateInviteCode(ctx context.Context, in v1.InviteCodeCreateInput) (*v1.InviteCodeCreateOutput, error) {
	if in.ExpireHours <= 0 {
		in.ExpireHours = 72
	}
	inviteRole := strings.TrimSpace(in.InviteRole)
	if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, in.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can create invite code")
	}
	var code string
	var expiredAt time.Time
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var err error
		code, err = uniqueInviteCode(ctx, tx)
		if err != nil {
			return err
		}
		expiredAt = time.Now().Add(time.Duration(in.ExpireHours) * time.Hour)
		_, err = tx.Model(consts.KidsInviteCodeTable).Ctx(ctx).Data(map[string]any{
			"code":             code,
			"circle_id":        in.CircleId,
			"inviter_user_id":  in.UserId,
			"invite_role":      inviteRole,
			"target_member_id": in.TargetMemberId,
			"expired_at":       expiredAt.Format(consts.MySQLTimeLayout),
		}).Insert()
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.InviteCodeCreateOutput{Code: code, ExpiredAt: expiredAt.Unix()}, nil
}

// PreviewInviteCode 查询邀请码对应的圈子和角色信息。
func (s *sKids) PreviewInviteCode(ctx context.Context, in v1.InviteCodePreviewInput) (*v1.InviteCodePreviewOutput, error) {
	record, err := inviteCodeRecord(ctx, nil, in.Code)
	if err != nil {
		return nil, err
	}
	circle, err := circleRecord(ctx, nil, record["circle_id"].Uint64())
	if err != nil {
		return nil, err
	}
	return &v1.InviteCodePreviewOutput{
		Code:       record["code"].String(),
		CircleId:   record["circle_id"].Uint64(),
		CircleName: circle["name"].String(),
		InviteRole: record["invite_role"].String(),
		ExpiredAt:  utils.ParseDBTime(record["expired_at"].Val()),
	}, nil
}

// JoinCircle 使用邀请码加入圈子，并按邀请码角色建立圈子用户关系。
func (s *sKids) JoinCircle(ctx context.Context, in v1.CircleJoinInput) (*v1.CircleJoinOutput, error) {
	var circle v1.CircleInfo
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		joinedCircle, err := joinCircleByInviteTx(ctx, tx, in.UserId, in.Code)
		if err != nil {
			return err
		}
		circle = joinedCircle
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.CircleJoinOutput{Circle: circle}, nil
}

// joinCircleByInviteTx 在同一事务中校验邀请码、建立圈子关系并标记邀请码已使用。
func joinCircleByInviteTx(ctx context.Context, tx gdb.TX, userId uint64, code string) (v1.CircleInfo, error) {
	invite, err := inviteCodeRecordForUpdate(ctx, tx, code)
	if err != nil {
		return v1.CircleInfo{}, err
	}
	usedByUserId := invite["used_by_user_id"].Uint64()
	if usedByUserId > 0 {
		return v1.CircleInfo{}, gerror.NewCode(gcode.CodeInvalidOperation, "invite code already used")
	}
	if utils.ParseDBTime(invite["expired_at"].Val()) <= time.Now().Unix() {
		return v1.CircleInfo{}, gerror.NewCode(gcode.CodeInvalidOperation, "invite code expired")
	}
	circle, err := circleRecord(ctx, tx, invite["circle_id"].Uint64())
	if err != nil {
		return v1.CircleInfo{}, err
	}
	if usedByUserId == 0 {
		if err = persistInviteJoin(ctx, tx, userId, invite); err != nil {
			return v1.CircleInfo{}, err
		}
	}
	return v1.CircleInfo{
		Id:     invite["circle_id"].Uint64(),
		Name:   circle["name"].String(),
		Icon:   circle["icon"].String(),
		Role:   invite["invite_role"].String(),
		Joined: true,
	}, nil
}

// persistInviteJoin 写入圈子关系、绑定目标家庭成员并消费邀请码。
func persistInviteJoin(ctx context.Context, tx gdb.TX, userId uint64, invite gdb.Record) error {
	now := time.Now().Format(consts.MySQLTimeLayout)
	memberId := invite["target_member_id"].Uint64()
	if _, err := tx.Model(consts.KidsCircleUserTable).Ctx(ctx).Data(map[string]any{
		"circle_id":  invite["circle_id"].Uint64(),
		"user_id":    userId,
		"role":       invite["invite_role"].String(),
		"member_id":  memberId,
		"deleted_at": nil,
		"left_at":    nil,
	}).Save(); err != nil {
		return err
	}
	if memberId > 0 {
		if _, err := tx.Model(consts.KidsFamilyMemberTable).Ctx(ctx).Where("id", memberId).Data(map[string]any{"bind_user_id": userId, "bound_at": now}).Update(); err != nil {
			return err
		}
	}
	_, err := tx.Model(consts.KidsInviteCodeTable).Ctx(ctx).Where("id", invite["id"].Uint64()).Data(map[string]any{
		"used_by_user_id": userId,
		"used_at":         now,
	}).Update()
	return err
}

// userIsCircleAdmin 判断用户是否是指定圈子的管理员。
func userIsCircleAdmin(ctx context.Context, tx gdb.TX, userId, circleId uint64) (bool, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsCircleUserTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsCircleUserTable).Ctx(ctx)
	}
	count, err := model.Where("user_id", userId).Where("circle_id", circleId).Where("role", v1.CircleRoleAdmin).Where("deleted_at IS NULL").Count()
	return count > 0, err
}

// uniqueInviteCode 生成数据库内唯一的六位数字邀请码。
func uniqueInviteCode(ctx context.Context, tx gdb.TX) (string, error) {
	for i := 0; i < 10; i++ {
		code, err := randomSixDigitCode()
		if err != nil {
			return "", err
		}
		count, err := tx.Model(consts.KidsInviteCodeTable).Ctx(ctx).Where("code", code).Count()
		if err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
	return "", gerror.NewCode(gcode.CodeInternalError, "generate invite code failed")
}

// randomSixDigitCode 生成六位数字邀请码。
func randomSixDigitCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// inviteCodeRecord 查询并校验邀请码是否存在。
func inviteCodeRecord(ctx context.Context, tx gdb.TX, code string) (gdb.Record, error) {
	code = strings.TrimSpace(code)
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsInviteCodeTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsInviteCodeTable).Ctx(ctx)
	}
	record, err := model.Where("code", code).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "invite code not found")
	}
	return record, nil
}

// inviteCodeRecordForUpdate 查询邀请码并加排他锁，避免同一邀请码被并发消费。
func inviteCodeRecordForUpdate(ctx context.Context, tx gdb.TX, code string) (gdb.Record, error) {
	if tx == nil {
		return inviteCodeRecord(ctx, tx, code)
	}
	code = strings.TrimSpace(code)
	record, err := tx.Model(consts.KidsInviteCodeTable).Ctx(ctx).Where("code", code).LockUpdate().One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "invite code not found")
	}
	return record, nil
}

// circleRecord 查询并校验圈子是否存在。
func circleRecord(ctx context.Context, tx gdb.TX, circleId uint64) (gdb.Record, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsCircleTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsCircleTable).Ctx(ctx)
	}
	record, err := model.Where("id", circleId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "circle not found")
	}
	return record, nil
}

// ListCircles 查询当前用户管理和加入的群组。
func (s *sKids) ListCircles(ctx context.Context, in v1.CircleListInput) (*v1.CircleListOutput, error) {
	records, err := utils.KidsDB(ctx).Model(consts.KidsCircleUserTable+" cu").Ctx(ctx).
		LeftJoin(consts.KidsCircleTable+" c", "c.id = cu.circle_id").
		Fields("c.id,c.name,c.icon,c.owner_user_id,cu.role").
		Where("cu.user_id", in.UserId).
		Where("cu.deleted_at IS NULL").
		Where("c.deleted_at IS NULL").
		OrderAsc("c.id").
		All()
	if err != nil {
		return nil, err
	}
	out := &v1.CircleListOutput{}
	for _, record := range records {
		item := v1.CircleInfo{Id: record["id"].Uint64(), Name: record["name"].String(), Icon: record["icon"].String(), Role: record["role"].String(), Joined: true}
		if item.Role == v1.CircleRoleAdmin {
			out.Managed = append(out.Managed, item)
		} else {
			out.Joined = append(out.Joined, item)
		}
	}
	return out, nil
}

// UpdateCircle 校验管理权限后更新群组名称和图标。
func (s *sKids) UpdateCircle(ctx context.Context, in v1.CircleUpdateInput) (*v1.CircleUpdateOutput, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "circle name is required")
	}
	if ok, err := userIsCircleAdmin(ctx, nil, in.UserId, in.Id); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle admin can manage circle")
	}
	_, err := utils.KidsDB(ctx).Model(consts.KidsCircleTable).Ctx(ctx).
		Where("id", in.Id).
		Where("deleted_at IS NULL").
		Data(map[string]any{"name": strings.TrimSpace(in.Name), "icon": strings.TrimSpace(in.Icon)}).
		Update()
	if err != nil {
		return nil, err
	}
	record, err := circleRecord(ctx, nil, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.CircleUpdateOutput{Circle: v1.CircleInfo{Id: in.Id, Name: record["name"].String(), Icon: record["icon"].String(), Role: v1.CircleRoleAdmin, Joined: true}}, nil
}

// DeleteCircle 仅群组所有者可以软删除群组。
func (s *sKids) DeleteCircle(ctx context.Context, in v1.CircleDeleteInput) (*v1.CircleDeleteOutput, error) {
	if ok, err := userIsCircleOwner(ctx, nil, in.UserId, in.Id); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle owner can delete circle")
	}
	err := utils.KidsDB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := time.Now().Format(consts.MySQLTimeLayout)
		if _, err := tx.Model(consts.KidsCircleTable).Ctx(ctx).Where("id", in.Id).Data(map[string]any{"deleted_at": now}).Update(); err != nil {
			return err
		}
		_, err := tx.Model(consts.KidsCircleUserTable).Ctx(ctx).Where("circle_id", in.Id).Where("deleted_at IS NULL").Data(map[string]any{"deleted_at": now, "left_at": now}).Update()
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.CircleDeleteOutput{Id: in.Id}, nil
}

// LeaveCircle 允许非所有者退出已加入群组。
func (s *sKids) LeaveCircle(ctx context.Context, in v1.CircleLeaveInput) (*v1.CircleLeaveOutput, error) {
	if ok, err := userIsCircleOwner(ctx, nil, in.UserId, in.Id); err != nil || ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "circle owner cannot leave circle")
	}
	result, err := utils.KidsDB(ctx).Model(consts.KidsCircleUserTable).Ctx(ctx).
		Where("circle_id", in.Id).
		Where("user_id", in.UserId).
		Where("deleted_at IS NULL").
		Data(map[string]any{"deleted_at": time.Now().Format(consts.MySQLTimeLayout), "left_at": time.Now().Format(consts.MySQLTimeLayout)}).
		Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "circle relation not found")
	}
	return &v1.CircleLeaveOutput{Id: in.Id}, nil
}

// ListCircleMembers 查询群组管理员和成员绑定状态。
func (s *sKids) ListCircleMembers(ctx context.Context, in v1.CircleMemberListInput) (*v1.CircleMemberListOutput, error) {
	if ok, err := userCanAccessCircle(ctx, nil, in.UserId, in.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "circle permission denied")
	}
	circle, err := circleRecord(ctx, nil, in.CircleId)
	if err != nil {
		return nil, err
	}
	out := &v1.CircleMemberListOutput{}
	adminRecords, err := utils.KidsDB(ctx).Model(consts.KidsCircleUserTable+" cu").Ctx(ctx).
		LeftJoin(consts.KidsUserTable+" u", "u.id = cu.user_id").
		Fields("cu.user_id,cu.member_id,cu.role,cu.created_at,u.nickname,u.avatar").
		Where("cu.circle_id", in.CircleId).
		Where("cu.deleted_at IS NULL").
		Where("cu.role", v1.CircleRoleAdmin).
		OrderAsc("cu.id").All()
	if err != nil {
		return nil, err
	}
	for _, record := range adminRecords {
		item := v1.CircleMemberItem{UserId: record["user_id"].Uint64(), MemberId: record["member_id"].Uint64(), Name: record["nickname"].String(), Avatar: record["avatar"].String(), Role: record["role"].String(), Bound: true, IsOwner: record["user_id"].Uint64() == circle["owner_user_id"].Uint64(), CreatedAt: utils.ParseDBTime(record["created_at"].Val())}
		if item.IsOwner {
			out.Owner = item
		} else {
			out.Managers = append(out.Managers, item)
		}
	}
	memberRecords, err := utils.KidsDB(ctx).Model(consts.KidsFamilyMemberTable).Ctx(ctx).
		Where("circle_id", in.CircleId).
		Where("deleted_at IS NULL").
		OrderAsc("id").All()
	if err != nil {
		return nil, err
	}
	for _, record := range memberRecords {
		out.Members = append(out.Members, v1.CircleMemberItem{UserId: record["bind_user_id"].Uint64(), MemberId: record["id"].Uint64(), Name: record["name"].String(), Avatar: record["avatar"].String(), Role: v1.CircleRoleMember, Bound: record["bind_user_id"].Uint64() > 0, IsOwner: record["owner"].Int() == 1, CreatedAt: utils.ParseDBTime(record["created_at"].Val())})
	}
	return out, nil
}

// RemoveCircleAdmin 仅群组所有者可以移除管理员。
func (s *sKids) RemoveCircleAdmin(ctx context.Context, in v1.CircleAdminRemoveInput) (*v1.CircleAdminRemoveOutput, error) {
	if ok, err := userIsCircleOwner(ctx, nil, in.OperatorUserId, in.CircleId); err != nil || !ok {
		if err != nil {
			return nil, err
		}
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "only circle owner can remove admin")
	}
	if in.OperatorUserId == in.AdminUserId {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "circle owner cannot remove self")
	}
	result, err := utils.KidsDB(ctx).Model(consts.KidsCircleUserTable).Ctx(ctx).
		Where("circle_id", in.CircleId).
		Where("user_id", in.AdminUserId).
		Where("role", v1.CircleRoleAdmin).
		Where("deleted_at IS NULL").
		Data(map[string]any{"deleted_at": time.Now().Format(consts.MySQLTimeLayout), "left_at": time.Now().Format(consts.MySQLTimeLayout)}).
		Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "circle admin not found")
	}
	return &v1.CircleAdminRemoveOutput{CircleId: in.CircleId, AdminUserId: in.AdminUserId}, nil
}

// userIsCircleOwner 判断用户是否是指定圈子的创建者。
func userIsCircleOwner(ctx context.Context, tx gdb.TX, userId, circleId uint64) (bool, error) {
	record, err := circleRecord(ctx, tx, circleId)
	if err != nil {
		return false, err
	}
	return record["owner_user_id"].Uint64() == userId, nil
}

// userCanAccessCircle 判断用户是否仍在圈子中。
func userCanAccessCircle(ctx context.Context, tx gdb.TX, userId, circleId uint64) (bool, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsCircleUserTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsCircleUserTable).Ctx(ctx)
	}
	count, err := model.Where("user_id", userId).Where("circle_id", circleId).Where("deleted_at IS NULL").Count()
	return count > 0, err
}
