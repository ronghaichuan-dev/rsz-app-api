package i18n

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

// TestTranslateByContextLanguage 验证上下文语言可以正确读取多语言文件。
func TestTranslateByContextLanguage(t *testing.T) {
	Init()

	zhCtx := gi18n.WithLanguage(context.Background(), "zh-CN")
	if got := T(zhCtx, "OK"); got != "成功" {
		t.Fatalf("期望中文成功文案，实际得到 %q", got)
	}
	if got := Tf(zhCtx, "error.kid_not_found", 3); got != "儿童 3 不存在" {
		t.Fatalf("期望中文格式化文案，实际得到 %q", got)
	}

	enCtx := gi18n.WithLanguage(context.Background(), "en")
	if got := T(enCtx, "authorization token is required"); got != "Authorization token is required" {
		t.Fatalf("期望英文鉴权文案，实际得到 %q", got)
	}
}

// TestValidationErrorMessage 验证 GoFrame 内置校验错误按上下文语言翻译。
func TestValidationErrorMessage(t *testing.T) {
	Init()

	zhCtx := gi18n.WithLanguage(context.Background(), "zh-CN")
	zhErr := g.Validator().Data("").Rules("required|min:1").Run(zhCtx)
	if got := ErrorMessage(zhCtx, zhErr); got != "不能为空" {
		t.Fatalf("期望中文校验文案，实际得到 %q", got)
	}

	enCtx := gi18n.WithLanguage(context.Background(), "en")
	enErr := g.Validator().Data("other").Rules("in:admin,member").Run(enCtx)
	if got := ErrorMessage(enCtx, enErr); got != "The field must be one of: admin,member" {
		t.Fatalf("期望英文校验文案，实际得到 %q", got)
	}
}
