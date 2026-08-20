package kids

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "rslytics-app-api/internal/api/kids/v1"
	commoni18n "rslytics-app-api/internal/common/i18n"
	commonlogin "rslytics-app-api/internal/common/login"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
	"rslytics-app-api/internal/utils"
)

// kidsLoginConfig 返回 kids 微服务调用公共登录模块所需的数据库和表配置。
func kidsLoginConfig() commonlogin.Config {
	return commonlogin.Config{
		DB: func(ctx context.Context) gdb.DB {
			return utils.KidsDB(ctx)
		},
		UserTable:      consts.KidsUserTable,
		UserAuthTable:  consts.KidsUserAuthTable,
		UserTokenTable: consts.KidsUserTokenTable,
		AccessTokenTTL: consts.KidsAccessTokenTTL,
		DefaultNickname: func(ctx context.Context, provider string) string {
			if provider == consts.LoginProviderGuest {
				return commoni18n.T(ctx, "Guest")
			}
			return commoni18n.Tf(ctx, "user.provider_name", utils.DefaultProviderName(provider))
		},
	}
}

// Login 根据登录方式完成 kids 用户登录，并委托公共登录模块持久化用户、授权身份和 token。
func (s *sKids) Login(ctx context.Context, in v1.UserLoginInput) (*v1.UserLoginOutput, error) {
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	deviceId := strings.TrimSpace(in.DeviceId)
	if deviceId == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "deviceId is required")
	}

	switch provider {
	case consts.LoginProviderGuest, consts.LoginProviderGoogle, consts.LoginProviderApple:
		out, err := service.Common().Login(ctx, kidsLoginConfig(), loginInputFromKids(in))
		if err != nil {
			return nil, err
		}
		return kidsLoginOutput(ctx, out, false)
	case consts.LoginProviderInvite:
		return s.loginInvite(ctx, deviceId, in)
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "unsupported login provider")
	}
}

// GetProfile 从数据库读取当前用户资料；正式鉴权接入前允许通过 userId 查询。
func (s *sKids) GetProfile(ctx context.Context, in v1.ProfileGetInput) (*v1.ProfileGetOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsUserTable).Ctx(ctx)
	if in.UserId > 0 {
		model = model.Where("id", in.UserId)
	}
	record, err := model.OrderDesc("id").One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return &v1.ProfileGetOutput{Id: 0, Nickname: commoni18n.T(ctx, "Guest"), Provider: consts.LoginProviderGuest, IsGuest: true}, nil
	}
	userId := record["id"].Uint64()
	return &v1.ProfileGetOutput{
		Id:       userId,
		Nickname: record["nickname"].String(),
		Avatar:   record["avatar"].String(),
		Provider: record["provider"].String(),
		IsGuest:  record["is_guest"].Uint() == 1,
		Stars:    getStarBalanceValue(ctx, userId),
	}, nil
}

// loginInvite 复用公共游客登录流程，并在同一事务中使用邀请码加入对应圈子。
func (s *sKids) loginInvite(ctx context.Context, deviceId string, in v1.UserLoginInput) (*v1.UserLoginOutput, error) {
	code := strings.TrimSpace(in.InviteCode)
	if code == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "inviteCode is required")
	}

	var joinedCircle v1.CircleInfo
	out, err := service.Common().LoginGuest(ctx, kidsLoginConfig(), loginInputFromKids(in), commonlogin.WithUserHook(func(ctx context.Context, tx gdb.TX, user *commonlogin.UserRecord) error {
		circle, err := joinCircleByInviteTx(ctx, tx, user.Id, code)
		if err != nil {
			return err
		}
		joinedCircle = circle
		return nil
	}))
	if err != nil {
		return nil, err
	}
	res, err := kidsLoginOutput(ctx, out, true)
	if err != nil {
		return nil, err
	}
	res.Provider = consts.LoginProviderInvite
	res.HasCircle = true
	res.CircleId = joinedCircle.Id
	res.CircleRole = joinedCircle.Role
	return res, nil
}

// loginInputFromKids 将 kids 登录入参转换为公共登录模块入参。
func loginInputFromKids(in v1.UserLoginInput) commonlogin.Input {
	return commonlogin.Input{
		Provider:      in.Provider,
		DeviceId:      in.DeviceId,
		AuthCode:      in.AuthCode,
		IdentityToken: in.IdentityToken,
		OpenId:        in.OpenId,
		Email:         in.Email,
		Nickname:      in.Nickname,
		Avatar:        in.Avatar,
	}
}

// kidsLoginOutput 为公共登录结果补充 kids 当前圈子信息。
func kidsLoginOutput(ctx context.Context, out *commonlogin.Output, joinedByInvite bool) (*v1.UserLoginOutput, error) {
	if out == nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "login output is empty")
	}
	circleId, circleRole, err := loginCircleInfo(ctx, nil, out.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.UserLoginOutput{
		UserId:         out.UserId,
		Token:          out.Token,
		Provider:       out.Provider,
		IsGuest:        out.IsGuest,
		BoundGuest:     out.BoundGuest,
		IsNewUser:      out.IsNewUser,
		DeviceId:       out.DeviceId,
		Nickname:       out.Nickname,
		Avatar:         out.Avatar,
		AccessExpire:   out.AccessExpire,
		HasCircle:      circleId > 0,
		CircleId:       circleId,
		CircleRole:     circleRole,
		JoinedByInvite: joinedByInvite,
	}, nil
}

// getStarBalanceValue 汇总星星流水得到当前余额。
func getStarBalanceValue(ctx context.Context, kidId uint64) int {
	value, err := utils.KidsDB(ctx).Model(consts.KidsStarRecordTable).Ctx(ctx).Where("kid_id", kidId).OrderDesc("id").Value("balance")
	if err != nil || value == nil {
		return 0
	}
	return gconv.Int(value)
}

// loginCircleInfo 查询登录用户当前加入的第一个圈子信息，用于客户端判断进入创建还是加入流程。
func loginCircleInfo(ctx context.Context, tx gdb.TX, userId uint64) (uint64, string, error) {
	var model *gdb.Model
	if tx != nil {
		model = tx.Model(consts.KidsCircleUserTable).Ctx(ctx)
	} else {
		model = utils.KidsDB(ctx).Model(consts.KidsCircleUserTable).Ctx(ctx)
	}
	record, err := model.Where("user_id", userId).OrderDesc("id").One()
	if err != nil || record.IsEmpty() {
		return 0, "", err
	}
	return record["circle_id"].Uint64(), record["role"].String(), nil
}

// UpdateProfile 更新当前用户昵称和头像。
func (s *sKids) UpdateProfile(ctx context.Context, in v1.ProfileUpdateInput) (*v1.ProfileUpdateOutput, error) {
	if strings.TrimSpace(in.Nickname) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "nickname is required")
	}
	result, err := utils.KidsDB(ctx).Model(consts.KidsUserTable).Ctx(ctx).
		Where("id", in.UserId).
		Data(map[string]any{"nickname": strings.TrimSpace(in.Nickname), "avatar": strings.TrimSpace(in.Avatar)}).
		Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "user not found")
	}
	profile, err := s.GetProfile(ctx, v1.ProfileGetInput{UserId: in.UserId})
	if err != nil {
		return nil, err
	}
	return &v1.ProfileUpdateOutput{Profile: *profile}, nil
}
