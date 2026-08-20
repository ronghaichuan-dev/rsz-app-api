package upload

import (
	"context"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/service"
)

// Controller 处理上传文件登记接口。
type Controller struct{}

// New 创建上传控制器实例。
func New() *Controller {
	return &Controller{}
}

// Create 登记已上传的照片或头像文件。
func (c *Controller) Create(ctx context.Context, req *v1.UploadFileReq) (res *v1.UploadFileRes, err error) {
	out, err := service.Kids().UploadFile(ctx, v1.UploadFileInput{UserId: authUserId(ctx), MemberId: req.MemberId, UsageType: req.UsageType, FileName: req.FileName, FileUrl: req.FileUrl, ContentType: req.ContentType, FileSize: req.FileSize})
	if err != nil {
		return nil, err
	}
	return &v1.UploadFileRes{Id: out.Id, FileUrl: out.FileUrl}, nil
}

// authUserId 从上下文读取当前登录用户ID。
func authUserId(ctx context.Context) uint64 {
	if value, ok := ctx.Value(consts.CtxUserIdKey).(uint64); ok {
		return value
	}
	return 0
}
