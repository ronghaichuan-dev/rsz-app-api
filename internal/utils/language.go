package utils

import (
	"strings"

	"rslytics-app-api/internal/consts"
)

// NormalizeLanguage 统一语言标识格式并过滤不支持的语言。
func NormalizeLanguage(language string) string {
	language = strings.TrimSpace(language)
	if language == "" {
		return ""
	}
	if index := strings.Index(language, ","); index >= 0 {
		language = language[:index]
	}
	if index := strings.Index(language, ";"); index >= 0 {
		language = language[:index]
	}
	language = strings.ReplaceAll(language, "_", "-")
	lower := strings.ToLower(language)
	switch lower {
	case "zh", "zh-cn", "zh-hans", "cn":
		language = consts.DefaultLanguage
	case "en", "en-us", "en-gb":
		language = consts.FallbackLanguage
	default:
		language = ""
	}
	if _, ok := consts.SupportedLanguages[language]; !ok {
		return ""
	}
	return language
}
