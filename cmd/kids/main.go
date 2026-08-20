package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	_ "rslytics-app-api/internal/logic"
	"rslytics-app-api/internal/server/kids"
)

// main 启动 kids 微服务入口。
func main() {
	kids.Main.Run(gctx.GetInitCtx())
}
