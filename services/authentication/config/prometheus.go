package config

import "github.com/spf13/viper"

type Prometheus struct {
	UrlList []string
}

func InitPrometheus(cfg *viper.Viper) *Prometheus {
	return &Prometheus{
		UrlList: cfg.GetStringSlice("url"),
	}
}

var PrometheusConfig = new(Prometheus)
