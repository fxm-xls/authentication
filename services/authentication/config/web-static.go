package config

import "github.com/spf13/viper"

type WebStatic struct {
	Path string
}

func InitWebStatic(cfg *viper.Viper) *WebStatic {
	return &WebStatic{
		Path: cfg.GetString("path"),
	}
}

var WebStaticConfig = new(WebStatic)
