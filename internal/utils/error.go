package utils

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// NewCodef 使用稳定多语言 key 创建带格式化参数的业务错误。
func NewCodef(code gcode.Code, key string, args ...any) error {
	return gerror.NewCodef(code, key, args...)
}
