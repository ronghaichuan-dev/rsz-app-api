package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
)

const DefaultEnv = "dev"

var supportedEnvs = map[string]struct{}{
	"dev":  {},
	"test": {},
	"prod": {},
}

// ResolveEnv 按命令行参数和环境变量优先级解析当前服务运行环境。
func ResolveEnv(parser *gcmd.Parser) string {
	env := ""
	if parser != nil {
		if v := parser.GetOpt("env"); v != nil {
			env = v.String()
		}
		if env == "" {
			if v := parser.GetOpt("e"); v != nil {
				env = v.String()
			}
		}
	}
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = os.Getenv("GF_ENV")
	}
	if env == "" {
		env = os.Getenv("GO_ENV")
	}

	env = strings.ToLower(strings.TrimSpace(env))
	if _, ok := supportedEnvs[env]; !ok {
		return DefaultEnv
	}
	return env
}

// SetFile 按环境选择配置文件；设置 APP_CONFIG_DIR 时优先使用服务器私有配置目录。
func SetFile(appName, env string) string {
	file := ""
	if configDir := strings.TrimSpace(os.Getenv("APP_CONFIG_DIR")); configDir != "" {
		file = filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
	} else {
		file = fmt.Sprintf("config/%s/config.%s.yaml", appName, env)
	}
	if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		adapter.SetFileName(file)
	}
	return file
}
