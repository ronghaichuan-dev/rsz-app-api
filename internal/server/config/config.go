package config

import (
	"fmt"
	"os"
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

func SetFile(appName, env string) string {
	file := fmt.Sprintf("config/%s/config.%s.yaml", appName, env)
	if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		adapter.SetFileName(file)
	}
	return file
}
