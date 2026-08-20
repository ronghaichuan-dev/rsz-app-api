package common

import (
	"context"
	commonlogin "rslytics-app-api/internal/common/login"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Login 根据登录方式完成通用登录，并把用户、授权身份和 token 持久化到调用方数据库。
func (s *sCommon) Login(ctx context.Context, cfg commonlogin.Config, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error) {
	cfg = commonlogin.NormalizeConfig(cfg)
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	deviceId := strings.TrimSpace(in.DeviceId)
	if deviceId == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "deviceId is required")
	}

	switch provider {
	case consts.LoginProviderGuest:
		return s.LoginGuest(ctx, cfg, in, opts...)
	case consts.LoginProviderGoogle, consts.LoginProviderApple:
		return s.loginOAuth(ctx, cfg, provider, deviceId, in, opts...)
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "unsupported login provider")
	}
}

// LoginGuest 创建或复用同设备游客账号，并允许调用方在同一事务内追加业务处理。
func (s *sCommon) LoginGuest(ctx context.Context, cfg commonlogin.Config, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error) {
	cfg = commonlogin.NormalizeConfig(cfg)
	deviceId := strings.TrimSpace(in.DeviceId)
	if deviceId == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "deviceId is required")
	}
	options := commonlogin.MergeOptions(opts...)

	var out *commonlogin.Output
	err := cfg.DB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		user, err := s.findGuestUserByDevice(ctx, tx, cfg, deviceId)
		if err != nil {
			return err
		}
		isNewUser := false
		if user == nil {
			user, err = s.createUser(ctx, tx, cfg, map[string]any{
				"device_id": deviceId,
				"provider":  consts.LoginProviderGuest,
				"nickname":  utils.DefaultString(in.Nickname, cfg.DefaultNickname(ctx, consts.LoginProviderGuest)),
				"avatar":    strings.TrimSpace(in.Avatar),
				"is_guest":  1,
			})
			if err != nil {
				return err
			}
			isNewUser = true
		}
		if options.UserHook != nil {
			if err = options.UserHook(ctx, tx, user); err != nil {
				return err
			}
			user, err = s.findUserById(ctx, tx, cfg, user.Id)
			if err != nil {
				return err
			}
		}
		out, err = s.buildLoginOutput(ctx, tx, cfg, user, false, isNewUser)
		return err
	})
	return out, err
}

// loginOAuth 在事务中完成授权登录，并按规则优先绑定同设备未绑定游客账号。
func (s *sCommon) loginOAuth(ctx context.Context, cfg commonlogin.Config, provider, deviceId string, in commonlogin.Input, opts ...commonlogin.Option) (*commonlogin.Output, error) {
	identity, err := s.verifyLoginOAuthIdentity(ctx, cfg, provider, in.IdentityToken)
	if err != nil {
		return nil, err
	}
	openId := identity.OpenId
	options := commonlogin.MergeOptions(opts...)

	var out *commonlogin.Output
	err = cfg.DB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		user, err := s.findUserByAuth(ctx, tx, cfg, provider, openId)
		if err != nil {
			return err
		}
		if user != nil {
			if err = s.updateExistingOAuthUser(ctx, tx, cfg, user.Id, deviceId, identity); err != nil {
				return err
			}
			user, err = s.findUserById(ctx, tx, cfg, user.Id)
			if err != nil {
				return err
			}
			if options.UserHook != nil {
				if err = options.UserHook(ctx, tx, user); err != nil {
					return err
				}
			}
			out, err = s.buildLoginOutput(ctx, tx, cfg, user, false, false)
			return err
		}

		boundGuest := false
		isNewUser := false
		user, err = s.findGuestUserByDevice(ctx, tx, cfg, deviceId)
		if err != nil {
			return err
		}
		if user != nil {
			if err = s.bindGuestUser(ctx, tx, cfg, user.Id, provider, identity); err != nil {
				return err
			}
			boundGuest = true
		} else {
			user, err = s.createUser(ctx, tx, cfg, map[string]any{
				"device_id": deviceId,
				"provider":  provider,
				"email":     identity.Email,
				"nickname":  utils.DefaultString(identity.Nickname, cfg.DefaultNickname(ctx, provider)),
				"avatar":    identity.Avatar,
				"is_guest":  0,
			})
			if err != nil {
				return err
			}
			isNewUser = true
		}
		if err = s.createUserAuth(ctx, tx, cfg, user.Id, provider, openId, identity.Email); err != nil {
			return err
		}
		user, err = s.findUserById(ctx, tx, cfg, user.Id)
		if err != nil {
			return err
		}
		if options.UserHook != nil {
			if err = options.UserHook(ctx, tx, user); err != nil {
				return err
			}
		}
		out, err = s.buildLoginOutput(ctx, tx, cfg, user, boundGuest, isNewUser)
		return err
	})
	return out, err
}

// findGuestUserByDevice 查询同一设备未绑定授权身份的游客账号。
func (s *sCommon) findGuestUserByDevice(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, deviceId string) (*commonlogin.UserRecord, error) {
	return s.scanLoginUserRecord(tx.Model(cfg.UserTable).Ctx(ctx).Where("device_id", deviceId).Where("is_guest", 1).OrderAsc("id"))
}

// findUserByAuth 根据授权提供方和 openId 查询已绑定账号。
func (s *sCommon) findUserByAuth(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, provider, openId string) (*commonlogin.UserRecord, error) {
	record, err := tx.Model(cfg.UserAuthTable+" ua").Ctx(ctx).
		LeftJoin(cfg.UserTable+" u", "u.id = ua.user_id").
		Fields("u.id,u.device_id,u.provider,u.email,u.nickname,u.avatar,u.is_guest").
		Where("ua.provider", provider).
		Where("ua.open_id", openId).
		One()
	if err != nil || record.IsEmpty() {
		return nil, err
	}
	return userRecordFromDB(record), nil
}

// findUserById 根据用户 ID 查询用户。
func (s *sCommon) findUserById(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, userId uint64) (*commonlogin.UserRecord, error) {
	return s.scanLoginUserRecord(tx.Model(cfg.UserTable).Ctx(ctx).Where("id", userId))
}

// scanLoginUserRecord 查询单条用户记录并转换为业务结构。
func (s *sCommon) scanLoginUserRecord(model *gdb.Model) (*commonlogin.UserRecord, error) {
	record, err := model.One()
	if err != nil || record.IsEmpty() {
		return nil, err
	}
	return userRecordFromDB(record), nil
}

// userRecordFromDB 将数据库记录转换为登录流程内部结构。
func userRecordFromDB(record gdb.Record) *commonlogin.UserRecord {
	return &commonlogin.UserRecord{
		Id:       record["id"].Uint64(),
		DeviceId: record["device_id"].String(),
		Provider: record["provider"].String(),
		Email:    record["email"].String(),
		Nickname: record["nickname"].String(),
		Avatar:   record["avatar"].String(),
		IsGuest:  record["is_guest"].Uint(),
	}
}

// createUser 创建用户并返回数据库生成的用户 ID。
func (s *sCommon) createUser(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, data map[string]any) (*commonlogin.UserRecord, error) {
	id, err := tx.Model(cfg.UserTable).Ctx(ctx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return s.findUserById(ctx, tx, cfg, uint64(id))
}

// updateExistingOAuthUser 更新已存在授权账号的最近设备和服务端可信展示资料。
func (s *sCommon) updateExistingOAuthUser(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, userId uint64, deviceId string, identity commonlogin.OAuthIdentity) error {
	data := map[string]any{"device_id": deviceId}
	if strings.TrimSpace(identity.Email) != "" {
		data["email"] = strings.TrimSpace(identity.Email)
	}
	if strings.TrimSpace(identity.Nickname) != "" {
		data["nickname"] = strings.TrimSpace(identity.Nickname)
	}
	if strings.TrimSpace(identity.Avatar) != "" {
		data["avatar"] = strings.TrimSpace(identity.Avatar)
	}
	_, err := tx.Model(cfg.UserTable).Ctx(ctx).Where("id", userId).Data(data).Update()
	return err
}

// bindGuestUser 将游客账号升级为正式授权账号，并保留原 userId。
func (s *sCommon) bindGuestUser(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, userId uint64, provider string, identity commonlogin.OAuthIdentity) error {
	data := map[string]any{
		"provider": provider,
		"email":    strings.TrimSpace(identity.Email),
		"is_guest": 0,
	}
	if strings.TrimSpace(identity.Nickname) != "" {
		data["nickname"] = strings.TrimSpace(identity.Nickname)
	}
	if strings.TrimSpace(identity.Avatar) != "" {
		data["avatar"] = strings.TrimSpace(identity.Avatar)
	}
	_, err := tx.Model(cfg.UserTable).Ctx(ctx).Where("id", userId).Data(data).Update()
	return err
}

// createUserAuth 持久化授权身份绑定关系。
func (s *sCommon) createUserAuth(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, userId uint64, provider, openId, email string) error {
	_, err := tx.Model(cfg.UserAuthTable).Ctx(ctx).Data(map[string]any{
		"user_id":  userId,
		"provider": provider,
		"open_id":  openId,
		"email":    strings.TrimSpace(email),
	}).Insert()
	return err
}

// buildLoginOutput 生成访问令牌、写入 token 表，并组装登录响应。
func (s *sCommon) buildLoginOutput(ctx context.Context, tx gdb.TX, cfg commonlogin.Config, user *commonlogin.UserRecord, boundGuest, isNewUser bool) (*commonlogin.Output, error) {
	if user == nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "user not found")
	}
	expiresAt := time.Now().Add(cfg.AccessTokenTTL)
	token, err := utils.GenerateJWT(user.Id, expiresAt.Unix(), cfg.JWTSecret(ctx))
	if err != nil {
		return nil, err
	}
	if _, err = tx.Model(cfg.UserTokenTable).Ctx(ctx).Data(map[string]any{
		"user_id":    user.Id,
		"token":      token,
		"expired_at": expiresAt.Format(consts.MySQLTimeLayout),
	}).Insert(); err != nil {
		return nil, err
	}
	return &commonlogin.Output{
		UserId:       user.Id,
		Token:        token,
		Provider:     user.Provider,
		IsGuest:      user.IsGuest == 1,
		BoundGuest:   boundGuest,
		IsNewUser:    isNewUser,
		DeviceId:     user.DeviceId,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		AccessExpire: expiresAt.Unix(),
	}, nil
}

// verifyLoginOAuthIdentity 调用 Google 或 Apple 后台公钥校验 identityToken。
func (s *sCommon) verifyLoginOAuthIdentity(ctx context.Context, cfg commonlogin.Config, provider, identityToken string) (commonlogin.OAuthIdentity, error) {
	var (
		verified *utils.OAuthIdentity
		err      error
	)
	switch provider {
	case consts.LoginProviderGoogle:
		verified, err = utils.VerifyGoogleIdentityToken(ctx, identityToken, cfg.OAuthClientId(ctx, provider))
	case consts.LoginProviderApple:
		verified, err = utils.VerifyAppleIdentityToken(ctx, identityToken, cfg.OAuthClientId(ctx, provider))
	default:
		err = gerror.NewCode(gcode.CodeInvalidParameter, "unsupported login provider")
	}
	if err != nil {
		return commonlogin.OAuthIdentity{}, err
	}
	return commonlogin.OAuthIdentity{
		OpenId:   strings.TrimSpace(verified.OpenId),
		Email:    strings.TrimSpace(verified.Email),
		Nickname: strings.TrimSpace(verified.Nickname),
		Avatar:   strings.TrimSpace(verified.Avatar),
	}, nil
}
