package consts

const (
	// HeaderLanguage 是客户端传递语言标识的请求头名称。
	HeaderLanguage = "Language"
	// HeaderAcceptLanguage 是 HTTP 标准语言请求头名称。
	HeaderAcceptLanguage = "Accept-Language"
	// DefaultLanguage 是未传语言或语言不支持时使用的默认语言。
	DefaultLanguage = "zh-CN"
	// FallbackLanguage 是默认语言缺少翻译时的兜底语言。
	FallbackLanguage = "en"
	// I18nPath 是多语言文件目录。
	I18nPath = "manifest/i18n"
)

// SupportedLanguages 是 kids 接口当前支持的语言列表。
var SupportedLanguages = map[string]struct{}{
	"zh-CN": {},
	"en":    {},
}
