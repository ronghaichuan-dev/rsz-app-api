package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	commonlogin "rslytics-app-api/internal/common/login"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

var publicKidsPaths = map[string]struct{}{
	"/api.json":            {},
	"/swagger":             {},
	"/v1/health":           {},
	"/v1/kids/users/login": {},
	"/favicon.ico":         {},
}

// KidsJWT 校验 kids 微服务除公开接口以外的 JWT，并把用户ID写入请求上下文。
func KidsJWT(r *ghttp.Request) {
	if isPublicKidsPath(r.URL.Path) {
		r.Middleware.Next()
		return
	}
	token := utils.ExtractBearerToken(r.Header.Get("Authorization"))
	if token == "" {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "authorization token is required"))
		return
	}
	claims, err := utils.ParseJWT(token, commonlogin.SecretFromConfig(r.Context()))
	if err != nil {
		r.SetError(err)
		return
	}
	if !commonlogin.TokenExists(r.Context(), utils.KidsDB(r.Context()), consts.KidsUserTokenTable, claims.UserId, token) {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "token is invalid or expired"))
		return
	}
	r.SetCtxVar(consts.CtxUserIdKey, claims.UserId)
	r.SetCtxVar(consts.CtxTokenKey, token)
	r.Middleware.Next()
}

// isPublicKidsPath 判断请求路径是否属于公开接口。
func isPublicKidsPath(path string) bool {
	if _, ok := publicKidsPaths[path]; ok {
		return true
	}
	// 合同路由使用自己的 principal 解析，不能由 legacy JWT 中间件提前拦截。
	if strings.HasPrefix(path, "/v1/") && !strings.HasPrefix(path, "/v1/kids/") && path != "/v1/health" {
		return true
	}
	return strings.HasPrefix(path, "/swagger/")
}
