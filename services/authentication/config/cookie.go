package config

import (
	"github.com/spf13/viper"
	"time"
)

type Cookie struct {
	MysqlTime time.Duration
}

func InitCookie(cfg *viper.Viper) *Cookie {
	return &Cookie{
		MysqlTime: cfg.GetDuration("mysqlTime"),
	}
}

var CookieConfig = new(Cookie)
