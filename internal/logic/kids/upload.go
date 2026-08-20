package kids

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// UploadFile 登记上传文件并返回文件地址。
func (s *sKids) UploadFile(ctx context.Context, in v1.UploadFileInput) (*v1.UploadFileOutput, error) {
	if strings.TrimSpace(in.FileUrl) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "file url is required")
	}
	if strings.TrimSpace(in.UsageType) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "usage type is required")
	}
	if !validUploadUsageType(in.UsageType) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "unsupported usage type")
	}
	id, err := utils.KidsDB(ctx).Model(consts.KidsUploadFileTable).Ctx(ctx).Data(map[string]any{
		"user_id":      in.UserId,
		"member_id":    in.MemberId,
		"usage_type":   strings.TrimSpace(in.UsageType),
		"file_name":    strings.TrimSpace(in.FileName),
		"file_url":     strings.TrimSpace(in.FileUrl),
		"content_type": strings.TrimSpace(in.ContentType),
		"file_size":    in.FileSize,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &v1.UploadFileOutput{Id: uint64(id), FileUrl: strings.TrimSpace(in.FileUrl)}, nil
}

// validUploadUsageType 校验上传文件用途枚举。
func validUploadUsageType(usageType string) bool {
	switch strings.TrimSpace(usageType) {
	case v1.UploadUsageAvatar, v1.UploadUsageTaskPhoto, v1.UploadUsageReward:
		return true
	default:
		return false
	}
}
