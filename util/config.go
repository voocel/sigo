package util

import (
	"github.com/go-ini/ini"
	"os"
)

var ProxyConfig map[string]string

type EnvConfig os.File

func init() {
	ProxyConfig = make(map[string]string)
	EnvConfig, err := ini.Load("env")
	if err != nil {
		panic(err)
	}
	proxy, _ := EnvConfig.GetSection("proxy")
	if proxy != nil {
		secs := proxy.ChildSections()
		for _, sec := range secs {
			path, _ := sec.GetKey("path")
			pass, _ := sec.GetKey("pass")
			if path != nil && pass != nil {
				ProxyConfig[path.Value()] = pass.Value()
			}
		}
	}
}
