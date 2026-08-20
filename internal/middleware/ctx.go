package middleware

import "github.com/gogf/gf/v2/net/ghttp"

func Ctx(r *ghttp.Request) {
	// Put trace id, auth user, tenant and other request-scoped data here later.
	r.Middleware.Next()
}
