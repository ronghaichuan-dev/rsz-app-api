package kids

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"

	"rslytics-app-api/internal/middleware"
)

// TestV1StatisticsRoutesRegistered 确保统计和统计对比路由命中 v1 处理器，并按未鉴权规则返回 401。
func TestV1StatisticsRoutesRegistered(t *testing.T) {
	server := g.Server(guid.S())
	server.SetPort(0)
	server.SetDumpRouterMap(false)
	server.Use(middleware.Ctx, middleware.V1Envelope)
	Register(context.Background(), server.Group(""))
	server.Start()
	defer server.Shutdown()

	time.Sleep(100 * time.Millisecond)
	client := &http.Client{Timeout: 5 * time.Second}
	paths := []string{
		"/v1/circles/circle:v1:test/statistics",
		"/v1/circles/circle:v1:test/statistics:compare",
	}
	for _, path := range paths {
		response, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d%s", server.GetListenedPort(), path))
		if err != nil {
			t.Fatalf("统计路由请求失败 path=%s err=%v", path, err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("统计路由没有进入 v1 鉴权处理 path=%s got=%d want=%d", path, response.StatusCode, http.StatusUnauthorized)
		}
	}
}
