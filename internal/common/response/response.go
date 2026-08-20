package response

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	commoni18n "rslytics-app-api/internal/common/i18n"
)

type JsonRes struct {
	Code    int    `json:"code" dc:"业务码"`
	Message string `json:"message" dc:"提示信息"`
	Data    any    `json:"data" dc:"响应数据"`
}

var streamContentTypes = map[string]struct{}{
	"text/event-stream":         {},
	"application/octet-stream":  {},
	"multipart/x-mixed-replace": {},
}

func Middleware(r *ghttp.Request) {
	r.Middleware.Next()
	if strings.HasPrefix(r.URL.Path, "/v1/") && !strings.HasPrefix(r.URL.Path, "/v1/kids/") && r.URL.Path != "/v1/health" {
		return
	}

	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}

	mediaType, _, _ := mime.ParseMediaType(r.Response.Header().Get("Content-Type"))
	if _, ok := streamContentTypes[mediaType]; ok {
		return
	}

	var (
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
		msg  string
	)
	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = commoni18n.ErrorMessage(r.Context(), err)
	} else {
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			r.SetError(gerror.NewCode(code))
		} else {
			code = gcode.CodeOK
		}
		msg = commoni18n.T(r.Context(), code.Message())
	}

	r.Response.WriteJson(JsonRes{
		Code:    code.Code(),
		Message: msg,
		Data:    res,
	})
}
