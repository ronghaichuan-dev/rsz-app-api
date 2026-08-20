package i18n

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"

	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

var initOnce sync.Once

// Init 初始化全局多语言管理器，服务启动时调用一次即可。
func Init() {
	initOnce.Do(func() {
		if err := gi18n.SetPath(resolvePath()); err != nil {
			panic(err)
		}
		gi18n.SetLanguage(consts.DefaultLanguage)
	})
}

// LanguageFromRequest 从 Language 请求头解析语言，兼容 Accept-Language。
func LanguageFromRequest(r *ghttp.Request) string {
	if r == nil {
		return consts.DefaultLanguage
	}
	language := utils.NormalizeLanguage(r.Header.Get(consts.HeaderLanguage))
	if language != "" {
		return language
	}
	language = utils.NormalizeLanguage(r.Header.Get(consts.HeaderAcceptLanguage))
	if language != "" {
		return language
	}
	return consts.DefaultLanguage
}

// WithRequestLanguage 把请求语言写入上下文，后续业务和响应统一使用该语言。
func WithRequestLanguage(r *ghttp.Request) {
	language := LanguageFromRequest(r)
	ctx := gi18n.WithLanguage(r.Context(), language)
	r.SetCtx(ctx)
}

// T 根据上下文语言翻译文本，缺少翻译时返回原文本。
func T(ctx context.Context, text string) string {
	if strings.TrimSpace(text) == "" {
		return text
	}
	return gi18n.T(ctx, text)
}

// Tf 根据上下文语言翻译格式化文本。
func Tf(ctx context.Context, format string, values ...any) string {
	if strings.TrimSpace(format) == "" {
		return format
	}
	return gi18n.Tf(ctx, format, values...)
}

// ErrorMessage 返回适合接口展示的错误文案。
func ErrorMessage(ctx context.Context, err error) string {
	if err == nil {
		return T(ctx, "OK")
	}
	if validErr, ok := err.(gvalid.Error); ok {
		if firstErr := validErr.FirstError(); firstErr != nil {
			return T(ctx, firstErr.Error())
		}
	}
	current := gerror.Current(err)
	if current != nil {
		if gfErr, ok := current.(*gerror.Error); ok {
			if gfErr.Text() == "" {
				return T(ctx, gerror.Code(err).Message())
			}
			return fmt.Sprintf(T(ctx, gfErr.Text()), gfErr.Args()...)
		}
		return T(ctx, current.Error())
	}
	return T(ctx, err.Error())
}

// resolvePath 从当前运行目录向上查找多语言文件目录，兼容测试和服务启动场景。
func resolvePath() string {
	wd, err := os.Getwd()
	if err != nil {
		return consts.I18nPath
	}
	for {
		path := filepath.Join(wd, consts.I18nPath)
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return path
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return consts.I18nPath
}
