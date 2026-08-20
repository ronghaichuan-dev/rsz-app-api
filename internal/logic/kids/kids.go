package kids

import (
	"context"

	"rslytics-app-api/internal/service"
)

// sKids 是 kids 微服务业务接口的默认实现。
type sKids struct{}

// init 在包加载时注册 kids 业务实现到 service 层。
func init() {
	service.RegisterKids(New())
}

// New 创建 kids 微服务业务实现实例。
func New() *sKids {
	return &sKids{}
}

// Get 保留 GoFrame service 示例方法，当前不承载业务逻辑。
func (s *sKids) Get(ctx context.Context) error {
	return nil
}
